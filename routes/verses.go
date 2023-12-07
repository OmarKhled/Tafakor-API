package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	controllers "tafakor.app/controllers"
)

func VersesRoutes(group fiber.Router, db *sql.DB) {
	group.Get("/", func(c *fiber.Ctx) error {
		verses := controllers.GetVerses(db)
		return c.JSON(verses)
	})

	group.Get("/one", func(c *fiber.Ctx) error {
		random := c.QueryBool("random")
		verse := controllers.GetVerse(random, db)

		return c.JSON(verse)
	})
}
