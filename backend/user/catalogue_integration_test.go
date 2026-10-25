package users_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	databases "stpNew/backend/database"
	schemas "stpNew/backend/schema"
	"stpNew/backend/testutil"
	users "stpNew/backend/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

// DevRoomSeedsForTest matches the authoritative seed defined in backend/cmd/seed/main.go
var DevRoomSeedsForTest = []schemas.Room{
	{
		ID:          "1",
		Name:        "Sirkuit Otomotif ITS",
		Description: "Fasilitas Pengujian & Balap Sirkuit Otomotif Kampus ITS Sukolilo, Surabaya",
		Capacity:    100,
		Price:       150000,
		ImageURL:    "https://images.unsplash.com/photo-1568605117036-5fe5e7bab0b7?auto=format&fit=crop&q=80&w=1000",
		Features:    "Sirkuit Balap & Pengujian",
	},
}

// Step 13: Room Seed Verification
func TestRoomSeed_Idempotency(t *testing.T) {
	gormDB, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	// 1. First run of seeder
	for _, r := range DevRoomSeedsForTest {
		if err := gormDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&r).Error; err != nil {
			t.Fatalf("seeder first run failed: %v", err)
		}
	}

	var firstCount int64
	gormDB.Model(&schemas.Room{}).Count(&firstCount)
	if firstCount != 1 {
		t.Fatalf("expected 1 room after first seed, got %d", firstCount)
	}

	var room schemas.Room
	gormDB.First(&room, "id = ?", "1")
	if room.Name != "Sirkuit Otomotif ITS" {
		t.Fatalf("expected room name 'Sirkuit Otomotif ITS', got '%s'", room.Name)
	}
	if room.Price <= 0 {
		t.Fatalf("expected non-negative room price, got %f", room.Price)
	}

	// 2. Second run of seeder (idempotency check)
	for _, r := range DevRoomSeedsForTest {
		if err := gormDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&r).Error; err != nil {
			t.Fatalf("seeder second run failed: %v", err)
		}
	}

	var secondCount int64
	gormDB.Model(&schemas.Room{}).Count(&secondCount)
	if secondCount != 1 {
		t.Fatalf("expected 1 room after second seed (no duplicates), got %d", secondCount)
	}

	// 3. User modifications are preserved
	gormDB.Model(&schemas.Room{}).Where("id = ?", "1").Update("description", "Deskripsi yang diubah pengguna")
	for _, r := range DevRoomSeedsForTest {
		gormDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&r)
	}
	var modifiedRoom schemas.Room
	gormDB.First(&modifiedRoom, "id = ?", "1")
	if modifiedRoom.Description != "Deskripsi yang diubah pengguna" {
		t.Errorf("seeder should not overwrite existing user-modified record, got: %s", modifiedRoom.Description)
	}
}

// Step 14: Database-Backed Catalogue Test (GET /api/Beranda)
func TestCatalogue_GetBerandaAgainstPostgres(t *testing.T) {
	gormDB, sqlDB := testutil.GetTestDB(t)
	testutil.ResetTestDBSchema(t, sqlDB)
	if err := databases.ApplyMigrationsUp(sqlDB); err != nil {
		t.Fatalf("migration up failed: %v", err)
	}

	databases.DB = gormDB

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/Beranda", users.DashboardPage)

	// Case 1: Empty catalogue returns HTTP 200 with data: []
	reqEmpty, _ := http.NewRequest(http.MethodGet, "/api/Beranda", nil)
	wEmpty := httptest.NewRecorder()
	router.ServeHTTP(wEmpty, reqEmpty)

	if wEmpty.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200 for empty catalogue, got %d", wEmpty.Code)
	}

	var emptyResp struct {
		Message string         `json:"message"`
		Data    []schemas.Room `json:"data"`
	}
	if err := json.Unmarshal(wEmpty.Body.Bytes(), &emptyResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if emptyResp.Data == nil || len(emptyResp.Data) != 0 {
		t.Fatalf("expected non-null empty array [], got %v", emptyResp.Data)
	}

	// Case 2: Seeded catalogue returns database records
	for _, r := range DevRoomSeedsForTest {
		gormDB.Clauses(clause.OnConflict{DoNothing: true}).Create(&r)
	}

	reqSeeded, _ := http.NewRequest(http.MethodGet, "/api/Beranda", nil)
	wSeeded := httptest.NewRecorder()
	router.ServeHTTP(wSeeded, reqSeeded)

	if wSeeded.Code != http.StatusOK {
		t.Fatalf("expected HTTP 200 for seeded catalogue, got %d", wSeeded.Code)
	}

	var seededResp struct {
		Message string         `json:"message"`
		Data    []schemas.Room `json:"data"`
	}
	if err := json.Unmarshal(wSeeded.Body.Bytes(), &seededResp); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if len(seededResp.Data) != 1 {
		t.Fatalf("expected 1 room returned from database, got %d", len(seededResp.Data))
	}
	if seededResp.Data[0].Name != "Sirkuit Otomotif ITS" {
		t.Errorf("expected room 'Sirkuit Otomotif ITS', got '%s'", seededResp.Data[0].Name)
	}
}
