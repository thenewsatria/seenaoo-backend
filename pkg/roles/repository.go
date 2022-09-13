package roles

import (
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateRole(r *models.Role) (*models.Role, error)
	ReadRoleById(rId *models.RoleById) (*models.Role, error)
	ReadRoleBySlug(rSlug *models.RoleBySlug) (*models.Role, error)
	ReadRolesByOwner(rOwner *models.RoleByOwner) (*[]models.Role, error)
	UpdateRole(r *models.Role) (*models.Role, error)
	DeleteRole(r *models.Role) (*models.Role, error)
}

type repository struct {
	Collection *mongo.Collection
}

func (r *repository) CreateRole(rl *models.Role) (*models.Role, error) {
	rl.ID = primitive.NewObjectID()
	rl.CreatedAt = time.Now()
	rl.UpdatedAt = time.Now()

	_, err := r.Collection.InsertOne(database.GetDBContext(), rl)
	if err != nil {
		return nil, err
	}

	return rl, nil
}

func (r *repository) DeleteRole(rl *models.Role) (*models.Role, error) {
	_, err := r.Collection.DeleteOne(database.GetDBContext(), bson.M{"_id": rl.ID})
	if err != nil {
		return nil, err
	}
	return rl, nil
}

func (r *repository) ReadRoleById(rlId *models.RoleById) (*models.Role, error) {
	var role = &models.Role{}
	roleId, err := primitive.ObjectIDFromHex(rlId.ID)
	if err != nil {
		return nil, err
	}
	err = r.Collection.FindOne(database.GetDBContext(), bson.M{"_id": roleId}).Decode(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *repository) ReadRoleBySlug(rlSlug *models.RoleBySlug) (*models.Role, error) {
	var role = &models.Role{}
	err := r.Collection.FindOne(database.GetDBContext(), bson.M{"slug": rlSlug.Slug}).Decode(role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (r *repository) ReadRolesByOwner(rlOwner *models.RoleByOwner) (*[]models.Role, error) {
	var roles = []models.Role{}
	cursor, err := r.Collection.Find(database.GetDBContext(), bson.M{"owner": rlOwner.Owner})
	if err != nil {
		return nil, err
	}

	for cursor.Next(database.GetDBContext()) {
		var role = &models.Role{}
		err := cursor.Decode(role)
		if err != nil {
			return nil, err
		}
		roles = append(roles, *role)
	}

	return &roles, nil
}

func (r *repository) UpdateRole(rl *models.Role) (*models.Role, error) {
	rl.UpdatedAt = time.Now()
	_, err := r.Collection.UpdateOne(database.GetDBContext(), bson.M{"_id": rl.ID}, bson.M{"$set": rl})
	if err != nil {
		return nil, err
	}

	return rl, nil
}

func NewRepo(c *mongo.Collection) Repository {
	return &repository{
		Collection: c,
	}
}
