package permissions

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
	CreatePermission(p *models.Permission) (*models.Permission, error, bool)
	ReadPermissionsByItemCategory(pItemCat *models.PermissionByItemCategory) (*[]models.Permission, error)
	ReadPermissionById(pId *models.PermissionById) (*models.Permission, error)
	ReadAllPermissions() (*[]models.Permission, error)
	ReadPermissionsDistinctItemCategory() (*[]string, error)
	ReadPermissionByName(pName *models.PermissionByName) (*models.Permission, error)
	UpdatePermission(p *models.Permission) (*models.Permission, error, bool)
	DeletePermission(p *models.Permission) (*models.Permission, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreatePermission(p *models.Permission) (*models.Permission, error, bool) {
	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	err := validator.ValidateStruct(p)
	if err != nil {
		return nil, err, true
	}

	_, err = r.Collection.InsertOne(database.GetDBContext(), p)
	if err != nil {
		return nil, err, false
	}
	return p, nil, false
}

func (r *repository) DeletePermission(p *models.Permission) (*models.Permission, error) {
	_, err := r.Collection.DeleteOne(database.GetDBContext(), bson.M{"_id": p.ID})
	if err != nil {
		return nil, err
	}
	return p, err
}

func (r *repository) ReadPermissionById(pId *models.PermissionById) (*models.Permission, error) {
	var permission = &models.Permission{}
	permId, err := primitive.ObjectIDFromHex(pId.ID)
	if err != nil {
		return nil, err
	}

	err = r.Collection.FindOne(database.GetDBContext(), bson.M{"_id": permId}).Decode(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (r *repository) ReadPermissionByName(pName *models.PermissionByName) (*models.Permission, error) {
	var permission = &models.Permission{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.M{"name": pName.Name}).Decode(permission)
	if err != nil {
		return nil, err
	}
	return permission, err
}

func (r *repository) ReadAllPermissions() (*[]models.Permission, error) {
	var permissions = []models.Permission{}
	cursor, err := r.Collection.Find(database.GetDBContext(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(database.GetDBContext()) {
		var permission = &models.Permission{}
		err := cursor.Decode(permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, *permission)
	}

	return &permissions, nil
}

func (r *repository) ReadPermissionsDistinctItemCategory() (*[]string, error) {
	uniqueItemCategories := []string{}
	result, err := r.Collection.Distinct(database.GetDBContext(), "item_category", bson.M{})
	if err != nil {
		return nil, err
	}

	for _, category := range result {
		itemCat := category.(string)
		uniqueItemCategories = append(uniqueItemCategories, itemCat)
	}
	return &uniqueItemCategories, nil
}

func (r *repository) ReadPermissionsByItemCategory(pItemCat *models.PermissionByItemCategory) (*[]models.Permission, error) {
	var permissions = []models.Permission{}
	cursor, err := r.Collection.Find(database.GetDBContext(), bson.M{"item_category": pItemCat.ItemCategory})
	if err != nil {
		return nil, err
	}

	for cursor.Next(database.GetDBContext()) {
		var permission = &models.Permission{}
		err := cursor.Decode(permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, *permission)
	}

	return &permissions, nil
}

func (r *repository) UpdatePermission(p *models.Permission) (*models.Permission, error, bool) {
	p.UpdatedAt = time.Now()

	err := validator.ValidateStruct(p)
	if err != nil {
		return nil, err, true
	}

	_, err = r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": p.ID}, bson.M{"$set": p})
	if err != nil {
		return nil, err, false
	}

	return p, nil, false
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		Collection: c,
	}
}
