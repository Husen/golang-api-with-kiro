package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func runRoleMiddleware(t *testing.T, roleInCtx string, setRole bool, allowedRoles ...string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/test", func(c *gin.Context) {
		if setRole {
			c.Set("role", roleInCtx)
		}
		c.Next()
	}, RoleMiddleware(allowedRoles...), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)
	return w
}

func TestRoleMiddleware_MissingRoleInContext(t *testing.T) {
	w := runRoleMiddleware(t, "", false, "admin")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "role claim missing") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestRoleMiddleware_InsufficientRole_Returns403(t *testing.T) {
	w := runRoleMiddleware(t, "user", true, "admin")
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "forbidden: insufficient role") {
		t.Errorf("unexpected body: %s", w.Body.String())
	}
}

func TestRoleMiddleware_AllowedRole_Passes(t *testing.T) {
	w := runRoleMiddleware(t, "admin", true, "admin")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestRoleMiddleware_UserRoleOnReadRoute_Passes(t *testing.T) {
	w := runRoleMiddleware(t, "user", true, "admin", "user")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthLoginRoute_AccessibleWithoutToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	// Register login route without any middleware (simulates auth handler registration)
	r.POST("/auth/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"access_token": "fake"})
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"username":"admin","password":"admin123"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for /auth/login without token, got %d", w.Code)
	}
}
