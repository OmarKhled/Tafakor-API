package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"tafakor.app/config"
	verses "tafakor.app/controllers"
)

func main() {
	// Loading enviroment
	config.LoadEnv()

	// Initiating DB
	config.DBConfig()

	// Initiating Fiber
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Tafakor API!")
	})

	app.Get("/verses", verses.GetVerses)

	log.Fatal(app.Listen(":3000"))
}
