package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	controllers "tafakor.app/controllers"
	"tafakor.app/models"
)

func PublishRoutes(group fiber.Router, db *sql.DB) {
	group.Post("/", func(c *fiber.Ctx) error {
		// Approval Reuest Email
		controllers.RequestPostApproval(c.Body())

		return c.JSON(true)
	})

	// ===================================
	// 					Acceptance Routes
	// ===================================

	group.Get("/accept", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.PublishmentParamaters
		c.QueryParser(&parameters)

		// Publish post to FB
		status, id := controllers.PublishToFB(parameters.PostingType, parameters.FileURL)

		if status == true { // Published Successfully

			// Saving Post & Stock data to DB
			postID := controllers.RecordPost(db, parameters.VerseID, true, "accepted", parameters.FileURL, id)
			controllers.RecordStock(db, parameters.StockID, postID, parameters.StockProvider, "accepted")

			// Response
			return c.JSON(postID)
		} else {
			return c.JSON(nil)
		}
	})

	// ===================================
	// 					Rejection Routes
	// ===================================

	reject := group.Group("/reject")

	reject.Get("/", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.PublishmentParamaters
		c.QueryParser(&parameters)

		// Saving Post & Stock data to DB
		postID := controllers.RecordPost(db, parameters.VerseID, true, "rejected", parameters.FileURL, "")
		controllers.RecordStock(db, parameters.StockID, postID, parameters.StockProvider, "rejected")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(postID)
	})

	reject.Get("/verse", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.PublishmentParamaters
		c.QueryParser(&parameters)

		// Saving Post data to DB
		postID := controllers.RecordPost(db, parameters.VerseID, true, "rejected", parameters.FileURL, "")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(postID)
	})

	reject.Get("/stock", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.PublishmentParamaters
		c.QueryParser(&parameters)

		// Saving Stock data to DB
		stockId := controllers.RecordStock(db, parameters.StockID, "", parameters.StockProvider, "rejected")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(stockId)
	})

	reject.Get("/stock-post", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.PublishmentParamaters
		c.QueryParser(&parameters)

		// Saving Post & Stock data to DB
		postID := controllers.RecordPost(db, parameters.VerseID, false, "pending", parameters.FileURL, "")
		stockId := controllers.RecordStock(db, parameters.StockID, postID, parameters.StockProvider, "rejected-once")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(stockId)
	})

}
