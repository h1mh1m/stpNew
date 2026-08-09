package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func LoginAlur(c *gin.Context) {
	var person schemas.LoginSchema
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data":    person,
	})
}

func RootWeb(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "halo"})
	// c.JSON(200, "/api/login")
}

func SignupAlur(c *gin.Context) {
	var person schemas.SignupSchema
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	Bytes, err := bcrypt.GenerateFromPassword([]byte(person.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi password"})
		return
	}
	var hashedPassword = string(Bytes)
	fmt.Sprintln("dbhbcqej %s", hashedPassword)

	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatalln("Error load .env file")
	}

	ResendAPI := os.Getenv("RESEND_API")
	if ResendAPI == "" {
		fmt.Sprintln("Cannot get api key from env file")
	}
	token := os.Getenv("SECRET")
	if token == "" {
		log.Fatalln("Cannot get Token from env file")
	}
	client := resend.NewClient(ResendAPI)
	verifyLink := fmt.Sprintf("http://localhost:8080/verify?token=%s", token)
	params := &resend.SendEmailRequest{
		From:    "Onboarding <onboarding@resend.dev>", // Domain default bawaan Resend untuk testing
		To:      []string{string(person.Email)},
		Subject: "Verifikasi Email Dev",
		Html:    fmt.Sprintf("<p>Klik link berikut untuk verifikasi: <a href='%s'>Verifikasi</a></p>", verifyLink),
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		log.Fatalf("Gagal kirim email: %v", err)
	}
	fmt.Sprintln("hbjdcbash %s", sent)
	c.JSON(http.StatusOK, gin.H{
		"message": "Anda sudah terdaftar",
	})
}

func ViewPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You're in View Page"})
}

func DashboardPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "You're in Dashboard Page"})
}

// var srv *calendar.Service

// func SchedulePage(c *gin.Context) {

// 	calendarId := "raushanfekra@gmail.com"
// 	t := time.Now().Format(time.RFC3339)
// 	events, err := srv.Events.List(calendarId).
// 		ShowDeleted(false).
// 		SingleEvents(true).
// 		TimeMin(t).
// 		MaxResults(10). // Ambil 10 jadwal terdekat
// 		OrderBy("startTime").
// 		Do()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Gagal menyambung ke google kalender",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "You are in Schedule Page",
// 		"data":    events.Items,
// 	})
// }

func SchedulePage(c *gin.Context) {
	ctx := context.Background()

	// 1. Inisialisasi Google Calendar Client langsung di dalam fungsi
	// Pastikan path ke credentials.json benar
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

	calendarId := "raushanfekra@gmail.com"
	t := time.Now().Format(time.RFC3339)

	// 2. Ambil data event dari Google Calendar
	events, err := srv.Events.List(calendarId).
		ShowDeleted(false).
		SingleEvents(true).
		TimeMin(t).
		MaxResults(10). // Ambil 10 jadwal terdekat
		OrderBy("startTime").
		Do()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyambung ke google kalender",
			"error":   err.Error(),
		})
		return
	}

	// 3. Kirim respon JSON
	c.JSON(http.StatusOK, gin.H{
		"message": "You are in Schedule Page",
		"data":    events.Items,
	})
}

func GetDataBooking(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "haveuckuqe"})
}

func main() {
	databases.ConnectDB()
	route := gin.Default()
	route.GET("/", RootWeb)
	route.GET("/api/Schedule", SchedulePage)
	route.GET("/api/View", ViewPage)
	route.POST("/api/login", LoginAlur)
	route.GET("/api/login", GetDataBooking)
	route.POST("/api/signup", SignupAlur)
	route.Run(":3000")

}
