package main

import (
	"fmt"
	"log"
	"os"

	"ProductService/internal/handlers"
	"ProductService/internal/repositories"
	"ProductService/internal/services"
	"ProductService/pkg/psql"

	"github.com/go-swagno/swagno"
	"github.com/go-swagno/swagno-fiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	fmt.Println("Hello World")

	err := godotenv.Load("../../.env")

	if err != nil {

		log.Fatal("Env Dosyasi Okunamadi", err)

	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	var db = psql.Connect(host, user, password, name, port)

	productRepository := repositories.NewProductRepository(db)

	productService := services.NewProductRepository(productRepository)

	productHandler := handlers.NewProductHandler(productService)

	app := fiber.New()

	productHandler.SetRoutesProduct(app)

	sw := swagno.New(swagno.Config{Title: "Testing API", Version: "v1.0.0"})

	sw.AddEndpoints(handlers.ProductGetEndpoints())

	swagger.SwaggerHandler(app, sw.MustToJson(), swagger.WithPrefix("/swagger"))

	log.Fatal(app.Listen(":5050"))

}

//RETRY MEKANİZMASI ÖNCELİKLİ LİST TEN YİNE GET İSTEĞİ ATCAZ TOKEN VAR MI KONTROL ETCEZ SONRA İŞLEMLERE DEVAM
