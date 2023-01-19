package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	ID           primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	ItemCategory string             `bson:"item_category" json:"itemCategory" validate:"required,uppercase,oneof=FLASHCARD"`
	Name         string             `bson:"name" json:"name" validate:"required,uppercase,min=5,max=64"`
	Description  string             `bson:"description" json:"description" validate:"omitempty,min=5,max=255"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type PermissionByItemCategory struct {
	ItemCategory string `bson:"item_category" json:"itemCategory" validate:"required,uppercase,oneof=FLASHCARD"`
}

type PermissionById struct {
	ID string `bson:"_id" json:"id" validate:"required"`
}

type PermissionByName struct {
	Name string `bson:"name" json:"name" validate:"required,uppercase,min=5,max=64"`
}
