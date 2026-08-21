package users

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm"
)

func LoginAlur(c *gin.Context) {
	var person schemas.LoginSchema
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user schemas.Customer
	err := databases.DB.Where("email = ?", person.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Email atau password salah"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(person.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email atau password salah"})
		return
	}
	LegthToken, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot generate jwt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"user":    LegthToken,
	})
}

func RootWeb(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login")
}

var JwtKey string = os.Getenv("JWT")

func GenerateToken(nama string) (string, error) {
	expirateTime := time.Now().Add(2 * time.Hour)

	claim := &schemas.Claims{
		IdUser: nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirateTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
	return token.SignedString(JwtKey)
}

func AuthMiddlewareJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Token needed"})
			c.Abort()
			return
		}
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Format Header tidak ada atau rusak"})
			c.Abort()
			return
		}
		tokenString := tokenParts[1]
		claims := &schemas.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("Method signning is not Valid")
			}
			return JwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak Valid"})
			c.Abort()
			return
		}
		c.Set("id", claims.IdUser)
		c.Next()
	}
}

func SignupAlur(c *gin.Context) {
	var person schemas.SignupSchema
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser schemas.Customer
	err := databases.DB.Where("email = ?", person.Email).First(&existingUser).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email sudah terdaftar, silakan login"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa email di database"})
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(person.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}

	addUser := schemas.Customer{
		ID:          uuid.NewString(),
		Nama:        person.Nama,
		Email:       person.Email,
		NomerTelpon: person.NomorTelpon,
		Password:    string(bytes),
		Status:      "Unverified",
	}

	if err := databases.DB.Create(&addUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data user"})
		return
	}

	resendAPI := os.Getenv("RESEND_API")
	secretToken := os.Getenv("SECRET")

	if resendAPI != "" && secretToken != "" {
		client := resend.NewClient(resendAPI)
		verifyLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", secretToken)
		params := &resend.SendEmailRequest{
			From:    "Onboarding <onboarding@resend.dev>",
			To:      []string{person.Email},
			Subject: "Verifikasi Email Dev",
			Html:    fmt.Sprintf("<p>Klik link berikut untuk verifikasi: <a href='%s'>Verifikasi</a></p>", verifyLink),
		}

		if _, err := client.Emails.Send(params); err != nil {
			log.Printf("Gagal kirim email: %v\n", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Anda sudah terdaftar. Silakan periksa email untuk verifikasi.",
	})
}

func ViewPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You're in View Page"})
}

func DashboardPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You're in Dashboard Page"})
}

func SchedulePage(c *gin.Context) {
	ctx := c.Request.Context()

	srv, err := calendar.NewService(ctx,
		option.WithCredentialsFile("credentials.json"),
		option.WithScopes(calendar.CalendarReadonlyScope),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal inisialisasi Google Calendar Client (periksa file credentials.json)",
			"error":   err.Error(),
		})
		return
	}

	calendarID := "raushanfekra@gmail.com"
	t := time.Now().Format(time.RFC3339)

	events, err := srv.Events.List(calendarID).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(10).
		OrderBy("startTime").
		Do()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyambung ke Google Calendar",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "You are in Schedule Page",
		"data":    events.Items,
	})
}

func GetDataBooking(c *gin.Context) {
	var input schemas.BookingInput

	// 1. Tangkap dan Validasi Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Input tidak valid: " + err.Error()})
		return
	}

	// Atur zona waktu ke WIB (Asia/Jakarta)
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		loc = time.Local // Fallback ke zona waktu server jika Asia/Jakarta tidak ditemukan
	}

	var waktuMulai, waktuSelesai time.Time
	var sesiStr *string
	var jumlahHari *int

	// 2. Proses Logika Berdasarkan Tipe Booking
	if input.TipeBooking == "SESI" {
		// Validasi: Jika pilih SESI, tanggal awal dan akhir WAJIB sama
		if input.TanggalAwal != input.TanggalAkhir {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Untuk tipe SESI, tanggal_awal dan tanggal_akhir harus sama"})
			return
		}

		// Validasi: Sesi tidak boleh kosong
		if input.Sesi == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Sesi harus diisi (1, 2, atau 3)"})
			return
		}

		var jamMulai, jamSelesai string
		s := ""

		// Pemetaan Jam sesuai Sesi
		switch *input.Sesi {
		case 1:
			jamMulai = "08:00:00"
			jamSelesai = "10:00:00"
			s = "PAGI"
		case 2:
			jamMulai = "10:00:00"
			jamSelesai = "12:00:00"
			s = "SIANG"
		case 3:
			jamMulai = "13:00:00"
			jamSelesai = "15:00:00"
			s = "SORE" // Asumsi sesi 3 masuk kategori enum MALAM di database
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Pilihan sesi tidak valid. Gunakan 1, 2, atau 3"})
			return
		}

		sesiStr = &s

		// Gabungkan string tanggal dan jam, lalu parse menjadi time.Time
		waktuMulai, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", input.TanggalAwal, jamMulai), loc)
		waktuSelesai, _ = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", input.TanggalAkhir, jamSelesai), loc)

	} else if input.TipeBooking == "HARIAN" {
		// Asumsi Jam Default untuk Booking Harian: Check-in 14:00, Check-out 12:00
		waktuMulai, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 14:00:00", input.TanggalAwal), loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_awal salah. Gunakan YYYY-MM-DD"})
			return
		}

		waktuSelesai, err = time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s 12:00:00", input.TanggalAkhir), loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_akhir salah. Gunakan YYYY-MM-DD"})
			return
		}

		// Hitung murni selisih hari dari kalender
		tAwal, _ := time.Parse("2006-01-02", input.TanggalAwal)
		tAkhir, _ := time.Parse("2006-01-02", input.TanggalAkhir)
		hari := int(tAkhir.Sub(tAwal).Hours() / 24)

		// Validasi agar tanggal akhir tidak lebih mundur dari tanggal awal
		if hari <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal akhir harus lebih besar dari tanggal awal"})
			return
		}

		jumlahHari = &hari

	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tipe_booking tidak valid, gunakan 'SESI' atau 'HARIAN'"})
		return
	}

	// 3. Mapping ke Struct Database (dari models yang dibuat sebelumnya)
	// Catatan: Pastikan CustomerID diisi, misalnya dari payload JWT token milik user yang login.
	newBooking := schemas.Booking{
		// CustomerID:  "...", // TODO: Ambil ID dari user yang sedang login
		TipeBooking:  input.TipeBooking,
		Sesi:         sesiStr,
		JumlahHari:   jumlahHari,
		WaktuMulai:   waktuMulai,
		WaktuSelesai: waktuSelesai,
		Status:       "PENDING",
	}

	// TODO: Simpan ke database menggunakan GORM
	// if err := config.DB.Create(&newBooking).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data booking"})
	// 	return
	// }

	// 4. Return Response Berhasil
	c.JSON(http.StatusOK, gin.H{
		"message": "Data booking berhasil diproses",
		"data":    newBooking,
	})
}
func EditProfilAccount(c *gin.Context) {

	user_id, exist := c.Get("IdUser")
	if !exist {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Unthorized data id"})
		return
	}
	var akun schemas.EditProfilAccountSchema
	if err := c.ShouldBindBodyWithJSON(&akun); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no data get"})
		return
	}
	err := databases.DB.Where("ID = ?", user_id).Updates(schemas.Customer{Nama: akun.Nama, Email: akun.Email, NomerTelpon: akun.NomorTelpon})
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot Update data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Can update data user",
		"data":    err,
	})

}
