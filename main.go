package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"serveresp32.com/sensor/model"
)

func main() {
	// Conectar a la base de datos MySQL
	// connectionString := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	connectionString := "root@tcp(127.0.0.1:3306)/db_esp32?charset=utf8mb4&parseTime=True&loc=Local"
	if os.Getenv("DB_USER") != "" {
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbName := os.Getenv("DB_NAME")

		connectionString = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbUser,
			dbPassword,
			dbHost,
			dbPort,
			dbName,
		)
	}

	fmt.Println(connectionString)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos")
	}
	// defer db.Close()

	// Migrar el esquema, esto creará la tabla si no existe
	db.AutoMigrate(&model.DataModel{})

	// Crear una instancia de Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Agregar middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// Manejador para la ruta POST que recibe los datos y los guarda en la base de datos
	app.Post("/guardar", func(c *fiber.Ctx) error {
		data := new(model.DataModel)
		if err := c.BodyParser(data); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Error al analizar la solicitud",
			})
		}

		// Crear una nueva entrada en la base de datos
		data.Timestamp = time.Now()
		db.Create(data)

		return c.JSON(fiber.Map{
			"message": "Datos guardados en la base de datos",
		})
	})

	app.Get("/consultar", func(c *fiber.Ctx) error {
		var datos []model.DataModel
		db.Find(&datos)

		return c.JSON(datos)
	})

	app.Get("/consultar/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var datos []model.DataModel
		if err := db.Where("id >= ?", id).Find(&datos).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al consultar la base de datos",
			})
		}

		return c.JSON(datos)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
