package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	schemas "stpNew/backend/schema"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// 1. RFC 3339 timestamp parsing and UTC normalization
func TestValidateBookingInterval_RFC3339ParsingAndUTC(t *testing.T) {
	now := time.Date(2026, 11, 1, 6, 0, 0, 0, time.UTC)

	// Test +07:00 (WIB) offset converted to UTC
	startStr := "2026-11-01T14:00:00+07:00"  // 07:00 UTC
	finishStr := "2026-11-01T16:00:00+07:00" // 09:00 UTC

	startUTC, finishUTC, err := ValidateBookingInterval(startStr, finishStr, now)
	if err != nil {
		t.Fatalf("unexpected error parsing interval: %v", err)
	}

	if startUTC.Location() != time.UTC {
		t.Errorf("expected start location UTC, got %v", startUTC.Location())
	}
	if finishUTC.Location() != time.UTC {
		t.Errorf("expected finish location UTC, got %v", finishUTC.Location())
	}

	expectedStartHour := 7
	if startUTC.Hour() != expectedStartHour {
		t.Errorf("expected start UTC hour %d, got %d", expectedStartHour, startUTC.Hour())
	}

	expectedFinishHour := 9
	if finishUTC.Hour() != expectedFinishHour {
		t.Errorf("expected finish UTC hour %d, got %d", expectedFinishHour, finishUTC.Hour())
	}
}

// 2. Finish must be strictly after start
func TestValidateBookingInterval_FinishAfterStart(t *testing.T) {
	now := time.Date(2026, 11, 1, 6, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		startStr    string
		finishStr   string
		expectError bool
		errMsg      string
	}{
		{
			name:        "finish is after start",
			startStr:    "2026-11-01T08:00:00Z",
			finishStr:   "2026-11-01T10:00:00Z",
			expectError: false,
		},
		{
			name:        "finish equals start",
			startStr:    "2026-11-01T08:00:00Z",
			finishStr:   "2026-11-01T08:00:00Z",
			expectError: true,
			errMsg:      "waktu selesai harus lebih lambat dari waktu mulai",
		},
		{
			name:        "finish is before start",
			startStr:    "2026-11-01T10:00:00Z",
			finishStr:   "2026-11-01T08:00:00Z",
			expectError: true,
			errMsg:      "waktu selesai harus lebih lambat dari waktu mulai",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ValidateBookingInterval(tt.startStr, tt.finishStr, now)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error containing '%s', got nil", tt.errMsg)
				}
				if err.Error() != tt.errMsg {
					t.Errorf("expected error '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// 3. Past start times are rejected, with a 5-minute clock tolerance
func TestValidateBookingInterval_PastStartRejected(t *testing.T) {
	referenceTime := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)

	// Past by 1 hour (outside 5-minute tolerance)
	pastStart := "2026-11-01T11:00:00Z"
	finish := "2026-11-01T13:00:00Z"
	_, _, err := ValidateBookingInterval(pastStart, finish, referenceTime)
	if err == nil {
		t.Fatal("expected error for start time 1 hour in past, got nil")
	}
	if err.Error() != "waktu mulai tidak boleh di masa lampau" {
		t.Errorf("expected past time error message, got: %s", err.Error())
	}

	// Past by 2 minutes (inside 5-minute clock drift tolerance)
	recentPastStart := "2026-11-01T11:58:00Z"
	_, _, err = ValidateBookingInterval(recentPastStart, finish, referenceTime)
	if err != nil {
		t.Fatalf("expected 2-minute past drift to be accepted within 5-minute tolerance, got: %v", err)
	}

	// Future start time
	futureStart := "2026-11-01T14:00:00Z"
	futureFinish := "2026-11-01T16:00:00Z"
	_, _, err = ValidateBookingInterval(futureStart, futureFinish, referenceTime)
	if err != nil {
		t.Fatalf("expected future interval to be accepted, got: %v", err)
	}
}

// 4. Invalid timestamp format rejection
func TestValidateBookingInterval_InvalidFormats(t *testing.T) {
	now := time.Date(2026, 11, 1, 6, 0, 0, 0, time.UTC)

	invalidInputs := []struct {
		name      string
		startStr  string
		finishStr string
	}{
		{
			name:      "slash date format",
			startStr:  "11/01/2026 08:00:00",
			finishStr: "11/01/2026 10:00:00",
		},
		{
			name:      "dash without time",
			startStr:  "2026-11-01",
			finishStr: "2026-11-02",
		},
		{
			name:      "empty start",
			startStr:  "",
			finishStr: "2026-11-01T10:00:00Z",
		},
		{
			name:      "empty finish",
			startStr:  "2026-11-01T08:00:00Z",
			finishStr: "",
		},
		{
			name:      "arbitrary non-date string",
			startStr:  "not-a-date",
			finishStr: "not-a-date",
		},
	}

	for _, tc := range invalidInputs {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ValidateBookingInterval(tc.startStr, tc.finishStr, now)
			if err == nil {
				t.Fatalf("expected error for invalid format in %s, got nil", tc.name)
			}
		})
	}
}

// 5. Overlap semantics: [start, finish) half-open intervals
func TestHasIntervalOverlap_Semantics(t *testing.T) {
	t0 := time.Date(2026, 11, 1, 8, 0, 0, 0, time.UTC)
	t1 := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 11, 1, 14, 0, 0, 0, time.UTC)

	// Adjacent intervals: [t0, t1) and [t1, t2)
	// Booking 1 ends exactly when Booking 2 begins -> NO CONFLICT
	if HasIntervalOverlap(t0, t1, t1, t2) {
		t.Errorf("adjacent intervals [t0, t1) and [t1, t2) must not conflict")
	}
	if HasIntervalOverlap(t1, t2, t0, t1) {
		t.Errorf("reverse adjacent intervals [t1, t2) and [t0, t1) must not conflict")
	}

	// Partial overlap: [08:00, 11:00) and [10:00, 12:00)
	midPoint := time.Date(2026, 11, 1, 11, 0, 0, 0, time.UTC)
	if !HasIntervalOverlap(t0, midPoint, t1, t2) {
		t.Errorf("partial overlap [t0, midPoint) and [t1, t2) must conflict")
	}

	// Complete containment: [08:00, 14:00) and [10:00, 12:00)
	if !HasIntervalOverlap(t0, t3, t1, t2) {
		t.Errorf("containment [t0, t3) containing [t1, t2) must conflict")
	}

	// Identical intervals: [t0, t1) and [t0, t1)
	if !HasIntervalOverlap(t0, t1, t0, t1) {
		t.Errorf("identical intervals must conflict")
	}

	// Completely disjoint intervals: [t0, t1) and [t2, t3)
	if HasIntervalOverlap(t0, t1, t2, t3) {
		t.Errorf("disjoint intervals must not conflict")
	}
}

