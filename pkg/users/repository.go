package users

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateUser(u *models.User) (*models.User, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateUser(u *models.User) (*models.User, error) {
	u.ID = primitive.NewObjectID()
	u.IsVerified = false
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
