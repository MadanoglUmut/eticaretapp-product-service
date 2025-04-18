package main

import (
	"fmt"
	"log"
	"os"

	"ProductService/internal/handlers"
	"ProductService/internal/repositories"
	"ProductService/internal/services"
	"ProductService/metrics"
	"ProductService/pkg/psql"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
)

func main() {

	fmt.Println("Hello World")

	app := fiber.New()

	//err := godotenv.Load("../../.env")

	err := godotenv.Load(".env")

	if err != nil {

		log.Fatal("Env Dosyasi Okunamadi", err)

	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	fmt.Printf("%s, %s, %s, %s, %s", host, port, user, password, name)

	var db = psql.Connect(host, user, password, name, port)

	productRepository := repositories.NewProductRepository(db)

	productService := services.NewProductRepository(productRepository)

	histogram := metrics.NewNamedHistogram("http_request_ProductService_duration_seconds", []float64{0.001, 0.005, 0.01, 0.05, 0.1})

	registry := prometheus.NewRegistry()
	registry.MustRegister(histogram.Histogram)

	productHandler := handlers.NewProductHandler(productService, histogram)

	productHandler.SetRoutesProduct(app)

	app.Get("/metrics", adaptor.HTTPHandler(metrics.GetHandler(registry)))

	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

	sw.AddEndpoints(handlers.ProductGetEndpoints())

	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	log.Fatal(app.Listen(":5050"))

}
