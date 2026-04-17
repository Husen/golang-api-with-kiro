package usecase

import (
	"go-products-api/internal/domain"
)

type productUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) GetAll() ([]domain.Product, error) {
	return u.repo.FindAll()
}

func (u *productUsecase) GetByID(id int) (*domain.Product, error) {
	return u.repo.FindByID(id)
}

func (u *productUsecase) Create(req domain.ProductRequest) (*domain.Product, error) {
	p := domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	return u.repo.Create(p)
}

func (u *productUsecase) Update(id int, req domain.ProductRequest) (*domain.Product, error) {
	p := domain.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	return u.repo.Update(id, p)
}

func (u *productUsecase) Delete(id int) error {
	return u.repo.Delete(id)
}
