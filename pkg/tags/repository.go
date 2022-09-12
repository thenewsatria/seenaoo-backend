package tags

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateTag(t *models.Tag) (*models.Tag, error)
	ReadTagById(tId *models.TagById) (*models.Tag, error)
	ReadTagByTagName(tName *models.TagByName) (*models.Tag, error)
	DeleteTag(t *models.Tag) (*models.Tag, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateTag(t *models.Tag) (*models.Tag, error) {
	t.ID = primitive.NewObjectID()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), t)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *repository) ReadTagById(tId *models.TagById) (*models.Tag, error) {
	tag := &models.Tag{}
	tagId, err := primitive.ObjectIDFromHex(tId.ID)
	if err != nil {
		return nil, err
	}
	err = r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "_id", Value: tagId}}).Decode(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *repository) ReadTagByTagName(tName *models.TagByName) (*models.Tag, error) {
	tag := &models.Tag{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "tag_name", Value: tName.TagName}}).Decode(tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

func (r *repository) DeleteTag(t *models.Tag) (*models.Tag, error) {
	_, err := r.Collection.DeleteOne(database.GetDBContext(), bson.M{"_id": t.ID})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		Collection: c,
	}
}
