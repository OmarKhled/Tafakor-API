package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"tafakor.app/controllers"
)

func StocksRoutes(group fiber.Router, db *sql.DB) {
	group.Get("/", func(c *fiber.Ctx) error {
		verseId := c.QueryInt("verseId")
		fmt.Println("ID", verseId)
		stocks := controllers.GetStocks(db, verseId)

		return c.JSON(stocks)
	})
}
