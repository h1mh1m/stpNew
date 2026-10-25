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
	"github.com/jackc/pgx/v5/pgconn"
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

	if user.Status != "Verified" {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Email belum diverifikasi",
			"error":   "Email belum diverifikasi",
		})
		return
	}

	tokenStr, err := GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot generate jwt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   tokenStr,
		// user is retained as a deprecated compatibility field for legacy frontend callers
		"user": tokenStr,
	})
}

func RootWeb(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login")
}

func getJwtSecret() ([]byte, error) {
	secret := strings.TrimSpace(os.Getenv("JWT"))
	if secret == "" {
		return nil, errors.New("JWT secret is empty or not configured")
	}
	if len(secret) < 8 {
		return nil, errors.New("JWT secret is unsafe (must be at least 8 characters)")
	}
	return []byte(secret), nil
}

func GenerateToken(userID string) (string, error) {
	if strings.TrimSpace(userID) == "" {
		return "", errors.New("customer ID cannot be empty")
	}

	secret, err := getJwtSecret()
	if err != nil {
		return "", err
	}

	expirateTime := time.Now().Add(2 * time.Hour)

	claim := &schemas.Claims{
		IdUser: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirateTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(secret)
}

func AuthMiddlewareJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" || strings.TrimSpace(tokenParts[1]) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing Bearer token format"})
			c.Abort()
			return
		}
		tokenString := strings.TrimSpace(tokenParts[1])

		secret, err := getJwtSecret()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication unavailable"})
			c.Abort()
			return
		}

		claims := &schemas.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method == nil || token.Method.Alg() != "HS256" || token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing algorithm: %v", token.Header["alg"])
			}
			return secret, nil
		}, jwt.WithExpirationRequired())

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
			c.Abort()
			return
		}

		if strings.TrimSpace(claims.IdUser) == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing customer ID claim"})
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
	if databases.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not available"})
		return
	}

	var rooms []schemas.Room
	if err := databases.DB.Find(&rooms).Error; err != nil {
		log.Printf("Error querying room catalogue: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch room catalogue"})
		return
	}

	if rooms == nil {
		rooms = []schemas.Room{}
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

	if databases.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not available"})
		return
	}

	var orders []schemas.Order
	if err := databases.DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&orders).Error; err != nil {
		log.Printf("Error fetching orders for user %v: %v\n", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	if orders == nil {
		orders = []schemas.Order{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    orders,
	})
}

// GetRoomAvailability returns occupied blocking intervals for a specific room.
// Only schedules with BlockingStatuses (Menunggu Pembayaran, Selesai) are returned.
// Other customers' user IDs, names, emails, and order prices are strictly omitted.
func GetRoomAvailability(c *gin.Context) {
	_, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	roomID := strings.TrimSpace(c.Param("room_id"))
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id wajib diisi"})
		return
	}

	if databases.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service unavailable"})
		return
	}

	// 1. Verify room exists
	var room schemas.Room
	if err := databases.DB.Where("id = ?", roomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ruangan tidak ditemukan"})
			return
		}
		log.Printf("Error verifying room %s: %v\n", roomID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi ruangan"})
		return
	}

	// 2. Parse optional bounded range
	// Default start: 5 minutes ago to include current active slots
	now := time.Now().UTC()
	rangeStart := now.Add(-5 * time.Minute)
	// Default finish: 90 days from now
	rangeFinish := now.AddDate(0, 3, 0)

	if startParam := strings.TrimSpace(c.Query("start")); startParam != "" {
		parsedStart, err := time.Parse(time.RFC3339, startParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format parameter start tidak valid: harus RFC 3339"})
			return
		}
		rangeStart = parsedStart.UTC()
	}

	if finishParam := strings.TrimSpace(c.Query("finish")); finishParam != "" {
		parsedFinish, err := time.Parse(time.RFC3339, finishParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format parameter finish tidak valid: harus RFC 3339"})
			return
		}
		rangeFinish = parsedFinish.UTC()
	}

	if rangeFinish.Before(rangeStart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parameter finish harus lebih lambat dari atau sama dengan start"})
		return
	}

	// Bounded limit: prevent querying more than 1 year span
	if rangeFinish.Sub(rangeStart) > 366*24*time.Hour {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rentang waktu ketersediaan maksimal 1 tahun"})
		return
	}

	// 3. Query blocking orders within range
	var slots []schemas.AvailabilitySlot
	if err := databases.DB.Model(&schemas.Order{}).
		Select("room_id, date_start, date_finish, status").
		Where("room_id = ? AND status IN (?) AND date_finish >= ? AND date_start <= ?",
			room.ID, schemas.BlockingStatuses, rangeStart, rangeFinish).
		Order("date_start ASC").
		Scan(&slots).Error; err != nil {
		log.Printf("Error fetching availability for room %s: %v\n", room.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ketersediaan ruangan"})
		return
	}

	if slots == nil {
		slots = []schemas.AvailabilitySlot{}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    slots,
	})
}

