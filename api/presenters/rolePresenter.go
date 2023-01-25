package presenters

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id"`
	Owner       string               `bson:"owner" json:"owner"`
	Name        string               `bson:"name" json:"name"`
	Slug        string               `bson:"slug" json:"slug"`
	Description string               `bson:"description" json:"description"`
	Permissions []primitive.ObjectID `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time            `bson:"updatedAt" json:"updatedAt"`
}

type RoleDetail struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Owner       User               `bson:"owner" json:"owner"`
	Name        string             `bson:"name" json:"name"`
	Slug        string             `bson:"slug" json:"slug"`
	Description string             `bson:"description" json:"description"`
	Permissions []Permission       `bson:"permissions" json:"permissions"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func RoleSuccessResponse(r *models.Role) *fiber.Map {
	role := &Role{
		ID:          r.ID,
		Owner:       r.Owner,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		Permissions: r.Permissions,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.CreatedAt,
	}
	return &fiber.Map{
		"success": true,
		"data":    role,
		"error":   nil,
	}
}

func RoleDetailSuccessResponse(r *models.Role, u *models.User, p *[]models.Permission) *fiber.Map {
	usr := &User{
		Username: u.Username,
		// DisplayName:     u.DisplayName,
		// AvatarImagePath: u.AvatarImagePath,
		// Biography:       u.Biography,
		// IsVerified:      u.IsVerified,
	}
	roleDetail := &RoleDetail{
		ID:          r.ID,
		Owner:       *usr,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		Permissions: []Permission{},
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	for _, perm := range *p {
		permit := &Permission{
			ID:           perm.ID,
			ItemCategory: perm.ItemCategory,
			Name:         perm.Name,
			Description:  perm.Description,
			CreatedAt:    perm.CreatedAt,
			UpdatedAt:    perm.UpdatedAt,
		}
		roleDetail.Permissions = append(roleDetail.Permissions, *permit)
	}
	return &fiber.Map{
		"success": true,
		"data":    roleDetail,
		"error":   nil,
	}
}

func RolesSuccessResponse(rs *[]models.Role) *fiber.Map {
	var roles = []Role{}
	for _, r := range *rs {
		role := &Role{
			ID:          r.ID,
			Owner:       r.Owner,
			Name:        r.Name,
			Slug:        r.Slug,
			Description: r.Description,
			Permissions: r.Permissions,
			CreatedAt:   r.CreatedAt,
			UpdatedAt:   r.CreatedAt,
		}
		roles = append(roles, *role)
	}
	return &fiber.Map{
		"success": true,
		"data":    roles,
		"error":   nil,
	}
}
