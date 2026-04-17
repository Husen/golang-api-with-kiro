package repository

import (
	"errors"
	"go-products-api/internal/domain"
	"sync"
)

type productRepository struct {
	mu       sync.Mutex
	products []domain.Product
	nextID   int
}

func NewProductRepository() domain.ProductRepository {
	return &productRepository{
		nextID: 4,
		products: []domain.Product{
			{ID: 1, Name: "Laptop", Price: 15000000, Stock: 10},
			{ID: 2, Name: "Mouse", Price: 150000, Stock: 50},
			{ID: 3, Name: "Keyboard", Price: 300000, Stock: 30},
		},
	}
}

func (r *productRepository) FindAll() ([]domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]domain.Product, len(r.products))
	copy(result, r.products)
	return result, nil
}

func (r *productRepository) FindByID(id int) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, p := range r.products {
		if p.ID == id {
			cp := p
			return &cp, nil
		}
	}
	return nil, errors.New("product not found")
}

func (r *productRepository) Create(p domain.Product) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	p.ID = r.nextID
	r.nextID++
	r.products = append(r.products, p)
	return &p, nil
}

func (r *productRepository) Update(id int, p domain.Product) (*domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, existing := range r.products {
		if existing.ID == id {
			p.ID = id
			r.products[i] = p
			return &p, nil
		}
	}
	return nil, errors.New("product not found")
}

func (r *productRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, p := range r.products {
		if p.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return nil
		}
	}
	return errors.New("product not found")
}
