// @title           Products API
// @version         1.0
// @description     Simple CRUD API for products using Gin + Clean Architecture
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT token only (without "Bearer " prefix). Example: eyJhbGci...

package main

import (
	"go-products-api/internal/handler"
	"go-products-api/internal/middleware"
	"go-products-api/internal/repository"
	"go-products-api/internal/usecase"
	"net/http"
	"os"

	_ "go-products-api/docs"

	"github.com/gin-gonic/gin"
)

func main() {
	productRepo := repository.NewProductRepository()
	productUC := usecase.NewProductUsecase(productRepo)
	productH := handler.NewProductHandler(productUC)

	userRepo := repository.NewUserRepository()
	authUC := usecase.NewAuthUsecase(userRepo)
	authH := handler.NewAuthHandler(authUC)

	r := gin.Default()

	r.GET("/docs/openapi.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", scalarHTML())
	})

	authH.RegisterRoutes(r)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret-key"
	}

	productH.RegisterRoutes(r, middleware.AuthMiddleware(jwtSecret))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func scalarHTML() []byte {
	title := os.Getenv("SCALAR_TITLE")
	if title == "" {
		title = "API Docs"
	}
	specURL := os.Getenv("SCALAR_SPEC_URL")
	if specURL == "" {
		specURL = "/docs/openapi.json"
	}
	html := `<!doctype html>
<html>
<head>
  <title>` + title + `</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>
<body>
  <script id="api-reference" data-url="` + specURL + `"></script>
  <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`
	return []byte(html)
}
