package tests

import (
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateTest(*models.Test) (*models.Test, error)
}

type repository struct {
	collection *mongo.Collection
}

func (r *repository) CreateTest(testData *models.Test) (*models.Test, error) {
	panic("unimplemented")
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		collection: c,
	}
}
