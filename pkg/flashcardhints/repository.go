package flashcardhints

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) CreateFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error) {
	flashcardHint.ID = primitive.NewObjectID()
	flashcardHint.CreatedAt = time.Now()
	flashcardHint.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), flashcardHint)
	if err != nil {
		return nil, err
	}

	return flashcardHint, nil
}
