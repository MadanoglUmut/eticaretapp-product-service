package handlers

import (
	"ProductService/internal/repositories"
	"ProductService/internal/services"
	"ProductService/pkg/psql"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

var app *fiber.App = fiber.New()

type TestDB struct {
	DB        *gorm.DB
	Container testcontainers.Container
}

type HandlerSetup struct {
	DB          *gorm.DB
	App         *fiber.App
	MockMetrics *MockMetrics
}

func (h *HandlerSetup) SetupProductHandler() {

	productRepository := repositories.NewProductRepository(h.DB)
	productService := services.NewProductRepository(productRepository)
	productHandler := NewProductHandler(productService, h.MockMetrics)
	productHandler.SetRoutesProduct(h.App)

}

func (t *TestDB) Setup() error {

	ctx := context.Background()

	dbConfig := map[string]string{
		"POSTGRES_USER":     "user",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_DB":       "users",
	}

	defaultPort := nat.Port("5432/tcp")
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{defaultPort.Port()},
			Env:          dbConfig,
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort(defaultPort),
			),
		},
		Started: true,
	})
	if err != nil {
		return err
	}
	t.Container = container
	port, err := container.MappedPort(ctx, defaultPort)
	if err != nil {
		return err
	}
	fmt.Println("Veritabanı başladı port numarasi:", port)
	t.DB = psql.Connect("0.0.0.0", dbConfig["POSTGRES_USER"], dbConfig["POSTGRES_PASSWORD"], dbConfig["POSTGRES_DB"], port.Port())
	return t.loadSQLFiles()
}

func (t *TestDB) loadSQLFiles() error {
	fileCreate, err := os.ReadFile("../../psql/create_tables.sql")
	if err != nil {
		return err
	}
	if err := t.DB.Exec(string(fileCreate)).Error; err != nil {
		return err
	}

	fileFill, err := os.ReadFile("../../psql/fill_tables.sql")
	if err != nil {
		return err
	}
	return t.DB.Exec(string(fileFill)).Error
}

func (t *TestDB) CleanUp() {
	t.Container.Terminate(context.Background())
}

type MockMetrics struct{}

func (m *MockMetrics) ObserveHandler(name string, startTime time.Time, status int) {

}

func TestMain(m *testing.M) {
	testDB := &TestDB{}
	if err := testDB.Setup(); err != nil {
		fmt.Println("Veritabanı bağlantısı başarısız", err)
		os.Exit(1)
	}
	defer testDB.CleanUp()

	handlerSetup := &HandlerSetup{
		DB:          testDB.DB,
		App:         app,
		MockMetrics: &MockMetrics{},
	}

	handlerSetup.SetupProductHandler()

	os.Exit(m.Run())
}
