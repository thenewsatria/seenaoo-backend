package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func Router(app *fiber.App) {
	var flashcardCollection = database.UseDB().Collection(flashcards.CollectionName)
	var flashcardRepo = flashcards.NewRepo(flashcardCollection)
	var flashcardService = flashcards.NewService(flashcardRepo)

	var flashcardHintCollection = database.UseDB().Collection(flashcardhints.CollectionName)
	var flashcardHintRepo = flashcardhints.NewRepo(flashcardHintCollection)
	var flashcardHintService = flashcardhints.NewService(flashcardHintRepo)

	var userCollection = database.UseDB().Collection(users.CollectionName)
	var userRepo = users.NewRepo(userCollection)
	var userService = users.NewService(userRepo)

	msg := os.Getenv("FIBER_ENV")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  200,
			"message": msg,
		})
	})

	api := app.Group("/api")
	apiV1 := api.Group("/v1")

	flashcardRouter(apiV1, flashcardService, flashcardHintService)
	flashcardHintRouter(apiV1, flashcardHintService)
	userRouter(apiV1, userService)
}
