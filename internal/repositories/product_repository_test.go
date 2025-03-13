package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductRepository(t *testing.T) {

	productRepository := NewProductRepository(db)

	t.Run("GetProducts", func(t *testing.T) {

		products, err := productRepository.GetProducts()

		assert.Nil(t, err)

		assert.Equal(t, 10, len(products))

	})

	t.Run("GetProducts V2", func(t *testing.T) {

		products, err := productRepository.GetProducts()

		assert.Nil(t, err)

		for _, product := range products {

			assert.NotEmpty(t, product.ID)

		}

	})

	t.Run("GetProduct", func(t *testing.T) {

		product, err := productRepository.GetProduct(5)

		assert.Nil(t, err)

		assert.Equal(t, "Klavye", product.Name)
		assert.Equal(t, 5, product.ID)

	})

	t.Run("GetProduct V2", func(t *testing.T) {

		product, err := productRepository.GetProduct(6)

		compreProduct, err := productRepository.GetProduct(1)

		assert.Nil(t, err)

		assert.NotEqual(t, compreProduct, product)

	})

}
