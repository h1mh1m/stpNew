package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	schemas "stpNew/backend/schema"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-jwt-secret-key-12345"

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Recovery())

	protected := router.Group("/api")
	protected.Use(AuthMiddlewareJwt())
	{
		protected.GET("/test-protected", func(c *gin.Context) {
			userID, _ := c.Get("id")
			c.JSON(http.StatusOK, gin.H{
				"message": "access granted",
				"user_id": userID,
			})
		})
	}

	return router
}

// 1. Token generation succeeds with a valid JWT secret.
func TestTokenGeneration_SuccessWithValidSecret(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	token, err := GenerateToken("customer-101")
	if err != nil {
		t.Fatalf("expected token generation to succeed, got error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token string")
	}
}

// 2. Generated token contains the expected customer ID.
func TestToken_ContainsExpectedCustomerID(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	expectedID := "customer-uuid-202"
	tokenStr, err := GenerateToken(expectedID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims := &schemas.Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(testSecret), nil
	})
	if err != nil || !parsedToken.Valid {
		t.Fatalf("failed to parse generated token: %v", err)
	}
	if claims.IdUser != expectedID {
		t.Fatalf("expected customer ID '%s', got '%s'", expectedID, claims.IdUser)
	}
}

// 3. Generated token uses HS256.
func TestToken_UsesHS256(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	tokenStr, err := GenerateToken("customer-303")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	parser := jwt.NewParser()
	claims := &schemas.Claims{}
	token, _, err := parser.ParseUnverified(tokenStr, claims)
	if err != nil {
		t.Fatalf("failed to parse token unverified: %v", err)
	}

	alg, ok := token.Header["alg"].(string)
	if !ok || alg != "HS256" {
		t.Fatalf("expected alg 'HS256', got '%v'", token.Header["alg"])
	}
	if token.Method.Alg() != "HS256" {
		t.Fatalf("expected signing method 'HS256', got '%s'", token.Method.Alg())
	}
}

// 4. Valid token is accepted by AuthMiddlewareJwt.
func TestAuthMiddleware_ValidTokenAccepted(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	router := setupTestRouter()
	tokenStr, err := GenerateToken("customer-404")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}
	if resp["user_id"] != "customer-404" {
		t.Fatalf("expected context user_id 'customer-404', got '%v'", resp["user_id"])
	}
}

// 5. Missing Authorization header is rejected.
func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

// 6. Incorrect authorization scheme is rejected.
func TestAuthMiddleware_IncorrectAuthorizationScheme(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	router := setupTestRouter()
	schemes := []string{
		"Basic dXNlcjpwYXNz",
		"Token some-token-value",
		"bearer lowercased-scheme",
		"Bearer", // missing token part
	}

	for _, headerVal := range schemes {
		req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
		req.Header.Set("Authorization", headerVal)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Fatalf("scheme '%s': expected status 401, got %d", headerVal, w.Code)
		}
	}
}

// 7. Malformed bearer token is rejected.
func TestAuthMiddleware_MalformedBearerToken(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.jwt.token")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

// 8. Wrong signature is rejected.
func TestAuthMiddleware_WrongSignature(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	// Generate token with a different secret
	otherSecret := []byte("completely-different-signing-key")
	claim := &schemas.Claims{
		IdUser: "customer-808",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	wrongToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := wrongToken.SignedString(otherSecret)
	if err != nil {
		t.Fatalf("failed to create wrong signature token: %v", err)
	}

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

// 9. Non-HS256 token is rejected (e.g. HS384).
func TestAuthMiddleware_NonHS256Token(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	// Sign with HS384 instead of HS256
	claim := &schemas.Claims{
		IdUser: "customer-909",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	hs384Token := jwt.NewWithClaims(jwt.SigningMethodHS384, claim)
	tokenStr, err := hs384Token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign token with HS384: %v", err)
	}

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for non-HS256 token, got %d", w.Code)
	}
}

// 10. Expired token is rejected.
func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	// Create expired token (expired 1 hour ago)
	claim := &schemas.Claims{
		IdUser: "customer-1010",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := expiredToken.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign expired token: %v", err)
	}

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for expired token, got %d", w.Code)
	}
}

// 11. Token missing the customer ID claim is rejected.
func TestAuthMiddleware_MissingCustomerIdClaim(t *testing.T) {
	os.Setenv("JWT", testSecret)
	defer os.Unsetenv("JWT")

	// Token with empty IdUser
	claim := &schemas.Claims{
		IdUser: "",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	noIdToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenStr, err := noIdToken.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign token without customer ID: %v", err)
	}

	router := setupTestRouter()
	req, _ := http.NewRequest(http.MethodGet, "/api/test-protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for token missing customer ID, got %d", w.Code)
	}
}

// 12. Empty or unsafe JWT secret causes token generation to fail safely.
func TestTokenGeneration_EmptyJWTSecretFailsSafely(t *testing.T) {
	// Case A: Unset / empty
	os.Unsetenv("JWT")
	token, err := GenerateToken("customer-1212")
	if err == nil || token != "" {
		t.Fatalf("expected error on unset JWT secret, got token: '%s'", token)
	}

	// Case B: Blank whitespace
	os.Setenv("JWT", "   ")
	token, err = GenerateToken("customer-1212")
	if err == nil || token != "" {
		t.Fatalf("expected error on whitespace JWT secret, got token: '%s'", token)
	}

	// Case C: Too short (< 8 chars)
	os.Setenv("JWT", "short")
	token, err = GenerateToken("customer-1212")
	if err == nil || token != "" {
		t.Fatalf("expected error on short JWT secret, got token: '%s'", token)
	}
	os.Unsetenv("JWT")
}
