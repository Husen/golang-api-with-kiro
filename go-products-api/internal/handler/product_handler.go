package handler

import (
	"go-products-api/internal/domain"
	"go-products-api/internal/middleware"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	usecase domain.ProductUsecase
}

func NewProductHandler(uc domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: uc}
}

func (h *ProductHandler) RegisterRoutes(r *gin.Engine, authMiddleware gin.HandlerFunc) {
	products := r.Group("/products", authMiddleware)
	{
		products.GET("", middleware.RoleMiddleware("admin", "user"), h.GetAll)
		products.GET("/:id", middleware.RoleMiddleware("admin", "user"), h.GetByID)
		products.POST("", middleware.RoleMiddleware("admin"), h.Create)
		products.PUT("/:id", middleware.RoleMiddleware("admin"), h.Update)
		products.DELETE("/:id", middleware.RoleMiddleware("admin"), h.Delete)
	}
}

// GetAll godoc
// @Summary      List all products
// @Tags         products
// @Produce      json
// @Success      200  {array}   domain.Product
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Security     BearerAuth
// @Router       /products [get]
func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetByID godoc
// @Summary      Get product by ID
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  domain.Product
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /products/{id} [get]
func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	product, err := h.usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// Create godoc
// @Summary      Create a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      domain.ProductRequest  true  "Product payload"
// @Success      201      {object}  domain.Product
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Security     BearerAuth
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var req domain.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.usecase.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

// Update godoc
// @Summary      Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "Product ID"
// @Param        product  body      domain.ProductRequest  true  "Product payload"
// @Success      200      {object}  domain.Product
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Security     BearerAuth
// @Router       /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.usecase.Update(id, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

// Delete godoc
// @Summary      Delete a product
// @Tags         products
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.usecase.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
