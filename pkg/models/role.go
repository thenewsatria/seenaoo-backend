package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID   `bson:"_id" json:"id" validate:"required"`
	Owner       string               `bson:"owner" json:"owner" validate:"required,lowercase,alphanum,max=25,min=5"`
	Name        string               `bson:"name" json:"name" validate:"required,max=64,min=5"`
	Slug        string               `bson:"slug" json:"slug" valodate:"required"`
	Description string               `bson:"description" json:"description" validate:"omitempty,min=5,max=255"`
	Permissions []primitive.ObjectID `bson:"permissions" json:"permissions" validate:"required,dive,required"`
	CreatedAt   time.Time            `bson:"createdAt" json:"createdAt" validate:"required"`
	UpdatedAt   time.Time            `bson:"updatedAt" json:"updatedAt" validate:"required"`
}

type RoleById struct {
	ID string `bson:"_id" json:"id" validate:"required"`
}

type RoleByOwner struct {
	Owner string `bson:"owner" json:"owner" validate:"required,lowercase,alphanum,max=25,min=5"`
}

type RoleBySlug struct {
	Slug string `bson:"slug" json:"slug" validate:"required"`
}
