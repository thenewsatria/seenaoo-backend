package collaborations

import (
	"context"
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateCollaboration(cl *models.Collaboration) (*models.Collaboration, error, bool)
	ReadCollaboration(clId *models.CollaborationById) (*models.Collaboration, error)
	UpdateCollaboration(cl *models.Collaboration) (*models.Collaboration, error, bool)
	DeleteCollaboration(cl *models.Collaboration) (*models.Collaboration, error)
	ReadCollaborationsByItemIdAndCollaborator(clIdAndC *models.CollaborationByItemIdAndCollaborator) (*models.Collaboration, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) UpdateCollaboration(cl *models.Collaboration) (*models.Collaboration, error, bool) {
	cl.UpdatedAt = time.Now()

	err := validator.ValidateStruct(cl)
	if err != nil {
		if validator.IsValidationError(err) {
			return nil, err, true
		}
		return nil, err, false
	}

	_, err = r.Collection.UpdateOne(context.Background(), bson.M{"_id": cl.ID}, bson.M{"$set": cl})
	if err != nil {
		return nil, err, false
	}

	return cl, nil, false
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

func (r *repository) CreateCollaboration(cl *models.Collaboration) (*models.Collaboration, error, bool) {
	cl.ID = primitive.NewObjectID()
	cl.CreatedAt = time.Now()
	cl.UpdatedAt = time.Now()
	cl.Status = "SENT"

	err := validator.ValidateStruct(cl)
	if err != nil {
		if validator.IsValidationError(err) {
			return nil, err, true
		}
		return nil, err, false
	}

	_, err = r.Collection.InsertOne(database.GetDBContext(), cl)

	if err != nil {
		return nil, err, false
	}

	return cl, nil, false
}

func (r *repository) ReadCollaborationsByItemIdAndCollaborator(clIdAndC *models.CollaborationByItemIdAndCollaborator) (*models.Collaboration, error) {
	collaboration := &models.Collaboration{}
	itemId, err := primitive.ObjectIDFromHex(clIdAndC.ItemID)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOne(database.GetDBContext(), bson.M{"item_id": itemId, "collaborator": clIdAndC.Collaborator}).Decode(collaboration)
	if err != nil {
		return nil, err
	}

	return collaboration, err
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		Collection: c,
	}
}
