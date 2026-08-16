package main

import (
	"log"
	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	users "stpNew/backend/user"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables (.env or ../.env)
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Println("Notice: .env file not found, using system environment variables")
	}

	// 2. Initialize Database Connection
	databases.ConnectDB()

	// 3. Run GORM AutoMigration for database models
	if databases.DB != nil {
		if err := databases.DB.AutoMigrate(&schemas.Customer{}, &schemas.Admin{}); err != nil {
			log.Printf("Warning: Auto Migration failed: %v\n", err)
		} else {
			log.Println("Database AutoMigration complete")
		}
	}

	// 4. Setup Gin Router
	route := gin.Default()

	// Enable CORS Middleware for frontend web integration
	// route.Use(corsMiddleware())

	authrized := route.Group("/")
	authrized.Use(users.AuthMiddlewareJwt())
	{
		route.GET("/Beranda")
		route.GET("/View")
		route.GET("/Schedule")
		route.GET("/Profile")
		route.GET("/Booking")
		route.POST("/Booking")
		route.PATCH("/Profile/Change-Password")
	}

	// 5. Read server port from environment or default to 3000
	route.Run("3000")
}

// corsMiddleware sets headers to allow cross-origin requests from web frontends
// func corsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 			return
// 		}

// 		c.Next()
// 	}
// }