// ValidateBookingInterval validates and parses RFC 3339 timestamps into UTC time.Time.
// It enforces finish > start and rejects past start times (allowing 5-minute clock tolerance).
func ValidateBookingInterval(startStr, finishStr string, now time.Time) (time.Time, time.Time, error) {
	if strings.TrimSpace(startStr) == "" || strings.TrimSpace(finishStr) == "" {
		return time.Time{}, time.Time{}, errors.New("date_start dan date_finish wajib diisi")
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("format date_start tidak valid: harus RFC 3339 (misal: 2026-11-01T08:00:00Z)")
	}

	finish, err := time.Parse(time.RFC3339, finishStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("format date_finish tidak valid: harus RFC 3339 (misal: 2026-11-01T10:00:00Z)")
	}

	startUTC := start.UTC()
	finishUTC := finish.UTC()

	if !finishUTC.After(startUTC) {
		return time.Time{}, time.Time{}, errors.New("waktu selesai harus lebih lambat dari waktu mulai")
	}

	// Clock tolerance of 5 minutes for slight server-client clock drift
	if startUTC.Before(now.UTC().Add(-5 * time.Minute)) {
		return time.Time{}, time.Time{}, errors.New("waktu mulai tidak boleh di masa lampau")
	}

	return startUTC, finishUTC, nil
}

// HasIntervalOverlap checks if [start1, finish1) overlaps with [start2, finish2).
// Semantics: start is inclusive, finish is exclusive.
func HasIntervalOverlap(start1, finish1, start2, finish2 time.Time) bool {
	return start1.Before(finish2) && finish1.After(start2)
}

// IsExclusionConstraintViolation checks if an error represents a PostgreSQL exclusion violation (SQLSTATE 23P01).
// Stable typed error classification only (errors.As on *pgconn.PgError).
func IsExclusionConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23P01"
	}
	return false
}

func CreateBooking(c *gin.Context) {
	userIdVal, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userId, ok := userIdVal.(string)
	if !ok || strings.TrimSpace(userId) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: user ID tidak valid"})
		return
	}

	var req schemas.BookingSchema
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid atau field wajib kosong"})
		return
	}

	if strings.TrimSpace(req.RoomID) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room_id wajib diisi"})
		return
	}

	startUTC, finishUTC, err := ValidateBookingInterval(req.DateStart, req.DateFinish, time.Now())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if databases.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database service unavailable"})
		return
	}

	ctx := c.Request.Context()
	tx := databases.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Printf("Failed to begin booking transaction: %v\n", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi reservasi"})
		return
	}
	defer tx.Rollback()

	// 1. Authoritative Room Lookup: room must exist in database
	var room schemas.Room
	if err := tx.Where("id = ?", req.RoomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ruangan tidak ditemukan"})
			return
		}
		log.Printf("Error querying room %s: %v\n", req.RoomID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memverifikasi ruangan"})
		return
	}

	// 2. Preliminary friendly overlap check
	var conflictCount int64
	if err := tx.Model(&schemas.Order{}).
		Where("room_id = ? AND status IN (?) AND date_start < ? AND date_finish > ?",
			room.ID, schemas.BlockingStatuses, finishUTC, startUTC).
		Count(&conflictCount).Error; err != nil {
		log.Printf("Error checking booking overlap: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memeriksa ketersediaan jadwal"})
		return
	}
	if conflictCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Jadwal bentrok: Ruangan sudah dipesan untuk jadwal yang dipilih"})
		return
	}

	// 3. Create Order with authoritative price and server-side status
	orderId := uuid.NewString()
	newOrder := schemas.Order{
		ID:         orderId,
		UserID:     userId,
		RoomID:     room.ID,
		DateStart:  startUTC,
		DateFinish: finishUTC,
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: room.Price,
	}

	if err := tx.Create(&newOrder).Error; err != nil {
		if IsExclusionConstraintViolation(err) {
			log.Printf("Exclusion constraint caught concurrent booking conflict for room %s\n", room.ID)
			c.JSON(http.StatusConflict, gin.H{"error": "Jadwal bentrok: Ruangan sudah dipesan untuk jadwal yang dipilih"})
			return
		}
		log.Printf("Failed to create order: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat reservasi"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		if IsExclusionConstraintViolation(err) {
			c.JSON(http.StatusConflict, gin.H{"error": "Jadwal bentrok: Ruangan sudah dipesan untuk jadwal yang dipilih"})
			return
		}
		log.Printf("Failed to commit booking transaction: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan reservasi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Booking created successfully",
		"data":    newOrder,
	})
}

// TODO: Replace MockPaymentWebhook with a signed POST payment provider webhook or administrator-controlled simulation endpoint.
// SECURITY NOTICE: This handler is unmounted from routing in backend/main.go during Phase 1 containment.
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
		"order":   order,
	})
}
