package admin

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AdminLogin(c *gin.Context) {
	var Admin schemas.LoginAdminSchema
	if err := c.ShouldBindJSON(&Admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var AdminDb schemas.Admin
	err := databases.DB.Where("Email = ? AND Nama = ?", Admin.Email, Admin.Nama).First(&AdminDb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadGateway, gin.H{"message": "eror data tidak ada di database"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot get data from database"})
		return
	}
	if Pass := bcrypt.CompareHashAndPassword([]byte(AdminDb.Password), []byte(Admin.Password)); Pass != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Your Password is not Correct"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "You are Logged in"})

}

var JwtAdminSign string = os.Getenv("JWTADMIN")

func GenerateTokenAdmin(nama string) (string, error) {
	expiratedTime := time.Now().Add(24 * time.Hour)
	claimsAdmin := &schemas.Claims{
		IdUser: nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiratedTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tokenAdmin := jwt.NewWithClaims(jwt.SigningMethodES256, claimsAdmin)
	return tokenAdmin.SignedString(JwtAdminSign)
}

func ValidateTokenJwt() gin.HandlerFunc {
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
			return JwtAdminSign, nil
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

func GetBookingsThisCalendarWeek(db *gorm.DB) ([]schemas.Booking, error) {
	var bookings []schemas.Booking

	now := time.Now()

	// Mencari hari Senin pada minggu ini
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Mengubah hari Minggu (0) menjadi 7
	}

	// Set ke hari Senin jam 00:00:00
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())

	// Query GORM antara Senin awal minggu sampai sekarang
	err := db.Where("created_at BETWEEN ? AND ?", startOfWeek, now).Find(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}
