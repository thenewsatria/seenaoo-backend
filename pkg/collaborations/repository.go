package collaborations

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
	CreateCollaboration(cl *models.Collaboration) (*models.Collaboration, error)
	ReadCollaboration(clId *models.CollaborationById) (*models.Collaboration, error)
	UpdateCollaboration(cl *models.Collaboration) (*models.Collaboration, error)
	DeleteCollaboration(cl *models.Collaboration) (*models.Collaboration, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) UpdateCollaboration(cl *models.Collaboration) (*models.Collaboration, error) {
	cl.UpdatedAt = time.Now()

	_, err := r.Collection.UpdateOne(context.Background(), bson.M{"_id": cl.ID}, bson.M{"$set": cl})
	if err != nil {
		return nil, err
	}

	return cl, err
}

func (r *repository) DeleteCollaboration(cl *models.Collaboration) (*models.Collaboration, error) {
	_, err := r.Collection.DeleteOne(context.Background(), bson.M{"_id": cl.ID})
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (r *repository) ReadCollaboration(clId *models.CollaborationById) (*models.Collaboration, error) {
	collaboration := &models.Collaboration{}
	id, err := primitive.ObjectIDFromHex(clId.ID)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "_id", Value: id}}).Decode(collaboration)

	if err != nil {
		return nil, err
	}

	return collaboration, err
}

func (r *repository) CreateCollaboration(cl *models.Collaboration) (*models.Collaboration, error) {
	cl.ID = primitive.NewObjectID()
	cl.CreatedAt = time.Now()
	cl.UpdatedAt = time.Now()
	cl.Status = "SENT"

	_, err := r.Collection.InsertOne(database.GetDBContext(), cl)

	if err != nil {
		return nil, err
	}

	return cl, nil
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		Collection: c,
	}
}
