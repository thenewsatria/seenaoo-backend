package models

import (
	"time"

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

type PermissionByItemCategory struct {
	ItemCategory string `bson:"item_category" json:"itemCategory"`
}

type PermissionById struct {
	ID string `bson:"_id" json:"id"`
}

type PermissionByName struct {
	Name string `bson:"name" json:"name"`
}
