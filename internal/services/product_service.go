package services

import "ProductService/internal/models"

type productRepository interface {
	GetProducts() ([]models.Product, error)
	GetProduct(id int) (models.Product, error)
}

type ProductService struct {
	productRepository productRepository
}

func NewProductRepository(productRepository productRepository) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) GetProducts() ([]models.Product, error) {
	return s.productRepository.GetProducts()
}

func (s *ProductService) GetProduct(id int) (models.Product, error) {

	return s.productRepository.GetProduct(id)

}
