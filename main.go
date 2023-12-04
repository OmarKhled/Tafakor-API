package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"
	"tafakor.app/config"
	"tafakor.app/routes"
)

func main() {
	// Loading enviroment
	config.LoadEnv()

	// Initiating DB
	db := config.DBConfig()

	// Initiating Fiber
	app := fiber.New()

	// Disabling CORS
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(`
     _____      __       _                   _    ____ ___ 
    |_   _|_ _ / _| __ _| | _____  _ __     / \  |  _ \_ _|
      | |/ _  | |_ / _  | |/ / _ \| '__|   / _ \ | |_) | | 
      | | (_| |  _| (_| |   < (_) | |     / ___ \|  __/| | 
      |_|\__,_|_|  \__,_|_|\_\___/|_|    /_/   \_\_|  |___|

      Developed by: Omar Khaled ":"`)
	})

	// Verses Routes
	routes.VersesRoutes(app.Group("/verses"), db)

	log.Fatal(app.Listen(":8080"))
}
