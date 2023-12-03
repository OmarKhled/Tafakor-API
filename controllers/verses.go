package verses

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"tafakor.app/migrate"
)

func GetVerses(db *sql.DB) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		migrate.Seed(db)

		return c.JSON(migrate.FetchPosts())
	}

}
