package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/refreshtokens"
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

	var refreshTokenCollection = database.UseDB().Collection(refreshtokens.CollectionName)
	var refreshTokenRepo = refreshtokens.NewRepo(refreshTokenCollection)
	var refreshTokenService = refreshtokens.NewService(refreshTokenRepo)

	api := app.Group("/api")
	apiV1 := api.Group("/v1")

	flashcardRouter(apiV1, flashcardService, flashcardHintService)
	flashcardHintRouter(apiV1, flashcardHintService, flashcardService)
	authenticationRouter(apiV1, userService, refreshTokenService)
}
