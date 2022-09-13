package models

import (
	"time"

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

type RoleById struct {
	ID string `bson:"_id" json:"id"`
}

type RoleByOwner struct {
	Owner string `bson:"owner" json:"owner"`
}

type RoleBySlug struct {
	Slug string `bson:"slug" json:"slug"`
}
