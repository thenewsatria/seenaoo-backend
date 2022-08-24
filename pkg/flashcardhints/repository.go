package flashcardhints

import (
	"context"
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error)
	ReadFlashcardHint(fhId *models.FlashcardHintByIdRequest) (*models.FlashcardHint, error)
	ReadFlashcardHintsByFlashcardId(fId *models.FlashcardByIdRequest) (*[]models.FlashcardHint, error)
	UpdateFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error)
	DeleteFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error)
	DeleteFlashcardHintsByFlashcardId(fId *models.FlashcardByIdRequest) (int64, error)
}

type repository struct {
	Collection *mongo.Collection
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
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

func (r *repository) ReadFlashcardHint(fhId *models.FlashcardHintByIdRequest) (*models.FlashcardHint, error) {
	flashcardHint := &models.FlashcardHint{}
	id, err := primitive.ObjectIDFromHex(fhId.ID)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "_id", Value: id}}).Decode(flashcardHint)
	if err != nil {
		return nil, err
	}

	return flashcardHint, nil
}

func (r *repository) ReadFlashcardHintsByFlashcardId(fId *models.FlashcardByIdRequest) (*[]models.FlashcardHint, error) {
	var hints = []models.FlashcardHint{}
	flashcardId, err := primitive.ObjectIDFromHex(fId.ID)
	if err != nil {
		return nil, err
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

func (r *repository) UpdateFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error) {
	fh.UpdatedAt = time.Now()
	_, err := r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": fh.ID}, bson.M{"$set": fh})
	if err != nil {
		return nil, err
	}

	return fh, nil
}

func (r *repository) DeleteFlashcardHint(fh *models.FlashcardHint) (*models.FlashcardHint, error) {
	_, err := r.Collection.DeleteOne(database.GetDBContext(), bson.M{"_id": fh.ID})
	if err != nil {
		return nil, err
	}

	return fh, nil
}

func (r *repository) DeleteFlashcardHintsByFlashcardId(fId *models.FlashcardByIdRequest) (int64, error) {
	id, err := primitive.ObjectIDFromHex(fId.ID)
	if err != nil {
		return -1, err
	}
	res, err := r.Collection.DeleteMany(database.GetDBContext(), bson.M{"flashcard_id": id})
	if err != nil {
		return -1, err
	}
	return res.DeletedCount, nil
}
