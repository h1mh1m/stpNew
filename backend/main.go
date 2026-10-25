package main

import (
	"log"
	"net/http"
	"os"
	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	users "stpNew/backend/user"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables (.env or ../.env)
	if err := godotenv.Load(".env", "../.env"); err != nil {
		log.Println("Notice: .env file not found, using system environment variables")
	}

	// 2. Initialize Database Connection
	if err := databases.ConnectDB(); err != nil || databases.DB == nil {
		log.Fatalf("Fatal: Database initialization failed: %v. Exiting.", err)
	}

	// 3. Schema Management Policy: Explicit migrations manage schema changes.
	// AutoMigrate is disabled by default to avoid uncontrolled runtime schema mutation.
	if strings.ToLower(strings.TrimSpace(os.Getenv("AUTO_MIGRATE"))) == "true" {
		log.Println("Notice: AUTO_MIGRATE is enabled. Running GORM AutoMigrate...")
		if err := databases.DB.AutoMigrate(&schemas.Customer{}, &schemas.Admin{}, &schemas.Room{}, &schemas.Order{}); err != nil {
			log.Printf("Warning: Auto Migration failed: %v\n", err)
		} else {
			log.Println("Database AutoMigration complete")
		}
	} else {
		log.Println("Notice: GORM AutoMigrate is disabled by default. Manage schema via 'make migrate-up' or set AUTO_MIGRATE=true.")
	}

	// 4. Setup Gin Router
	route := gin.Default()

	// Enable CORS Middleware for frontend web integration
	route.Use(corsMiddleware())

	// Public routes
	route.POST("/api/login", users.LoginAlur)
	route.POST("/api/signup", users.SignupAlur)
	route.GET("/api/verify", users.VerifyEmail)

	// SECURITY CONTAINMENT: Route disabled to prevent unauthorized order completion.
	// TODO: Replace with a signed POST payment provider webhook or administrator-controlled simulation endpoint.
	// route.GET("/api/payments/pay", users.MockPaymentWebhook)

	authrized := route.Group("/api")
	authrized.Use(users.AuthMiddlewareJwt())
	{
		authrized.GET("/Beranda", users.DashboardPage)
		authrized.GET("/View", users.ViewPage)
		authrized.GET("/Schedule", users.SchedulePage)
		authrized.GET("/Profile", users.ViewPage) // placeholder
		authrized.GET("/Bookings", users.GetBookings)
		authrized.POST("/Bookings", users.CreateBooking)
		authrized.GET("/rooms/:room_id/availability", users.GetRoomAvailability)
		authrized.PATCH("/Profile/Change-Password", users.ViewPage) // placeholder
	}

	// 5. Read HTTP server port from APP_PORT (separate from database PORT)
	appPort := strings.TrimSpace(os.Getenv("APP_PORT"))
	if appPort == "" {
		appPort = "8080"
	}
	cleanPort := strings.TrimPrefix(appPort, ":")
	portNum, err := strconv.Atoi(cleanPort)
	if err != nil || portNum < 1 || portNum > 65535 {
		log.Printf("Warning: Invalid APP_PORT '%s', defaulting to 8080\n", appPort)
		cleanPort = "8080"
	}

	listenAddr := ":" + cleanPort
	log.Printf("Starting HTTP server on %s\n", listenAddr)
	if err := route.Run(listenAddr); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

// corsMiddleware sets headers to allow cross-origin requests from the configured frontend origin
func corsMiddleware() gin.HandlerFunc {
	frontendOrigin := strings.TrimSpace(os.Getenv("FRONTEND_ORIGIN"))
	if frontendOrigin == "" {
		frontendOrigin = "http://localhost:3000"
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == frontendOrigin || origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", frontendOrigin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
