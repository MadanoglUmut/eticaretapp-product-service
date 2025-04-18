package handlers

import (
	"ProductService/internal/models"
	"math/rand"
	"time"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type productService interface {
	GetProducts() ([]models.Product, error)
	GetProduct(id int) (models.Product, error)
}

type metric interface {
	ObserveHandler(name string, startTime time.Time, status int)
}

type ProductHandler struct {
	productService productService
	metric         metric
}

func NewProductHandler(productService productService, metric metric) *ProductHandler {

	return &ProductHandler{
		productService: productService,
		metric:         metric,
	}

}

func (h *ProductHandler) GetProductsHandler(c *fiber.Ctx) error {

	defer func() {

		h.metric.ObserveHandler("ProductHandler_GetProducts", time.Now(), c.Response().StatusCode())

	}()

	products, err := h.productService.GetProducts()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{
			Error:   "Ürünler Getirilemedi",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{
		SuccesData: products,
	})

}

func (h *ProductHandler) GetProductHandler(c *fiber.Ctx) error {

	defer func() {

		h.metric.ObserveHandler("ProductHandler_GetProduct", time.Now(), c.Response().StatusCode())

	}()

	randomNumber := rand.Intn(100)

	if randomNumber < 40 {
		return c.Status(fiber.StatusServiceUnavailable).JSON(models.FailResponse{
			Error:   "Sunucu  hizmet veremiyor",
			Details: "Tekrar deneyin",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.FailResponse{
			Error:   "ID Parse Hatası",
			Details: err.Error(),
		})
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.FailResponse{
			Error:   "Ürün Bulunamadı",
			Details: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SuccesResponse{
		SuccesData: product,
	})

}

func (h *ProductHandler) SetRoutesProduct(app *fiber.App) {

	productGroup := app.Group("/products")

	productGroup.Get("", h.GetProductsHandler)

	productGroup.Get("/:id", h.GetProductHandler)

}

func ProductGetEndpoints() []*endpoint.EndPoint {

	return []*endpoint.EndPoint{

		endpoint.New(
			endpoint.GET,
			"/products",
			endpoint.WithTags("products"),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccesResponse{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{response.New(models.FailResponse{}, "500", "Internal Server")}),
		),

		endpoint.New(
			endpoint.GET,
			"/products/{id}",
			endpoint.WithTags("products"),
			endpoint.WithParams(parameter.IntParam("id", parameter.Path, parameter.WithRequired())),
			endpoint.WithSuccessfulReturns([]response.Response{response.New(models.SuccesResponse{}, "200", "OK")}),
			endpoint.WithErrors([]response.Response{
				response.New(models.FailResponse{}, "400", "Bad Request"),
				response.New(models.FailResponse{}, "500", "Internal Server")}),
		),
	}

}
