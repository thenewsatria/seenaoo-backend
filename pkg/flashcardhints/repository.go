package flashcardhints

import (
	"context"
	"errors"
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error)
	PopulateFlashcard(fId *models.ReadFlashcardRequest) (*[]models.FlashcardHint, error)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}

func (r *repository) PopulateFlashcard(fId *models.ReadFlashcardRequest) (*[]models.FlashcardHint, error) {
	var hints = []models.FlashcardHint{}
	flashcardId, err := primitive.ObjectIDFromHex(fId.ID)
	if err != nil {
		return nil, errors.New("Invalid flashcard_id")
	}
	cursor, err := r.Collection.Find(database.GetDBContext(), bson.D{{Key: "flashcard_id", Value: flashcardId}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var hint = models.FlashcardHint{}
		err := cursor.Decode(&hint)
		if err != nil {
			return nil, err
		}
		hints = append(hints, hint)
	}

	return &hints, nil
}

func (r *repository) CreateFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error) {
	fh.ID = primitive.NewObjectID()
	fh.CreatedAt = time.Now()
	fh.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), fh)
	if err != nil {
		return nil, err
	}

	return fh, nil
}
