package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"tafakor.app/config"
	verses "tafakor.app/controllers"
)

func main() {
	// Loading enviroment
	config.LoadEnv()

	// Initiating DB
	db := config.DBConfig()

	// Initiating Fiber
	app := fiber.New()

	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Tafakor API!")
	})

	app.Get("/verse", verses.GetVerse(db))

	log.Fatal(app.Listen(":8080"))
}
