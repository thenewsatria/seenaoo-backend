package presenters

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlashcardCover struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Slug        string               `bson:"slug" json:"slug"`
	Title       string               `bson:"title" json:"title"`
	Description string               `bson:"description" json:"description"`
	ImagePath   string               `bson:"image_path" json:"imagePath"`
	Tags        []primitive.ObjectID `bson:"tags" json:"tags"`
	Author      string               `bson:"author" json:"author"`
	CreatedAt   time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updated_at" json:"updatedAt"`
}

type FlashcardCoverDetail struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Slug        string             `bson:"slug" json:"slug"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	ImagePath   string             `bson:"image_path" json:"imagePath"`
	Tags        []Tag              `bson:"tags" json:"tags"`
	Flashcards  []Flashcard        `bson:"flashcards" json:"flashcards"`
	Author      UserDetail         `bson:"author" json:"author"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

func FlashcardCoverSuccessResponse(fcCover *models.FlashcardCover) *fiber.Map {
	flashcardCvr := &FlashcardCover{
		ID:          fcCover.ID,
		Slug:        fcCover.Slug,
		Title:       fcCover.Title,
		Description: fcCover.Description,
		ImagePath:   fcCover.ImagePath,
		Tags:        fcCover.Tags,
		Author:      fcCover.Author,
		CreatedAt:   fcCover.CreatedAt,
		UpdatedAt:   fcCover.UpdatedAt,
	}
	return &fiber.Map{
		"success": true,
		"data":    flashcardCvr,
		"error":   nil,
	}
}

func FlashcardCoverDetailSuccessResponse(fcCover *models.FlashcardCover, tags *[]models.Tag,
	flashcards *[]models.Flashcard, author *models.User, authorProfile *models.UserProfile) *fiber.Map {

	userProfile := &UserProfile{
		DisplayName:     authorProfile.DisplayName,
		AvatarImagePath: authorProfile.AvatarImagePath,
		BannerImagePath: authorProfile.BannerImagePath,
		Biography:       authorProfile.Biography,
		IsVerified:      authorProfile.IsVerified,
	}

	owner := &UserDetail{
		Username: author.Username,
		Email:    author.Email,
		Profile:  *userProfile,
		// DisplayName:     author.DisplayName,
		// AvatarImagePath: author.AvatarImagePath,
		// Biography:       author.Biography,
		// IsVerified:      author.IsVerified,
	}

	flashcardCvrDetail := &FlashcardCoverDetail{
		ID:          fcCover.ID,
		Slug:        fcCover.Slug,
		Title:       fcCover.Title,
		Description: fcCover.Description,
		ImagePath:   fcCover.ImagePath,
		Tags:        []Tag{},
		Flashcards:  []Flashcard{},
		Author:      *owner,
		CreatedAt:   fcCover.CreatedAt,
		UpdatedAt:   fcCover.UpdatedAt,
	}

	for _, tag := range *tags {
		tagDetail := &Tag{
			ID:      tag.ID,
			TagName: tag.TagName,
		}
		flashcardCvrDetail.Tags = append(flashcardCvrDetail.Tags, *tagDetail)
	}

	for _, flashcard := range *flashcards {
		flashcardDetail := &Flashcard{
			ID:               flashcard.ID,
			FrontImagePath:   flashcard.FrontImagePath,
			BackImagePath:    flashcard.BackImagePath,
			FrontText:        flashcard.FrontText,
			BackText:         flashcard.BackImagePath,
			Question:         flashcard.Question,
			FlashCardCoverId: flashcard.FlashCardCoverId,
			CreatedAt:        flashcard.CreatedAt,
			UpdatedAt:        flashcard.UpdatedAt,
		}
		flashcardCvrDetail.Flashcards = append(flashcardCvrDetail.Flashcards, *flashcardDetail)
	}

	return &fiber.Map{
		"success": true,
		"data":    flashcardCvrDetail,
		"error":   nil,
	}
}
