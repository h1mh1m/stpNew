package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	"stpNew/backend/testutil"
	users "stpNew/backend/user"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Helper to seed a baseline customer and room
func seedCustomerAndRoom(t *testing.T, db *gorm.DB, custID, roomID string) {
	t.Helper()
	cust := schemas.Customer{
		ID:    custID,
		Nama:  "Test Customer " + custID,
		Email: custID + "@example.com",
	}
	if err := db.Create(&cust).Error; err != nil {
		t.Fatalf("failed to seed customer: %v", err)
	}

	room := schemas.Room{
		ID:          roomID,
		Name:        "Sirkuit Otomotif ITS",
		Description: "Fasilitas Pengujian & Balap",
		Capacity:    100,
		Price:       150000,
	}
	if err := db.Create(&room).Error; err != nil {
		t.Fatalf("failed to seed room: %v", err)
	}
}

// Step 10: Database Constraint Tests
func TestDatabaseConstraints(t *testing.T) {
	gormDB, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	seedCustomerAndRoom(t, gormDB, "cust-const", "room-const")
	now := time.Now().UTC().Truncate(time.Second)

	// 1. Valid order is accepted
	ord1 := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(2 * time.Hour),
		DateFinish: now.Add(4 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 300000,
	}
	if err := gormDB.Create(&ord1).Error; err != nil {
		t.Fatalf("valid order should be accepted, got: %v", err)
	}

	// 2. Missing customer is rejected by foreign key (fk_orders_customer)
	ordNoCust := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "nonexistent-cust",
		RoomID:     "room-const",
		DateStart:  now.Add(10 * time.Hour),
		DateFinish: now.Add(12 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordNoCust).Error; err == nil {
		t.Fatal("expected foreign key violation for nonexistent customer, got nil")
	}

	// 3. Missing room is rejected by foreign key (fk_orders_room)
	ordNoRoom := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "nonexistent-room",
		DateStart:  now.Add(10 * time.Hour),
		DateFinish: now.Add(12 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordNoRoom).Error; err == nil {
		t.Fatal("expected foreign key violation for nonexistent room, got nil")
	}

	// 4. finish before start is rejected (chk_orders_date_order)
	ordFinishBefore := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(12 * time.Hour),
		DateFinish: now.Add(10 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordFinishBefore).Error; err == nil {
		t.Fatal("expected check constraint violation for finish before start, got nil")
	}

	// 5. finish equal to start is rejected (chk_orders_date_order)
	ordFinishEqual := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(10 * time.Hour),
		DateFinish: now.Add(10 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordFinishEqual).Error; err == nil {
		t.Fatal("expected check constraint violation for finish equal start, got nil")
	}

	// 6. Unknown status is rejected (chk_orders_status)
	ordBadStatus := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(10 * time.Hour),
		DateFinish: now.Add(12 * time.Hour),
		Status:     "InvalidStatusString",
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordBadStatus).Error; err == nil {
		t.Fatal("expected check constraint violation for invalid status, got nil")
	}

	// 7. Allowed statuses are accepted: Selesai and Dibatalkan
	ordSelesai := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(10 * time.Hour),
		DateFinish: now.Add(12 * time.Hour),
		Status:     schemas.StatusSelesai,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordSelesai).Error; err != nil {
		t.Fatalf("expected status 'Selesai' to be accepted, got: %v", err)
	}

	// 8. Same-room partially overlapping booking is rejected (SQLSTATE 23P01)
	// ord1 is [2h, 4h). Overlap at [3h, 5h).
	ordOverlapPartial := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(3 * time.Hour),
		DateFinish: now.Add(5 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	err := gormDB.Create(&ordOverlapPartial).Error
	if err == nil {
		t.Fatal("expected exclusion constraint violation for partial overlap, got nil")
	}
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23P01" {
		t.Fatalf("expected SQLSTATE 23P01, got: %v", err)
	}

	// 9. Same-room containing interval is rejected ([1h, 5h) contains [2h, 4h))
	ordOverlapContaining := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(1 * time.Hour),
		DateFinish: now.Add(5 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	err = gormDB.Create(&ordOverlapContaining).Error
	if err == nil {
		t.Fatal("expected exclusion constraint violation for containing interval, got nil")
	}
	if !errors.As(err, &pgErr) || pgErr.Code != "23P01" {
		t.Fatalf("expected SQLSTATE 23P01, got: %v", err)
	}

	// 10. Same-room contained interval is rejected ([2.5h, 3.5h) inside [2h, 4h))
	ordOverlapContained := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(2*time.Hour + 30*time.Minute),
		DateFinish: now.Add(3*time.Hour + 30*time.Minute),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	err = gormDB.Create(&ordOverlapContained).Error
	if err == nil {
		t.Fatal("expected exclusion constraint violation for contained interval, got nil")
	}
	if !errors.As(err, &pgErr) || pgErr.Code != "23P01" {
		t.Fatalf("expected SQLSTATE 23P01, got: %v", err)
	}

	// 11. Same-room identical interval is rejected ([2h, 4h))
	ordOverlapIdentical := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(2 * time.Hour),
		DateFinish: now.Add(4 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	err = gormDB.Create(&ordOverlapIdentical).Error
	if err == nil {
		t.Fatal("expected exclusion constraint violation for identical interval, got nil")
	}
	if !errors.As(err, &pgErr) || pgErr.Code != "23P01" {
		t.Fatalf("expected SQLSTATE 23P01, got: %v", err)
	}

	// 12. Same-room adjacent interval is ACCEPTED ([4h, 6h) touches [2h, 4h))
	ordAdjacent := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(4 * time.Hour),
		DateFinish: now.Add(6 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordAdjacent).Error; err != nil {
		t.Fatalf("adjacent booking [4h, 6h) should be accepted, got: %v", err)
	}

	// 13. Different-room overlapping interval is ACCEPTED
	seedCustomerAndRoom(t, gormDB, "cust-2", "room-2")
	ordDiffRoom := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-2",
		RoomID:     "room-2",
		DateStart:  now.Add(2 * time.Hour),
		DateFinish: now.Add(4 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordDiffRoom).Error; err != nil {
		t.Fatalf("different-room overlapping booking should be accepted, got: %v", err)
	}

	// 14. Dibatalkan status DOES NOT block another booking
	ordCancelled := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(14 * time.Hour),
		DateFinish: now.Add(16 * time.Hour),
		Status:     schemas.StatusDibatalkan,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordCancelled).Error; err != nil {
		t.Fatalf("failed to insert Dibatalkan order: %v", err)
	}
	ordOverCancelled := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-const",
		RoomID:     "room-const",
		DateStart:  now.Add(14 * time.Hour),
		DateFinish: now.Add(16 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	if err := gormDB.Create(&ordOverCancelled).Error; err != nil {
		t.Fatalf("order overlapping a Dibatalkan order should be accepted, got: %v", err)
	}
}

// Step 12: Live Concurrency Test (with go test -race)
func TestConcurrentBookings_SameRoomAndInterval(t *testing.T) {
	gormDB, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	seedCustomerAndRoom(t, gormDB, "cust-race", "room-race")
	databases.DB = gormDB

	now := time.Now().UTC()
	startStr := now.Add(24 * time.Hour).Format(time.RFC3339)
	finishStr := now.Add(26 * time.Hour).Format(time.RFC3339)

	// Set Gin test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/Bookings", func(c *gin.Context) {
		c.Set("id", "cust-race")
		users.CreateBooking(c)
	})

	const numAttempts = 10
	startBarrier := make(chan struct{})
	var wg sync.WaitGroup
	var successCount int32
	var conflictCount int32
	var otherCount int32

	for i := 0; i < numAttempts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startBarrier // synchronize start

			body := map[string]string{
				"room_id":     "room-race",
				"date_start":  startStr,
				"date_finish": finishStr,
			}
			jsonBody, _ := json.Marshal(body)
			req, _ := http.NewRequest(http.MethodPost, "/api/Bookings", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			switch w.Code {
			case http.StatusCreated:
				atomic.AddInt32(&successCount, 1)
			case http.StatusConflict:
				atomic.AddInt32(&conflictCount, 1)
			default:
				atomic.AddInt32(&otherCount, 1)
				t.Logf("Unexpected response code %d: %s", w.Code, w.Body.String())
			}
		}()
	}

	// Release all attempts simultaneously
	close(startBarrier)
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected EXACTLY 1 success among 10 concurrent requests, got %d", successCount)
	}
	if conflictCount != 9 {
		t.Fatalf("expected EXACTLY 9 conflicts among 10 concurrent requests, got %d (other=%d)", conflictCount, otherCount)
	}

	// Verify database row count is exactly 1
	var dbCount int64
	gormDB.Model(&schemas.Order{}).Where("room_id = ?", "room-race").Count(&dbCount)
	if dbCount != 1 {
		t.Fatalf("expected exactly 1 row in database for room-race, got %d", dbCount)
	}

	// 10 Adjacent non-overlapping bookings - all 10 must succeed
	var adjacentSuccess int32
	for i := 0; i < 10; i++ {
		adjStart := now.Add(time.Duration(30+i*2) * time.Hour)
		adjFinish := now.Add(time.Duration(30+(i+1)*2) * time.Hour)
		ord := schemas.Order{
			ID:         uuid.NewString(),
			UserID:     "cust-race",
			RoomID:     "room-race",
			DateStart:  adjStart,
			DateFinish: adjFinish,
			Status:     schemas.StatusMenungguPembayaran,
			TotalPrice: 150000,
		}
		if err := gormDB.Create(&ord).Error; err != nil {
			t.Fatalf("adjacent booking %d failed: %v", i, err)
		}
		adjacentSuccess++
	}
	if adjacentSuccess != 10 {
		t.Fatalf("expected 10 adjacent bookings to succeed, got %d", adjacentSuccess)
	}
}

