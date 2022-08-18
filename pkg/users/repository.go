package users

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateUser(u *models.User) (*models.User, error)
	ReadUserByEmail(uEmail *models.UserByEmailRequest) (*models.User, error)
	ReadUserByUsername(uUname *models.UserByUsernameRequest) (*models.User, error)
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

func (r *repository) ReadUserByEmail(uEmail *models.UserByEmailRequest) (*models.User, error) {
	user := &models.User{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "email", Value: uEmail.Email}}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) ReadUserByUsername(uUname *models.UserByUsernameRequest) (*models.User, error) {
	user := &models.User{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "username", Value: uUname.Username}}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func NewRepo(collection *mongo.Collection) Repository {
	return &repository{
		Collection: collection,
	}
}
