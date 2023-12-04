package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	versesController "tafakor.app/controllers"
)

func VersesRoutes(group fiber.Router, db *sql.DB) {
	group.Get("/", func(c *fiber.Ctx) error {
		verses := versesController.GetVerses(db)

		return c.JSON(verses)
	})

	group.Get("/one", func(c *fiber.Ctx) error {
		random := c.QueryBool("random")
		verse := versesController.GetVerse(random, db)

		return c.JSON(verse)
	})
}
