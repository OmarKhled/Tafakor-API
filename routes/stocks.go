package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"tafakor.app/controllers"
)

func StocksRoutes(group fiber.Router, db *sql.DB) {
	group.Get("/", func(c *fiber.Ctx) error {
		postID := c.QueryInt("postID")
		fmt.Println("ID", postID)
		stocks := controllers.GetStocks(db, postID)

		return c.JSON(stocks)
	})
}
