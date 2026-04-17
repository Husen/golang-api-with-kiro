package usecase

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"go-products-api/internal/domain"
	"go-products-api/internal/middleware"
	"go-products-api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"pgregory.net/rapid"
)

// defaultCredentials maps known valid username → plaintext password
var defaultCredentials = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

// makeSignedToken creates a signed JWT with given claims and secret.
func makeSignedToken(secret string, claims *domain.JWTClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

// validClaimsFor builds a valid JWTClaims for the given user data.
func validClaimsFor(userID int, username, role string) *domain.JWTClaims {
	return &domain.JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

// Feature: jwt-auth, Property 1: Login dengan kredensial valid menghasilkan JWT yang lengkap
func TestProperty1_ValidLoginProducesCompleteJWT(t *testing.T) {
	os.Unsetenv("JWT_SECRET")
	rapid.Check(t, func(t *rapid.T) {
		// Pick a random valid username from known credentials
		usernames := []string{"admin", "user"}
		username := rapid.SampledFrom(usernames).Draw(t, "username")
		password := defaultCredentials[username]

		repo := repository.NewUserRepository()
		uc := NewAuthUsecase(repo)

		tokenStr, err := uc.Login(username, password)
		if err != nil {
			t.Fatalf("unexpected error for valid credentials (%s): %v", username, err)
		}

		claims := &domain.JWTClaims{}
		tok, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("dev-secret-key"), nil
		})
		if err != nil || !tok.Valid {
			t.Fatalf("token not parseable: %v", err)
		}

		if claims.Username != username {
			t.Fatalf("expected username '%s', got '%s'", username, claims.Username)
		}
		if claims.Role == "" {
			t.Fatal("role claim is empty")
		}
		if claims.UserID == 0 {
			t.Fatal("user_id claim is zero")
		}
		if tok.Method.Alg() != "HS256" {
			t.Fatalf("expected HS256 signing, got %s", tok.Method.Alg())
		}
		exp := claims.ExpiresAt.Time
		diff := exp.Sub(time.Now())
		if diff < 23*time.Hour || diff > 25*time.Hour {
			t.Fatalf("expiry not ~24h: %v", diff)
		}
	})
}

// Feature: jwt-auth, Property 2: Login dengan kredensial tidak valid selalu ditolak
func TestProperty2_InvalidCredentialsAlwaysRejected(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		username := rapid.StringN(1, 50, -1).Draw(t, "username")
		password := rapid.StringN(1, 50, -1).Draw(t, "password")

		// Skip valid combinations to test only invalid ones
		if pw, ok := defaultCredentials[username]; ok && pw == password {
			t.Skip()
		}

		os.Unsetenv("JWT_SECRET")
		repo := repository.NewUserRepository()
		uc := NewAuthUsecase(repo)

		_, err := uc.Login(username, password)
		if err == nil {
			t.Fatalf("expected error for invalid credentials (%s / %s), got nil", username, password)
		}
	})
}

// Feature: jwt-auth, Property 3: Role pengguna yang tersimpan selalu bernilai valid
func TestProperty3_StoredRolesAreAlwaysValid(t *testing.T) {
	validRoles := map[string]bool{"admin": true, "user": true}
	rapid.Check(t, func(t *rapid.T) {
		repo := repository.NewUserRepository()
		for _, username := range []string{"admin", "user"} {
			u, err := repo.FindByUsername(username)
			if err != nil {
				t.Fatalf("user '%s' not found: %v", username, err)
			}
			if !validRoles[u.Role] {
				t.Fatalf("invalid role '%s' for user '%s'", u.Role, username)
			}
		}
	})
}

// Feature: jwt-auth, Property 4: FindByUsername adalah round-trip yang konsisten
func TestProperty4_FindByUsernameRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		registeredUsernames := []string{"admin", "user"}
		username := rapid.SampledFrom(registeredUsernames).Draw(t, "username")

		repo := repository.NewUserRepository()
		u, err := repo.FindByUsername(username)
		if err != nil {
			t.Fatalf("expected to find '%s', got error: %v", username, err)
		}
		if u.Username != username {
			t.Fatalf("round-trip failed: queried '%s', got '%s'", username, u.Username)
		}
	})
}

// Feature: jwt-auth, Property 4 (complement): unregistered usernames always return error
func TestProperty4_UnregisteredUsernameAlwaysErrors(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		username := rapid.StringN(10, 30, -1).Draw(t, "username")
		// Avoid accidentally matching registered users
		if username == "admin" || username == "user" {
			t.Skip()
		}
		repo := repository.NewUserRepository()
		_, err := repo.FindByUsername(username)
		if err == nil {
			t.Fatalf("expected error for unregistered username '%s', got nil", username)
		}
	})
}

