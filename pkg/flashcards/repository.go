package flashcards

import (
	"errors"
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error)
	// 	ReadFlashcards() ([]*models.Flashcard, error)
	ReadFlashcard(flashcardId *models.ReadFlashcardRequest) (*models.Flashcard, error)
	// 	UpdateFlashcard(flashcardId *models.ReadFlashcardRequest, flashcard *models.Flashcard) (*models.Flashcard, error)
	// 	DeleteFlashcard(flashcardId *models.DeleteFlashcardRequest) (*models.Flashcard, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error) {
	flashcard.ID = primitive.NewObjectID()
	flashcard.CreatedAt = time.Now()
	flashcard.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), flashcard)
	if err != nil {
		return nil, err
	}
	return flashcard, nil
}

func (r *repository) ReadFlashcard(flashcardId *models.ReadFlashcardRequest) (*models.Flashcard, error) {
	flashcard := &models.Flashcard{}
	id, err := primitive.ObjectIDFromHex(flashcardId.ID)
	if err != nil {
		return nil, errors.New("Invalid flashcard id")
	}
	err = r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "_id", Value: id}}).Decode(flashcard)
	if err != nil {
		return nil, err
	}

	return flashcard, nil
}

// // ReadFlashcards implements FlashcardRepository
// func (*flashcardRepository) ReadFlashcards() ([]*models.Flashcard, error) {
// 	panic("unimplemented")
// }

// // UpdateFlashcard implements FlashcardRepository
// func (*flashcardRepository) UpdateFlashcard(flashcardId *models.ReadFlashcardRequest, flashcard *models.Flashcard) (*models.Flashcard, error) {
// 	panic("unimplemented")
// }

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
