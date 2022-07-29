package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
)

var CreateFlashcard = func(c *fiber.Ctx) error {
	collection := database.UseDB().Collection(models.FlashCardCollectionName)
	newFlashcard := new(models.Flashcard)

	if err := c.BodyParser(newFlashcard); err != nil {
		return c.JSON(fiber.Map{
			"status":     "error",
			"statusCode": 400,
			"message":    err.Error(),
		})
	}

	_, err := collection.InsertOne(database.GetDBContext(), newFlashcard)
	if err != nil {
		return c.JSON(fiber.Map{
			"status":     "error",
			"statusCode": 400,
			"message":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "item created",
	})
}
