package presenters

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	ItemCategory string             `bson:"item_category" json:"itemCategory"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description" json:"description"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

func PermissionsGroupedByItemCateogory(groupedPermissions map[string]*[]models.Permission) *fiber.Map {
	permits := map[string]*[]Permission{}
	for key, permissions := range groupedPermissions {
		currentCategory := []Permission{}
		for _, permission := range *permissions {
			permit := Permission{
				ID:           permission.ID,
				ItemCategory: permission.ItemCategory,
				Name:         permission.Name,
				Description:  permission.Description,
				CreatedAt:    permission.CreatedAt,
				UpdatedAt:    permission.CreatedAt,
			}
			currentCategory = append(currentCategory, permit)
		}
		permits[key] = &currentCategory

	}
	return &fiber.Map{
		"status": "success",
		"data":   permits,
	}
}
