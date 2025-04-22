package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestProductHandler(t *testing.T) {

	t.Run("SuccesGetProductsHandler", func(t *testing.T) {

		req := httptest.NewRequest("GET", "/products", nil)

		resp, err := app.Test(req)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	})

	t.Run("SuccesGetProductHandler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/products/1", nil)

		resp, err := app.Test(req)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("ErrorGetProductHandler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/products/999", nil)

		resp, err := app.Test(req)

		assert.Nil(t, err)

		assert.Equal(t, fiber.StatusServiceUnavailable, resp.StatusCode)
	})

}
