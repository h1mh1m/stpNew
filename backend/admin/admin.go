package admin

import (
	"errors"
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

func getAdminJwtSecret() ([]byte, error) {
	secret := strings.TrimSpace(os.Getenv("JWTADMIN"))
	if secret == "" {
		return nil, errors.New("JWTADMIN secret is empty or not configured")
	}
	if len(secret) < 8 {
		return nil, errors.New("JWTADMIN secret is unsafe (must be at least 8 characters)")
	}
	return []byte(secret), nil
}

func GenerateTokenAdmin(nama string) (string, error) {
	if strings.TrimSpace(nama) == "" {
		return "", errors.New("admin ID/name cannot be empty")
	}

	secret, err := getAdminJwtSecret()
	if err != nil {
		return "", err
	}

	expiratedTime := time.Now().Add(24 * time.Hour)
	claimsAdmin := &schemas.Claims{
		IdUser: nama,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiratedTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tokenAdmin := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAdmin)
	return tokenAdmin.SignedString(secret)
}
