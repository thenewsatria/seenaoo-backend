package presenters

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collaboration struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Inviter      string             `bson:"inviter" json:"inviter"`
	Collaborator string             `bson:"collaborator" json:"collaborator"`
	ItemID       primitive.ObjectID `bson:"item_id" json:"itemId"`
	ItemType     string             `bson:"item_type" json:"itemType"`
	Status       string             `bson:"status" json:"status"`
	RoleId       primitive.ObjectID `bson:"roleId" json:"roleId"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

type CollaborationFlashcardCoverDetail struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Inviter      UserDetail         `bson:"inviter" json:"inviter"`
	Collaborator UserDetail         `bson:"collaborator" json:"collaborator"`
	Item         FlashcardCover     `bson:"item" json:"item"`
	ItemType     string             `bson:"item_type" json:"itemType"`
	Status       string             `bson:"status" json:"status"`
	Role         Role               `bson:"role" json:"role"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

func CollaborationSuccessResponse(collaboration *models.Collaboration) *fiber.Map {
	collab := &Collaboration{
		ID:           collaboration.ID,
		Inviter:      collaboration.Inviter,
		Collaborator: collaboration.Collaborator,
		ItemID:       collaboration.ItemID,
		ItemType:     collaboration.ItemType,
		Status:       collaboration.Status,
		CreatedAt:    collaboration.CreatedAt,
		UpdatedAt:    collaboration.UpdatedAt,
	}

	return &fiber.Map{
		"status": true,
		"data":   collab,
		"error":  nil,
	}
}

func CollaborationFlashcardDetailSuccessResponse(collaboration *models.Collaboration, inviter *models.User,
	collaborator *models.User, inviterProfile *models.UserProfile, collaboratorProfile *models.UserProfile,
	flashcardCover *models.FlashcardCover, r *models.Role) *fiber.Map {

	invUserProfile := &UserProfile{
		DisplayName:     inviterProfile.DisplayName,
		AvatarImagePath: inviterProfile.AvatarImagePath,
		BannerImagePath: inviterProfile.BannerImagePath,
		Biography:       inviterProfile.Biography,
		IsVerified:      inviterProfile.IsVerified,
	}

	invUser := &UserDetail{
		Username: inviter.Username,
		Profile:  *invUserProfile,
	}

	collabUserProfile := &UserProfile{
		DisplayName:     collaboratorProfile.DisplayName,
		AvatarImagePath: collaboratorProfile.AvatarImagePath,
		BannerImagePath: collaboratorProfile.BannerImagePath,
		Biography:       collaboratorProfile.Biography,
		IsVerified:      collaboratorProfile.IsVerified,
	}

	collabUser := &UserDetail{
		Username: collaborator.Username,
		Profile:  *collabUserProfile,
	}

	fcCover := &FlashcardCover{
		ID:          flashcardCover.ID,
		Slug:        flashcardCover.Slug,
		Title:       flashcardCover.Title,
		Description: flashcardCover.Description,
		ImagePath:   flashcardCover.ImagePath,
		Tags:        flashcardCover.Tags,
		Author:      flashcardCover.Author,
		CreatedAt:   flashcardCover.CreatedAt,
		UpdatedAt:   flashcardCover.UpdatedAt,
	}

	role := &Role{
		ID:          r.ID,
		Owner:       r.Owner,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		Permissions: r.Permissions,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}

	collabFlashcardCoverDetail := &CollaborationFlashcardCoverDetail{
		ID:           collaboration.ID,
		Inviter:      *invUser,
		Collaborator: *collabUser,
		Item:         *fcCover,
		ItemType:     collaboration.ItemType,
		Status:       collaboration.Status,
		Role:         *role,
		CreatedAt:    collaboration.CreatedAt,
		UpdatedAt:    collaboration.UpdatedAt,
	}

	return &fiber.Map{
		"status": "success",
		"data":   collabFlashcardCoverDetail,
		"error":  nil,
	}
}
