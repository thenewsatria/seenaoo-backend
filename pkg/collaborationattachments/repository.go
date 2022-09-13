package collaborationattachments

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateCollabAttachment(ca *models.CollaborationAttachment) (*models.CollaborationAttachment, error)
	ReadCollabAttachment(caId *models.CollaborationAttachmentById) (*models.CollaborationAttachment, error)
	DeleteCollabAttachment(ca *models.CollaborationAttachment) (*models.CollaborationAttachment, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateCollabAttachment(ca *models.CollaborationAttachment) (*models.CollaborationAttachment, error) {
	ca.ID = primitive.NewObjectID()
	ca.CreatedAt = time.Now()
	ca.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), ca)
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (r *repository) DeleteCollabAttachment(ca *models.CollaborationAttachment) (*models.CollaborationAttachment, error) {
	_, err := r.Collection.DeleteOne(database.GetDBContext(), bson.M{"_id": ca.ID})
	if err != nil {
		return nil, err
	}
	return ca, nil
}

func (r *repository) ReadCollabAttachment(caId *models.CollaborationAttachmentById) (*models.CollaborationAttachment, error) {
	var collabAttachment = &models.CollaborationAttachment{}
	id, err := primitive.ObjectIDFromHex(caId.ID)
	if err != nil {
		return nil, err
	}
	err = r.Collection.FindOne(database.GetDBContext(), bson.M{"_id": id}).Decode(collabAttachment)
	if err != nil {
		return nil, err
	}

	return collabAttachment, err
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		Collection: c,
	}
}
