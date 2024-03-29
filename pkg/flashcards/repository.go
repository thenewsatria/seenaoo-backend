package flashcards

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateFlashcard(f *models.Flashcard) (*models.Flashcard, error, bool)
	ReadFlashcard(fId *models.FlashcardByIdRequest) (*models.Flashcard, error)
	ReadFlashcardsByFlashcardCoverId(fFCoverId *models.FlashcardCoverById) (*[]models.Flashcard, error)
	UpdateFlashcard(f *models.Flashcard) (*models.Flashcard, error, bool)
	DeleteFlashcard(f *models.Flashcard) (*models.Flashcard, error)
	DeleteFlashcardsByFlashcardCoverId(fcCoverId *models.FlashcardCoverById) (int64, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateFlashcard(f *models.Flashcard) (*models.Flashcard, error, bool) {
	f.ID = primitive.NewObjectID()
	f.CreatedAt = time.Now()
	f.UpdatedAt = time.Now()

	err := validator.ValidateStruct(f)
	if err != nil {
		if validator.IsValidationError(err) {
			return nil, err, true
		}
		return nil, err, false
	}

	_, err = r.Collection.InsertOne(database.GetDBContext(), f)
	if err != nil {
		return nil, err, false
	}
	return f, nil, false
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

func (r *repository) UpdateFlashcard(f *models.Flashcard) (*models.Flashcard, error, bool) {
	f.UpdatedAt = time.Now()

	err := validator.ValidateStruct(f)
	if err != nil {
		if validator.IsValidationError(err) {
			return nil, err, true
		}
		return nil, err, false
	}

	_, err = r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": f.ID}, bson.M{"$set": f})
	if err != nil {
		return nil, err, false
	}
	return f, nil, false
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

func (r *repository) DeleteFlashcardsByFlashcardCoverId(fcCoverId *models.FlashcardCoverById) (int64, error) {
	fcCvrId, err := primitive.ObjectIDFromHex(fcCoverId.ID)
	if err != nil {
		return -1, err
	}
	res, err := r.Collection.DeleteMany(database.GetDBContext(), bson.M{"flashcard_cover_id": fcCvrId})
	if err != nil {
		return -1, err
	}

	return res.DeletedCount, nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
