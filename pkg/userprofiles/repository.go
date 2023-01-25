package userprofiles

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
	CreateProfile(up *models.UserProfile) (*models.UserProfile, error, bool)
	ReadProfileById(upId *models.UserProfileByID) (*models.UserProfile, error)
	ReadProfileByOwner(upOwner *models.UserProfileByOwner) (*models.UserProfile, error)
	UpdateProfile(up *models.UserProfile) (*models.UserProfile, error, bool)
}

type repository struct {
	Collection *mongo.Collection
}

// CreateProfile implements Repository
func (r *repository) CreateProfile(up *models.UserProfile) (*models.UserProfile, error, bool) {
	up.ID = primitive.NewObjectID()
	up.IsVerified = false
	up.CreatedAt = time.Now()
	up.UpdatedAt = time.Now()

	err := validator.ValidateStruct(up)
	if err != nil {
		if validator.IsValidationError(err) {
			err = validator.TranslateError(err)
			return nil, err, true
		}

		return nil, err, false
	}

	_, err = r.Collection.InsertOne(database.GetDBContext(), up)
	if err != nil {
		return nil, err, false
	}

	return up, nil, false
}

// ReadProfileById implements Repository
func (r *repository) ReadProfileById(upId *models.UserProfileByID) (*models.UserProfile, error) {
	up := &models.UserProfile{}
	profileId, err := primitive.ObjectIDFromHex(upId.ID)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "_id", Value: profileId}}).Decode(up)
	if err != nil {
		return nil, err
	}

	return up, nil
}

// ReadProfileByOwner implements Repository
func (r *repository) ReadProfileByOwner(upOwner *models.UserProfileByOwner) (*models.UserProfile, error) {
	up := &models.UserProfile{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "owner", Value: upOwner.Owner}}).Decode(up)
	if err != nil {
		return nil, err
	}

	return up, nil
}

func (r *repository) UpdateProfile(up *models.UserProfile) (*models.UserProfile, error, bool) {
	up.UpdatedAt = time.Now()

	err := validator.ValidateStruct(up)
	if err != nil {
		if validator.IsValidationError(err) {
			err = validator.TranslateError(err)
			return nil, err, true
		}
		return nil, err, false
	}

	_, err = r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": up.ID}, bson.M{"$set": up})
	if err != nil {
		return nil, err, false
	}

	return up, nil, false
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		Collection: c,
	}
}
