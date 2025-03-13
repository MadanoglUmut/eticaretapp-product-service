package handlers

import (
	"ProductService/internal/models"

	"github.com/go-swagno/swagno/components/endpoint"
	"github.com/go-swagno/swagno/components/http/response"
	"github.com/go-swagno/swagno/components/parameter"
	"github.com/gofiber/fiber/v2"
)

type productService interface {
	GetProducts() ([]models.Product, error)
	GetProduct(id int) (models.Product, error)
}

type ProductHandler struct {
	productService productService
}

func NewProductHandler(productService productService) *ProductHandler {

	return &ProductHandler{
		productService: productService,
	}

}

func (h *ProductHandler) GetProductsHandler(c *fiber.Ctx) error {

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
	/*
		randomNumber := rand.Intn(100)


		if randomNumber < 55 {
			return c.Status(fiber.StatusServiceUnavailable).JSON(models.FailResponse{
				Error:   "Sunucu  hizmet veremiyor",
				Details: "Tekrar deneyin",
			})
		}
	*/
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
