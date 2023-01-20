package tags

import (
	"strings"
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateTag(t *models.Tag) (*models.Tag, error, bool)
	ReadTagById(tId *models.TagById) (*models.Tag, error)
	ReadTagByTagName(tName *models.TagByName) (*models.Tag, error)
	DeleteTag(t *models.Tag) (*models.Tag, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateTag(t *models.Tag) (*models.Tag, error, bool) {
	t.ID = primitive.NewObjectID()
	t.TagName = strings.ToLower(t.TagName) //mengecilkan seluruh huruf sebelum dilakukan validasi
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()

	err := validator.ValidateStruct(t)
	if err != nil {
		if validator.IsValidationError(err) {
			err = validator.TranslateError(err)
			return nil, err, true
		}
		return nil, err, false
	}

	_, err = r.Collection.InsertOne(database.GetDBContext(), t)
	if err != nil {
		return nil, err, false
	}

	return t, nil, false
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
	tagName := strings.ToLower(tName.TagName)
	err := r.Collection.FindOne(database.GetDBContext(), bson.D{{Key: "tag_name", Value: tagName}}).Decode(tag)
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
