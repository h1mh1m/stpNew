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
	c.JSON(http.StatusOK, gin.H{"message": "GetDataBooking response"})
}
