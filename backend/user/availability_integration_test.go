package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	"stpNew/backend/testutil"
	users "stpNew/backend/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func TestGetRoomAvailability_Integration(t *testing.T) {
	gormDB, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	databases.DB = gormDB

	// Seed customer
	customer1 := schemas.Customer{
		ID:          "cust-avail-1",
		Nama:        "Customer One",
		Email:       "cust1@example.com",
		NomerTelpon: "081234567890",
		Status:      "Verified",
	}
	customer2 := schemas.Customer{
		ID:          "cust-avail-2",
		Nama:        "Customer Two",
		Email:       "cust2@example.com",
		NomerTelpon: "081234567891",
		Status:      "Verified",
	}
	gormDB.Create(&customer1)
	gormDB.Create(&customer2)

	// Seed 2 rooms
	room1 := schemas.Room{
		ID:          "room-target-1",
		Name:        "Sirkuit Otomotif ITS",
		Description: "Target Room",
		Capacity:    100,
		Price:       150000,
	}
	room2 := schemas.Room{
		ID:          "room-other-2",
		Name:        "Ruang Rapat",
		Description: "Other Room",
		Capacity:    20,
		Price:       50000,
	}
	gormDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&room1)
	gormDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&room2)

	now := time.Now().UTC()

	// Seed orders for room 1:
	// Order A: Tomorrow 09:00 - 11:00 (Menunggu Pembayaran) - SHOULD BE INCLUDED
	orderA := schemas.Order{
		ID:         "order-a",
		UserID:     customer1.ID,
		RoomID:     room1.ID,
		DateStart:  now.Add(24 * time.Hour),
		DateFinish: now.Add(26 * time.Hour),
		Status:     schemas.StatusMenungguPembayaran,
		TotalPrice: 150000,
	}
	// Order B: Day after tomorrow 14:00 - 16:00 (Selesai) - SHOULD BE INCLUDED
	orderB := schemas.Order{
		ID:         "order-b",
		UserID:     customer2.ID,
		RoomID:     room1.ID,
		DateStart:  now.Add(48 * time.Hour),
		DateFinish: now.Add(50 * time.Hour),
		Status:     schemas.StatusSelesai,
		TotalPrice: 150000,
	}
	// Order C: Tomorrow 12:00 - 13:00 (Dibatalkan) - MUST BE EXCLUDED
	orderC := schemas.Order{
		ID:         "order-c",
		UserID:     customer1.ID,
		RoomID:     room1.ID,
		DateStart:  now.Add(27 * time.Hour),
		DateFinish: now.Add(28 * time.Hour),
		Status:     schemas.StatusDibatalkan,
		TotalPrice: 150000,
	}
	// Order D: Room 2 Tomorrow (Selesai) - MUST BE EXCLUDED (different room)
	orderD := schemas.Order{
		ID:         "order-d",
		UserID:     customer1.ID,
		RoomID:     room2.ID,
		DateStart:  now.Add(24 * time.Hour),
		DateFinish: now.Add(26 * time.Hour),
		Status:     schemas.StatusSelesai,
		TotalPrice: 50000,
	}
	// Order E: Past order 10 days ago (Selesai) - MUST BE EXCLUDED from default future query
	orderE := schemas.Order{
		ID:         "order-e",
		UserID:     customer1.ID,
		RoomID:     room1.ID,
		DateStart:  now.Add(-240 * time.Hour),
		DateFinish: now.Add(-238 * time.Hour),
		Status:     schemas.StatusSelesai,
		TotalPrice: 150000,
	}

	gormDB.Create(&orderA)
	gormDB.Create(&orderB)
	gormDB.Create(&orderC)
	gormDB.Create(&orderD)
	gormDB.Create(&orderE)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Set("id", "calling-customer-id")
		c.Next()
	})
	router.GET("/api/rooms/:room_id/availability", users.GetRoomAvailability)

	// 1. Missing room returns HTTP 404
	{
		req, _ := http.NewRequest(http.MethodGet, "/api/rooms/non-existent-room/availability", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Fatalf("expected HTTP 404 for missing room, got %d (%s)", w.Code, w.Body.String())
		}
	}

	// 2. Empty availability returns data: []
	{
		// room2 has an order, but let's query for a time window where room2 has no orders
		farFutureStart := now.Add(1000 * time.Hour).Format(time.RFC3339)
		farFutureFinish := now.Add(1002 * time.Hour).Format(time.RFC3339)
		req, _ := http.NewRequest(http.MethodGet, "/api/rooms/"+room2.ID+"/availability?start="+farFutureStart+"&finish="+farFutureFinish, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected HTTP 200, got %d", w.Code)
		}

		var resp struct {
			Message string                     `json:"message"`
			Data    []schemas.AvailabilitySlot `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse json: %v", err)
		}
		if resp.Data == nil || len(resp.Data) != 0 {
			t.Fatalf("expected empty non-null array [], got %v", resp.Data)
		}
	}

	// 3. Room 1 query: verifies filtering, sorting, privacy, and status
	{
		req, _ := http.NewRequest(http.MethodGet, "/api/rooms/"+room1.ID+"/availability", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected HTTP 200, got %d", w.Code)
		}

		// Inspect raw JSON for privacy check: user_id, customer names, total_price MUST NOT exist
		var rawResp struct {
			Message string                   `json:"message"`
			Data    []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &rawResp); err != nil {
			t.Fatalf("failed to parse raw json: %v", err)
		}

		// Exactly 2 active future blocking orders for room 1: orderA and orderB
		if len(rawResp.Data) != 2 {
			t.Fatalf("expected 2 active blocking slots for room 1, got %d (%s)", len(rawResp.Data), w.Body.String())
		}

		// Check sorting: orderA (now + 24h) comes before orderB (now + 48h)
		for i, slot := range rawResp.Data {
			// Room ID must match
			if slot["room_id"] != room1.ID {
				t.Errorf("slot %d: expected room_id %s, got %v", i, room1.ID, slot["room_id"])
			}
			// Privacy check: forbidden keys
			for _, forbidden := range []string{"user_id", "userId", "total_price", "totalPrice", "nama", "email", "nomer_telpon"} {
				if _, exists := slot[forbidden]; exists {
					t.Fatalf("PRIVACY VIOLATION: availability slot exposes '%s'", forbidden)
				}
			}
			// Status check: must be in BlockingStatuses
			statusStr, ok := slot["status"].(string)
			if !ok || !schemas.IsBlockingStatus(statusStr) {
				t.Errorf("slot %d: status %v is not a blocking status", i, slot["status"])
			}
		}

		// Verify chronological ascending order
		firstStart, _ := time.Parse(time.RFC3339, rawResp.Data[0]["date_start"].(string))
		secondStart, _ := time.Parse(time.RFC3339, rawResp.Data[1]["date_start"].(string))
		if !secondStart.After(firstStart) {
			t.Errorf("expected date_start ascending order, got first=%v second=%v", firstStart, secondStart)
		}
	}

	// 4. Invalid date range parameters return HTTP 400
	{
		// Invalid start string
		req1, _ := http.NewRequest(http.MethodGet, "/api/rooms/"+room1.ID+"/availability?start=invalid-date", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)
		if w1.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid start param, got %d", w1.Code)
		}

		// finish before start
		startStr := now.Add(24 * time.Hour).Format(time.RFC3339)
		finishStr := now.Add(12 * time.Hour).Format(time.RFC3339)
		req2, _ := http.NewRequest(http.MethodGet, "/api/rooms/"+room1.ID+"/availability?start="+startStr+"&finish="+finishStr, nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)
		if w2.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for finish before start, got %d", w2.Code)
		}

		// Range > 1 year
		tooFarFinish := now.Add(400 * 24 * time.Hour).Format(time.RFC3339)
		req3, _ := http.NewRequest(http.MethodGet, "/api/rooms/"+room1.ID+"/availability?start="+startStr+"&finish="+tooFarFinish, nil)
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)
		if w3.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for range > 1 year, got %d", w3.Code)
		}
	}
}
