package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	controllers "tafakor.app/controllers"
	"tafakor.app/models"
)

func PublishRoutes(group fiber.Router, db *sql.DB) {
	group.Post("/", func(c *fiber.Ctx) error {
		var body models.PublishmentParamaters
		json.Unmarshal(c.Body(), &body)

		// Record Posting as Pending
		postId := controllers.RecordPost(db, body.VerseID, false, "pending", body.PostURL, 0, body.ReelURL)
		fmt.Println("Post id", postId)
		controllers.RecordStock(db, body.StockID, postId, body.StockProvider, "pending")
		fmt.Println("stock")

		// Approval Reuest Email
		controllers.RequestPostApproval(postId, body.PostURL, body.ReelURL)
		fmt.Println("verify")

		return c.JSON(postId)
	})

	// ===================================
	// 					Acceptance Routes
	// ===================================

	group.Get("/accept", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.EmailSubmissionParameters
		c.QueryParser(&parameters)

		var post models.Post = controllers.GetPost(db, parameters.PostID)
		fmt.Println("Post retrival done")
		var stock models.Stock = controllers.GetStockByPostID(db, parameters.PostID)
		fmt.Println("Stock retrival done")

		// Publish post to FB
		status, _ := controllers.SocialPublishment(post.PostURL, post.ReelURL, parameters.Platform)

		// status := true

		if status == true { // Published Successfully

			// Accepting Post saved in DB
			controllers.ChangePostStatus(db, post.ID, "accepted")

			controllers.ChangeStockStatus(db, stock.ID, "accepted")

			// Response
			return c.JSON(parameters.PostID)
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
		var parameters models.EmailSubmissionParameters
		c.QueryParser(&parameters)

		// Retreiving pending stock
		var stock models.Stock = controllers.GetStockByPostID(db, parameters.PostID)

		// Rejecting Post & Stock
		controllers.ChangePostStatus(db, parameters.PostID, "rejected")
		controllers.ChangeStockStatus(db, stock.ID, "rejected")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(parameters.PostID)
	})

	reject.Get("/verse", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.EmailSubmissionParameters
		c.QueryParser(&parameters)

		// Retreiving pending stock
		var stock models.Stock = controllers.GetStockByPostID(db, parameters.PostID)

		// Saving Post data to DB
		controllers.ChangePostStatus(db, parameters.PostID, "rejected")
		controllers.ChangeStockStatus(db, stock.ID, "discarded")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(parameters.PostID)
	})

	reject.Get("/stock", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.EmailSubmissionParameters
		c.QueryParser(&parameters)

		// Retreiving pending stock
		var stock models.Stock = controllers.GetStockByPostID(db, parameters.PostID)

		// Saving Stock data to DB
		controllers.ChangePostStatus(db, parameters.PostID, "discarded")
		controllers.ChangeStockStatus(db, stock.ID, "rejected")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(stock.ID)
	})

	reject.Get("/stock-post", func(c *fiber.Ctx) error {
		// Extracting request query
		var parameters models.EmailSubmissionParameters
		c.QueryParser(&parameters)

		// Retreiving pending stock
		var stock models.Stock = controllers.GetStockByPostID(db, parameters.PostID)

		// Saving Stock data to DB
		controllers.ChangePostStatus(db, parameters.PostID, "discarded")
		controllers.ChangeStockStatus(db, stock.ID, "rejected-once")

		// Trigger new render workflow
		controllers.TriggerRender()

		// Response
		return c.JSON(stock.ID)
	})

}