// Step 15: Booking Handler Integration Test
func TestBookingHandler_PostAndGetIntegration(t *testing.T) {
	gormDB, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	seedCustomerAndRoom(t, gormDB, "cust-alpha", "room-alpha")
	seedCustomerAndRoom(t, gormDB, "cust-beta", "room-beta")
	databases.DB = gormDB

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/Bookings", func(c *gin.Context) {
		uid := c.GetHeader("X-User-ID")
		c.Set("id", uid)
		users.CreateBooking(c)
	})
	router.GET("/api/Bookings", func(c *gin.Context) {
		uid := c.GetHeader("X-User-ID")
		c.Set("id", uid)
		users.GetBookings(c)
	})

	now := time.Now().UTC()
	startStr := now.Add(2 * time.Hour).Format(time.RFC3339)
	finishStr := now.Add(4 * time.Hour).Format(time.RFC3339)

	// 1. Valid booking creation - verifies client overrides for user_id, price, and status are IGNORED
	body := map[string]interface{}{
		"room_id":     "room-alpha",
		"date_start":  startStr,
		"date_finish": finishStr,
		"user_id":     "cust-attacker", // must be ignored
		"price":       10,              // must be ignored
		"total_price": 10,              // must be ignored
		"status":      "Selesai",       // must be ignored
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/api/Bookings", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "cust-alpha")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected HTTP 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Message string        `json:"message"`
		Data    schemas.Order `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Assert authoritative values
	if resp.Data.UserID != "cust-alpha" {
		t.Errorf("user_id override attempt succeeded! expected 'cust-alpha', got '%s'", resp.Data.UserID)
	}
	if resp.Data.TotalPrice != 150000 {
		t.Errorf("price override attempt succeeded! expected 150000 (room price), got %f", resp.Data.TotalPrice)
	}
	if resp.Data.Status != schemas.StatusMenungguPembayaran {
		t.Errorf("status override attempt succeeded! expected 'Menunggu Pembayaran', got '%s'", resp.Data.Status)
	}

	// 2. Missing room returns HTTP 404
	badRoomBody, _ := json.Marshal(map[string]string{
		"room_id":     "nonexistent-room",
		"date_start":  now.Add(10 * time.Hour).Format(time.RFC3339),
		"date_finish": now.Add(12 * time.Hour).Format(time.RFC3339),
	})
	req404, _ := http.NewRequest(http.MethodPost, "/api/Bookings", bytes.NewReader(badRoomBody))
	req404.Header.Set("Content-Type", "application/json")
	req404.Header.Set("X-User-ID", "cust-alpha")
	w404 := httptest.NewRecorder()
	router.ServeHTTP(w404, req404)
	if w404.Code != http.StatusNotFound {
		t.Fatalf("expected HTTP 404 for missing room, got %d", w404.Code)
	}

	// 3. Existing overlap returns HTTP 409
	req409, _ := http.NewRequest(http.MethodPost, "/api/Bookings", bytes.NewReader(bodyBytes))
	req409.Header.Set("Content-Type", "application/json")
	req409.Header.Set("X-User-ID", "cust-alpha")
	w409 := httptest.NewRecorder()
	router.ServeHTTP(w409, req409)
	if w409.Code != http.StatusConflict {
		t.Fatalf("expected HTTP 409 for overlapping booking, got %d", w409.Code)
	}

	// 4. Test GET /api/Bookings ownership scoping
	// Insert order for cust-beta
	ordBeta := schemas.Order{
		ID:         uuid.NewString(),
		UserID:     "cust-beta",
		RoomID:     "room-beta",
		DateStart:  now.Add(20 * time.Hour),
		DateFinish: now.Add(22 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	gormDB.Create(&ordBeta)

	// Fetch bookings for cust-alpha
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/Bookings", nil)
	reqGet.Header.Set("X-User-ID", "cust-alpha")
	wGet := httptest.NewRecorder()
	router.ServeHTTP(wGet, reqGet)

	if wGet.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200 on GET /api/Bookings, got %d", wGet.Code)
	}

	var getResp struct {
		Message string          `json:"message"`
		Data    []schemas.Order `json:"data"`
	}
	if err := json.Unmarshal(wGet.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to parse GET bookings response: %v", err)
	}

	if len(getResp.Data) != 1 {
		t.Fatalf("expected exactly 1 order for cust-alpha, got %d", len(getResp.Data))
	}
	if getResp.Data[0].UserID != "cust-alpha" {
		t.Errorf("expected booking to belong to cust-alpha, got %s", getResp.Data[0].UserID)
	}
}