// 6. Exclusion constraint error detection (SQLSTATE 23P01)
func TestIsExclusionConstraintViolation(t *testing.T) {
	// A: Typed pgconn.PgError with 23P01
	pgExclusionErr := &pgconn.PgError{Code: "23P01", Message: "conflicting key value violates exclusion constraint"}
	if !IsExclusionConstraintViolation(pgExclusionErr) {
		t.Error("expected IsExclusionConstraintViolation to return true for pgconn.PgError Code 23P01")
	}

	// B: Wrapped pgconn.PgError with 23P01
	wrappedErr := errors.Join(errors.New("transaction commit failed"), pgExclusionErr)
	if !IsExclusionConstraintViolation(wrappedErr) {
		t.Error("expected IsExclusionConstraintViolation to return true for wrapped pgconn.PgError Code 23P01")
	}

	// C: Typed pgconn.PgError with different code (e.g. 23505 Unique Violation)
	pgUniqueErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value"}
	if IsExclusionConstraintViolation(pgUniqueErr) {
		t.Error("expected IsExclusionConstraintViolation to return false for Code 23505")
	}

	// D: Generic error with code string (must NOT match; only typed errors are accepted)
	errWithStringCode := errors.New("ERROR: conflicting key value (SQLSTATE 23P01)")
	if IsExclusionConstraintViolation(errWithStringCode) {
		t.Error("expected IsExclusionConstraintViolation to return false for non-typed string containing 23P01")
	}

	// E: Generic error with constraint name (must NOT match)
	errWithConstraintName := errors.New("violates exclusion constraint 'exclude_overlapping_room_bookings'")
	if IsExclusionConstraintViolation(errWithConstraintName) {
		t.Error("expected IsExclusionConstraintViolation to return false for non-typed string with constraint name")
	}

	// F: Unrelated database error
	otherErr := errors.New("connection refused")
	if IsExclusionConstraintViolation(otherErr) {
		t.Error("expected IsExclusionConstraintViolation to return false for generic connection error")
	}

	// G: Nil error
	if IsExclusionConstraintViolation(nil) {
		t.Error("expected IsExclusionConstraintViolation to return false for nil")
	}
}

