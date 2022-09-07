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
	ReadFlashcard(fId *models.FlashcardByIdRequest) (*models.Flashcard, error)
	ReadFlashcardsByFlashcardCoverId(fFCoverId *models.FlashcardCoverById) (*[]models.Flashcard, error)
	UpdateFlashcard(f *models.Flashcard) (*models.Flashcard, error)
	DeleteFlashcard(f *models.Flashcard) (*models.Flashcard, error)
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

func (r *repository) ReadFlashcard(fId *models.FlashcardByIdRequest) (*models.Flashcard, error) {
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

func (r *repository) UpdateFlashcard(f *models.Flashcard) (*models.Flashcard, error) {
	f.UpdatedAt = time.Now()
	_, err := r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": f.ID}, bson.M{"$set": f})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *repository) DeleteFlashcard(f *models.Flashcard) (*models.Flashcard, error) {
	_, err := r.Collection.DeleteOne(database.GetDBContext(), bson.M{"_id": f.ID})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (r *repository) ReadFlashcardsByFlashcardCoverId(fFCoverId *models.FlashcardCoverById) (*[]models.Flashcard, error) {
	flashcards := []models.Flashcard{}
	flashCvrId, err := primitive.ObjectIDFromHex(fFCoverId.ID)
	if err != nil {
		return nil, err
	}
	cursor, err := r.Collection.Find(database.GetDBContext(), bson.D{{Key: "flashcard_cover_id", Value: flashCvrId}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(database.GetDBContext()) {
		var flashcard = models.Flashcard{}
		err := cursor.Decode(&flashcard)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}

	return &flashcards, nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