// Feature: jwt-auth, Property 5: Password tersimpan sebagai bcrypt hash yang dapat diverifikasi
func TestProperty5_PasswordStoredAsBcryptHash(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		cases := []struct{ username, plaintext string }{
			{"admin", "admin123"},
			{"user", "user123"},
		}
		idx := rapid.IntRange(0, len(cases)-1).Draw(t, "idx")
		tc := cases[idx]

		repo := repository.NewUserRepository()
		u, err := repo.FindByUsername(tc.username)
		if err != nil {
			t.Fatalf("user not found: %v", err)
		}
		if u.PasswordHash == tc.plaintext {
			t.Fatal("password stored as plaintext, not hash")
		}
		// bcrypt.CompareHashAndPassword validation is done indirectly through Login
		os.Unsetenv("JWT_SECRET")
		uc := NewAuthUsecase(repo)
		_, loginErr := uc.Login(tc.username, tc.plaintext)
		if loginErr != nil {
			t.Fatalf("bcrypt verification failed: %v", loginErr)
		}
	})
}

// Feature: jwt-auth, Property 6: Inisialisasi repository bersifat idempotent
func TestProperty6_RepositoryInitIsIdempotent(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		r1 := repository.NewUserRepository()
		r2 := repository.NewUserRepository()

		for _, username := range []string{"admin", "user"} {
			u1, err1 := r1.FindByUsername(username)
			u2, err2 := r2.FindByUsername(username)
			if err1 != nil || err2 != nil {
				t.Fatalf("user '%s' missing in one of the repos", username)
			}
			if u1.Username != u2.Username || u1.Role != u2.Role || u1.ID != u2.ID {
				t.Fatalf("repos differ for '%s': %+v vs %+v", username, u1, u2)
			}
		}
	})
}

// Feature: jwt-auth, Property 7: Token valid melewati middleware dan claims tersimpan di context
func TestProperty7_ValidTokenPassesMiddlewareAndSetsContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "prop7-secret"

	rapid.Check(t, func(t *rapid.T) {
		usernames := []string{"admin", "user"}
		roles := []string{"admin", "user"}
		username := rapid.SampledFrom(usernames).Draw(t, "username")
		role := rapid.SampledFrom(roles).Draw(t, "role")
		userID := rapid.IntRange(1, 100).Draw(t, "userID")

		claims := validClaimsFor(userID, username, role)
		tokenStr := makeSignedToken(secret, claims)

		var (
			gotUID      interface{}
			gotUsername interface{}
			gotRole     interface{}
			handlerCalled bool
		)

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/test", middleware.AuthMiddleware(secret), func(c *gin.Context) {
			gotUID, _ = c.Get("user_id")
			gotUsername, _ = c.Get("username")
			gotRole, _ = c.Get("role")
			handlerCalled = true
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		r.ServeHTTP(w, req)

		if !handlerCalled {
			t.Fatal("handler was not called for valid token")
		}
		if gotUID != userID {
			t.Fatalf("expected user_id=%d, got %v", userID, gotUID)
		}
		if gotUsername != username {
			t.Fatalf("expected username='%s', got %v", username, gotUsername)
		}
		if gotRole != role {
			t.Fatalf("expected role='%s', got %v", role, gotRole)
		}
	})
}

// Feature: jwt-auth, Property 8: Token dengan format atau signature yang salah selalu ditolak
func TestProperty8_InvalidTokenAlwaysRejected(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "prop8-secret"

	rapid.Check(t, func(t *rapid.T) {
		// Generate random strings that are not valid tokens signed with our secret
		invalidToken := rapid.StringN(1, 200, -1).Draw(t, "invalidToken")

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		r.GET("/test", middleware.AuthMiddleware(secret), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+invalidToken)
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Fatalf("expected rejection for token '%s', got 200", invalidToken)
		}
	})
}

// Feature: jwt-auth, Property 9: Izin write bergantung pada role admin
func TestProperty9_WritePermissionRequiresAdminRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rapid.Check(t, func(t *rapid.T) {
		roles := []string{"admin", "user"}
		role := rapid.SampledFrom(roles).Draw(t, "role")

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)

		r.POST("/products", func(c *gin.Context) {
			c.Set("role", role)
		}, middleware.RoleMiddleware("admin"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodPost, "/products", nil)
		r.ServeHTTP(w, req)

		if role == "admin" {
			if w.Code != http.StatusOK {
				t.Fatalf("admin should pass write middleware, got %d", w.Code)
			}
		} else {
			if w.Code != http.StatusForbidden {
				t.Fatalf("non-admin should get 403, got %d", w.Code)
			}
		}
	})
}

// Feature: jwt-auth, Property 10: Semua role valid dapat melakukan read operations
func TestProperty10_AllValidRolesCanRead(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rapid.Check(t, func(t *rapid.T) {
		roles := []string{"admin", "user"}
		role := rapid.SampledFrom(roles).Draw(t, "role")

		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)

		r.GET("/products", func(c *gin.Context) {
			c.Set("role", role)
		}, middleware.RoleMiddleware("admin", "user"), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("role '%s' should be allowed to read, got %d", role, w.Code)
		}
	})
}
