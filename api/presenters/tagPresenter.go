package presenters

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ItemList struct {
	FlashcardCovers []FlashcardCover `bson:"flashcards" json:"flashcards"`
}

type Tag struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	TagName string             `bson:"tag_name" json:"tagName"`
}

type TagDetails struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	TagName string             `bson:"tag_name" json:"tagName"`
	Items   ItemList           `bson:"items" json:"items"`
}

func TagSuccessResponse(t *models.Tag) *fiber.Map {
	tag := &Tag{
		ID:      t.ID,
		TagName: t.TagName,
	}
	return &fiber.Map{
		"success": true,
		"data":    tag,
		"error":   nil,
	}
}

func TagDetailSuccessResponse(t *models.Tag, fcs *[]models.FlashcardCover) *fiber.Map {
	itemList := &ItemList{
		FlashcardCovers: []FlashcardCover{},
	}

	for _, fc := range *fcs {
		fcCover := &FlashcardCover{
			ID:          fc.ID,
			Slug:        fc.Slug,
			Title:       fc.Title,
			Description: fc.Description,
			Image_path:  fc.Image_path,
			Tags:        fc.Tags,
			Author:      fc.Author,
			CreatedAt:   fc.CreatedAt,
			UpdatedAt:   fc.UpdatedAt,
		}

		itemList.FlashcardCovers = append(itemList.FlashcardCovers, *fcCover)

	}

	tag := &TagDetails{
		ID:      t.ID,
		TagName: t.TagName,
		Items:   *itemList,
	}

	return &fiber.Map{
		"success": true,
		"data":    tag,
		"error":   nil,
	}
}
