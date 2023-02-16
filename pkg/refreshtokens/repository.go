package refreshtokens

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
	CreateRefreshToken(rt *models.RefreshToken) (*models.RefreshToken, error, bool)
	ReadRefreshTokenByUsersUsername(rtu *models.RefreshTokenByUsersUsernameRequest) (*models.RefreshToken, error)
	UpdateRefreshToken(rt *models.RefreshToken) (*models.RefreshToken, error, bool)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateRefreshToken(rt *models.RefreshToken) (*models.RefreshToken, error, bool) {
	rt.ID = primitive.NewObjectID()
	rt.CreatedAt = time.Now()
	rt.UpdatedAt = time.Now()
	rt.IsBlocked = false

	err := validator.ValidateStruct(rt)
	if err != nil {
		if validator.IsValidationError(err) {
			return nil, err, true
		}
		return nil, err, false
	}

	_, err = r.Collection.InsertOne(database.GetDBContext(), rt)
	if err != nil {
		return nil, err, false
	}

	return rt, nil, false
}

func (r *repository) ReadRefreshTokenByUsersUsername(rtu *models.RefreshTokenByUsersUsernameRequest) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "username", Value: rtu.Username}}).Decode(refreshToken)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}

func (r *repository) UpdateRefreshToken(rt *models.RefreshToken) (*models.RefreshToken, error, bool) {
	rt.UpdatedAt = time.Now()

	err := validator.ValidateStruct(rt)
	if err != nil {
		if validator.IsValidationError(err) {
			return nil, err, true
		}
		return nil, err, false
	}
	_, err = r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": rt.ID}, bson.M{"$set": rt})
	if err != nil {
		return nil, err, false
	}
	return rt, nil, false
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
