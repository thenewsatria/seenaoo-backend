package flashcards

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateFlashcard(f *models.Flashcard) (*models.Flashcard, error)
	ReadFlashcard(fId *models.ReadFlashcardRequest) (*models.Flashcard, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateFlashcard(f *models.Flashcard) (*models.Flashcard, error) {
	f.ID = primitive.NewObjectID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *repository) ReadFlashcard(fId *models.ReadFlashcardRequest) (*models.Flashcard, error) {
	flashcard := &models.Flashcard{}
	id, err := primitive.ObjectIDFromHex(fId.ID)
	if err != nil {
		return nil, err
	}
	err = r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "_id", Value: id}}).Decode(flashcard)
	if err != nil {
		return nil, err
	}

	return flashcard, nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
