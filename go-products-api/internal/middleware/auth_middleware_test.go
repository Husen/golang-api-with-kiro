package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-products-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret"

func makeToken(secret string, claims *domain.JWTClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func validClaims() *domain.JWTClaims {
	return &domain.JWTClaims{
		UserID:   1,
		Username: "admin",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func runMiddleware(t *testing.T, authHeader string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)
	r.GET("/test", AuthMiddleware(testSecret), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	c.Request = req
	r.ServeHTTP(w, req)
	return w
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	w := runMiddleware(t, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !contains(w.Body.String(), "authorization header required") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestAuthMiddleware_InvalidFormat_NoBearer(t *testing.T) {
	w := runMiddleware(t, "Token sometoken")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !contains(w.Body.String(), "invalid authorization format") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestAuthMiddleware_InvalidFormat_BearerOnly(t *testing.T) {
	w := runMiddleware(t, "Bearer")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	w := runMiddleware(t, "Bearer notavalidtoken")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !contains(w.Body.String(), "invalid token") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	claims := &domain.JWTClaims{
		UserID:   1,
		Username: "admin",
		Role:     "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := makeToken(testSecret, claims)
	w := runMiddleware(t, "Bearer "+token)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !contains(w.Body.String(), "token expired") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestAuthMiddleware_MissingRoleClaim(t *testing.T) {
	claims := &domain.JWTClaims{
		UserID:   1,
		Username: "admin",
		Role:     "", // empty role
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := makeToken(testSecret, claims)
	w := runMiddleware(t, "Bearer "+token)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !contains(w.Body.String(), "invalid token claims") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestAuthMiddleware_ValidToken_SetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	var (
		gotUserID   int
		gotUsername string
		gotRole     string
	)

	r.GET("/test", AuthMiddleware(testSecret), func(c *gin.Context) {
		gotUserID, _ = c.Get("user_id").(int)
		gotUsername, _ = c.Get("username").(string)
		gotRole, _ = c.Get("role").(string)
		c.Status(http.StatusOK)
	})

	token := makeToken(testSecret, validClaims())
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if gotUserID != 1 {
		t.Errorf("expected user_id=1, got %d", gotUserID)
	}
	if gotUsername != "admin" {
		t.Errorf("expected username='admin', got '%s'", gotUsername)
	}
	if gotRole != "admin" {
		t.Errorf("expected role='admin', got '%s'", gotRole)
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	token := makeToken("other-secret", validClaims())
	w := runMiddleware(t, "Bearer "+token)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