// 7. Client-supplied fields cannot override UserID, Price, or Status
func TestBookingSchema_NoClientOverride(t *testing.T) {
	payload := `{
		"room_id": "room-uuid-1",
		"date_start": "2026-11-01T08:00:00Z",
		"date_finish": "2026-11-01T10:00:00Z",
		"user_id": "malicious-user-override",
		"price": 10,
		"total_price": 10,
		"status": "Selesai"
	}`

	var schema schemas.BookingSchema
	if err := json.Unmarshal([]byte(payload), &schema); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if schema.RoomID != "room-uuid-1" {
		t.Errorf("expected room_id 'room-uuid-1', got '%s'", schema.RoomID)
	}
	if schema.DateStart != "2026-11-01T08:00:00Z" {
		t.Errorf("expected date_start '2026-11-01T08:00:00Z', got '%s'", schema.DateStart)
	}
	if schema.DateFinish != "2026-11-01T10:00:00Z" {
		t.Errorf("expected date_finish '2026-11-01T10:00:00Z', got '%s'", schema.DateFinish)
	}

	// The BookingSchema struct does not have user_id, price, total_price, or status fields,
	// so client-provided overrides are entirely ignored by the schema parser.
}

// 8. HTTP Handler: CreateBooking rejects unauthorized requests (missing auth context)
func TestCreateBooking_UnauthorizedWithoutAuthContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := bytes.NewBufferString(`{"room_id": "1", "date_start": "2026-11-01T08:00:00Z", "date_finish": "2026-11-01T10:00:00Z"}`)
	req, _ := http.NewRequest(http.MethodPost, "/api/Bookings", body)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	CreateBooking(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected HTTP 401 Unauthorized for missing auth context, got %d", w.Code)
	}
}

// 9. HTTP Handler: CreateBooking rejects invalid request bodies
func TestCreateBooking_ValidationFailures(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{
			name:         "malformed JSON",
			body:         `{invalid-json}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "missing room_id",
			body:         `{"date_start": "2026-11-01T08:00:00Z", "date_finish": "2026-11-01T10:00:00Z"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid timestamp format",
			body:         `{"room_id": "1", "date_start": "01/11/2026", "date_finish": "01/11/2026"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "finish before start",
			body:         `{"room_id": "1", "date_start": "2026-11-01T12:00:00Z", "date_finish": "2026-11-01T10:00:00Z"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "past start time",
			body:         `{"room_id": "1", "date_start": "2020-01-01T08:00:00Z", "date_finish": "2020-01-01T10:00:00Z"}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("id", "test-user-id-123")

			req, _ := http.NewRequest(http.MethodPost, "/api/Bookings", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			CreateBooking(c)

			if w.Code != tt.expectedCode {
				t.Fatalf("%s: expected HTTP %d, got %d (%s)", tt.name, tt.expectedCode, w.Code, w.Body.String())
			}
		})
	}
}

// 10. HTTP Handler: GetBookings rejects unauthorized requests
func TestGetBookings_UnauthorizedWithoutAuthContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/api/Bookings", nil)
	c.Request = req

	GetBookings(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected HTTP 401 Unauthorized for missing auth context in GetBookings, got %d", w.Code)
	}
}

// 11. HTTP Handler: DashboardPage returns HTTP 500 when DB is nil
func TestDashboardPage_NilDatabase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest(http.MethodGet, "/api/Beranda", nil)
	c.Request = req

	DashboardPage(c)

	// Since databases.DB is nil during isolated unit tests, DashboardPage safely handles it
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected HTTP 500 when databases.DB is nil, got %d", w.Code)
	}
}
