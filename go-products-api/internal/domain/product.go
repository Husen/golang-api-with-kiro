package domain

// Product represents the product entity
// @Description Product data
type Product struct {
	ID    int     `json:"id" example:"1"`
	Name  string  `json:"name" example:"Laptop"`
	Price float64 `json:"price" example:"15000000"`
	Stock int     `json:"stock" example:"10"`
}

// ProductRequest is used for create/update payload
// @Description Product request payload
type ProductRequest struct {
	Name  string  `json:"name" binding:"required" example:"Laptop"`
	Price float64 `json:"price" binding:"required" example:"15000000"`
	Stock int     `json:"stock" binding:"required" example:"10"`
}

// ProductRepository defines the contract for data access
type ProductRepository interface {
	FindAll() ([]Product, error)
	FindByID(id int) (*Product, error)
	Create(p Product) (*Product, error)
	Update(id int, p Product) (*Product, error)
	Delete(id int) error
}

// ProductUsecase defines the contract for business logic
type ProductUsecase interface {
	GetAll() ([]Product, error)
	GetByID(id int) (*Product, error)
	Create(req ProductRequest) (*Product, error)
	Update(id int, req ProductRequest) (*Product, error)
	Delete(id int) error
}
