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
		verifyLink := fmt.Sprintf("http://localhost:3000/api/verify?token=%s&email=%s", secretToken, person.Email)
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

func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")
	secretToken := os.Getenv("SECRET")

	if token == "" || email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token atau email tidak valid"})
		return
	}

	if token != secretToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token verifikasi salah"})
		return
	}

	var user schemas.Customer
	if err := databases.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	if user.Status == "Verified" {
		c.JSON(http.StatusOK, gin.H{"message": "Email sudah terverifikasi sebelumnya"})
		return
	}

	user.Status = "Verified"
	if err := databases.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan status verifikasi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email berhasil diverifikasi! Anda sekarang dapat login."})
}

func DashboardPage(c *gin.Context) {
	// Mock rooms data since database isn't fully seeded yet
	rooms := []schemas.Room{
		{
			ID:          "1",
			Name:        "Lab Pemrograman 1",
			Description: "Departemen Teknik Informatika ITS - Zona A",
			Capacity:    25,
			Price:       90000,
			ImageURL:    "",
			Features:    "Jangka Pendek",
		},
		{
			ID:          "2",
			Name:        "Lab Pemrograman 2",
			Description: "Departemen Teknik Informatika ITS - Zona A",
			Capacity:    25,
			Price:       100000,
			ImageURL:    "",
			Features:    "Jangka Pendek",
		},
		{
			ID:          "3",
			Name:        "Gedung NASDEC",
			Description: "Institut Teknologi Sepuluh Nopember - Zona C",
			Capacity:    100,
			Price:       5000000,
			ImageURL:    "",
			Features:    "Jangka Panjang",
		},
		{
			ID:          "4",
			Name:        "Menara Sains",
			Description: "Institut Teknologi Sepuluh Nopember - Zona C",
			Capacity:    50,
			Price:       10000000,
			ImageURL:    "",
			Features:    "Jangka Panjang",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    rooms,
	})
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

func GetBookings(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var orders []schemas.Order
	if err := databases.DB.Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    orders,
	})
}

func CreateBooking(c *gin.Context) {
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req schemas.BookingSchema
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Schedule Conflict Validation
	// Overlap occurs if: new_start < existing_finish AND new_finish > existing_start
	var existingOrder schemas.Order
	conflictErr := databases.DB.Where("room_id = ? AND status IN (?, ?) AND date_start < ? AND date_finish > ?", 
		req.RoomID, "Selesai", "Menunggu Pembayaran", req.DateFinish, req.DateStart).First(&existingOrder).Error
	
	if conflictErr == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Room is already booked for the selected schedule. Schedule Conflict!"})
		return
	} else if !errors.Is(conflictErr, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking schedule conflicts"})
		return
	}

	// If no conflict, mock finding the room for price (using a dummy price for now)
	var room schemas.Room
	databases.DB.Where("id = ?", req.RoomID).First(&room)
	price := room.Price
	if price == 0 {
		price = 100000 // default dummy price if not seeded
	}

	orderId := uuid.NewString()
	newOrder := schemas.Order{
		ID:         orderId,
		UserID:     userId.(string),
		RoomID:     req.RoomID,
		DateStart:  req.DateStart,
		DateFinish: req.DateFinish,
		Status:     "Menunggu Pembayaran",
		TotalPrice: price,
	}

	if err := databases.DB.Create(&newOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// System Payment Mock
	paymentLink := fmt.Sprintf("http://localhost:3000/api/payments/pay?order_id=%s", orderId)

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking created successfully",
		"payment_link": paymentLink,
		"data": newOrder,
	})
}

func MockPaymentWebhook(c *gin.Context) {
	orderId := c.Query("order_id")
	if orderId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	var order schemas.Order
	if err := databases.DB.Where("id = ?", orderId).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if order.Status == "Selesai" {
		c.JSON(http.StatusOK, gin.H{"message": "Order is already paid"})
		return
	}

	order.Status = "Selesai"
	if err := databases.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment successful. Order is now Selesai.",
		"order": order,
	})
}
