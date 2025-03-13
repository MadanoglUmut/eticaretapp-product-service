package repositories

import (
	"ProductService/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProducts() ([]models.Product, error) {

	var products []models.Product

	if err := r.db.Table("products").Debug().Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil

}

func (r *ProductRepository) GetProduct(id int) (models.Product, error) {

	var product models.Product
	if err := r.db.Table("products").Debug().First(&product, id).Error; err != nil {

		return models.Product{}, err

	}

	return product, nil

}
