package flashcardcovers

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateFlashcardCover(fcCover *models.FlashcardCover) (*models.FlashcardCover, error)
	ReadFlashcardCoverBySlug(fcCoverSlug *models.FlashcardCoverBySlug) (*models.FlashcardCover, error)
	ReadFlashcardCoverById(fcCoverId *models.FlashcardCoverById) (*models.FlashcardCover, error)
	UpdateFlashcardCover(fcCover *models.FlashcardCover) (*models.FlashcardCover, error)
	DeleteFlashcardCover(fcCover *models.FlashcardCover) (*models.FlashcardCover, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateFlashcardCover(fcCover *models.FlashcardCover) (*models.FlashcardCover, error) {
	fcCover.ID = primitive.NewObjectID()
	fcCover.CreatedAt = time.Now()
	fcCover.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), fcCover)
	if err != nil {
		return nil, err
	}

	return fcCover, nil
}

func (r *repository) DeleteFlashcardCover(fcCover *models.FlashcardCover) (*models.FlashcardCover, error) {
	_, err := r.Collection.DeleteOne(database.GetDBContext(), bson.M{"_id": fcCover.ID})
	if err != nil {
		return nil, err
	}

	return fcCover, nil
}

func (r *repository) ReadFlashcardCoverById(fcCoverId *models.FlashcardCoverById) (*models.FlashcardCover, error) {
	flashcardCover := &models.FlashcardCover{}
	fccId, err := primitive.ObjectIDFromHex(fcCoverId.ID)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "_id", Value: fccId}}).Decode(flashcardCover)
	if err != nil {
		return nil, err
	}

	return flashcardCover, nil
}

func (r *repository) ReadFlashcardCoverBySlug(fcCoverSlug *models.FlashcardCoverBySlug) (*models.FlashcardCover, error) {
	fcCover := &models.FlashcardCover{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "slug", Value: fcCoverSlug.Slug}}).Decode(fcCover)
	if err != nil {
		return nil, err
	}

	return fcCover, nil
}

func (r *repository) UpdateFlashcardCover(fcCover *models.FlashcardCover) (*models.FlashcardCover, error) {
	fcCover.UpdatedAt = time.Now()
	_, err := r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": fcCover.ID}, bson.M{"$set": fcCover})
	if err != nil {
		return nil, err
	}

	return fcCover, nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
