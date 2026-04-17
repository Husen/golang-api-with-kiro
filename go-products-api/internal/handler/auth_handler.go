package handler

import (
	"go-products-api/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	usecase domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/auth/login", h.Login)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate with username and password, returns a JWT access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      domain.LoginRequest   true  "Login credentials"
// @Success      200          {object}  domain.LoginResponse
// @Failure      400          {object}  map[string]string
// @Failure      401          {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.usecase.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.LoginResponse{AccessToken: token})
}
