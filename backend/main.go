package main

import (
	"log"
	"net/http"
	admin "stpNew/backend/admin"
	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	users "stpNew/backend/user"
	"os"

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
		models := []any{
			&schemas.Customer{},
			&schemas.Admin{},
			&schemas.Room{},
			&schemas.Order{},
		}
		for _, m := range models {
			if err := databases.DB.AutoMigrate(m); err != nil {
				log.Printf("Notice: Auto Migration note for %T: %v\n", m, err)
			}
		}
		log.Println("Database AutoMigration complete")

		var roomCount int64
		databases.DB.Model(&schemas.Room{}).Count(&roomCount)
		if roomCount == 0 {
			sampleRooms := []schemas.Room{
				{
					ID:          "1",
					Name:        "Lab Pemrograman 1",
					Description: "Departemen Teknik Informatika ITS - Zona A",
					Capacity:    25,
					Price:       90000,
					ImageURL:    "https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?auto=format&fit=crop&q=80&w=1000",
					Features:    "Jangka Pendek",
				},
				{
					ID:          "2",
					Name:        "Lab Pemrograman 2",
					Description: "Departemen Teknik Informatika ITS - Zona A",
					Capacity:    25,
					Price:       100000,
					ImageURL:    "https://images.unsplash.com/photo-1522071820081-009f0129c71c?auto=format&fit=crop&q=80&w=1000",
					Features:    "Jangka Pendek",
				},
				{
					ID:          "3",
					Name:        "Gedung NASDEC",
					Description: "Institut Teknologi Sepuluh Nopember - Zona C",
					Capacity:    100,
					Price:       5000000,
					ImageURL:    "https://images.unsplash.com/photo-1497366216548-37526070297c?auto=format&fit=crop&q=80&w=1000",
					Features:    "Jangka Panjang",
				},
				{
					ID:          "4",
					Name:        "Menara Sains",
					Description: "Institut Teknologi Sepuluh Nopember - Zona C",
					Capacity:    50,
					Price:       10000000,
					ImageURL:    "https://images.unsplash.com/photo-1486406146926-c627a92ad1ab?auto=format&fit=crop&q=80&w=1000",
					Features:    "Jangka Panjang",
				},
			}
			for _, r := range sampleRooms {
				_ = databases.DB.Create(&r)
			}
			log.Println("Database room seeding complete")
		}
	}

	// 4. Setup Gin Router
	route := gin.Default()

	// Enable CORS Middleware for frontend web integration
	route.Use(corsMiddleware())

	// Health check / root info
	route.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "STP Backend API is running",
		})
	})

	// Public routes
	registerPublic := func(rg *gin.RouterGroup) {
		rg.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		rg.POST("/login", users.LoginAlur)
		rg.POST("/signup", users.SignupAlur)
		rg.GET("/verify", users.VerifyEmail)
		rg.GET("/payments/pay", users.MockPaymentWebhook)
		rg.POST("/admin/login", admin.AdminLogin)
		rg.GET("/Beranda", users.DashboardPage)
		rg.GET("/beranda", users.DashboardPage)
		rg.GET("/rooms", users.DashboardPage)
	}

	registerPublic(route.Group(""))
	registerPublic(route.Group("/api"))

	// Authorized routes
	registerAuthorized := func(rg *gin.RouterGroup) {
		rg.GET("/View", users.ViewPage)
		rg.GET("/view", users.ViewPage)
		rg.GET("/Schedule", users.SchedulePage)
		rg.GET("/schedule", users.SchedulePage)
		rg.GET("/Profile", users.ViewPage)
		rg.GET("/profile", users.ViewPage)
		rg.GET("/Bookings", users.GetBookings)
		rg.GET("/booking", users.GetBookings)
		rg.GET("/bookings", users.GetBookings)
		rg.POST("/Bookings", users.CreateBooking)
		rg.POST("/booking", users.CreateBooking)
		rg.POST("/bookings", users.CreateBooking)
		rg.PATCH("/Profile/Change-Password", users.ViewPage)
		rg.PATCH("/profile/Change-Password", users.ViewPage)
		rg.PATCH("/profile/change-password", users.ViewPage)
	}

	appAuth := route.Group("/app")
	appAuth.Use(users.AuthMiddlewareJwt())
	registerAuthorized(appAuth)

	apiAuth := route.Group("/api")
	apiAuth.Use(users.AuthMiddlewareJwt())
	registerAuthorized(apiAuth)

	// 5. Read server port from environment or default to 8000
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}
	route.Run(":" + port)
}

// corsMiddleware sets headers to allow cross-origin requests from web frontends
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
