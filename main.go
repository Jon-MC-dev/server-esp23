package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"serveresp32.com/sensor/model"
)

func main() {
	// Conectar a la base de datos MySQL
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := "root@tcp(127.0.0.1:3306)/db_esp32?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos")
	}
	// defer db.Close()

	// Migrar el esquema, esto creará la tabla si no existe
	db.AutoMigrate(&model.DataModel{})

	// Crear una instancia de Fiber
	app := fiber.New()

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

	// Iniciar el servidor en el puerto 8080
	log.Fatal(app.Listen(":8080"))
}
