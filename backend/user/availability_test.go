package users

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestGetRoomAvailability_Unauthorized verifies missing auth context returns HTTP 401
func TestGetRoomAvailability_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/api/rooms/room-123/availability", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "room_id", Value: "room-123"}}

	GetRoomAvailability(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected HTTP 401 Unauthorized for missing auth context, got %d", w.Code)
	}
}

// TestGetRoomAvailability_NilDatabase verifies HTTP 500 when database connection is nil
func TestGetRoomAvailability_NilDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("id", "cust-test-123")
	req, _ := http.NewRequest(http.MethodGet, "/api/rooms/room-123/availability", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "room_id", Value: "room-123"}}

	GetRoomAvailability(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected HTTP 500 when database is nil, got %d (%s)", w.Code, w.Body.String())
	}
}

// TestGetRoomAvailability_EmptyRoomID verifies HTTP 400 when room_id is empty or whitespace
func TestGetRoomAvailability_EmptyRoomID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("id", "cust-test-123")
	req, _ := http.NewRequest(http.MethodGet, "/api/rooms/%20/availability", nil)
	c.Request = req
	c.Params = gin.Params{{Key: "room_id", Value: "   "}}

	GetRoomAvailability(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected HTTP 400 for empty room_id, got %d", w.Code)
	}
}
