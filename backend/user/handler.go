package users

import (
	"net/http"
	schemas "stpNew/backend/schema"

	"github.com/gin-gonic/gin"
)

func loginAlur(c *gin.Context) {
	var person schemas.LoginSchema
	if err := c.ShouldBindQuery(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data":    person,
	})
}
