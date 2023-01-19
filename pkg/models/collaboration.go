package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collaboration struct {
	ID           primitive.ObjectID `bson:"_id" json:"id" validate:"required"`
	Inviter      string             `bson:"inviter" json:"inviter" validate:"required,lowercase,alphanum,max=25,min=5"`
	Collaborator string             `bson:"collaborator" json:"collaborator" validate:"required,lowercase,alphanum,max=25,min=5"`
	ItemID       primitive.ObjectID `bson:"item_id" json:"itemId" validate:"required"`
	ItemType     string             `bson:"item_type" json:"itemType" validate:"required,oneof=FLASHCARD"`
	Status       string             `bson:"status" json:"status" validate:"required,oneof=REJECTED SENT ACCEPTED"` //rejected, sent, accepted
	RoleId       primitive.ObjectID `bson:"role_id" json:"role_id" validate:"required"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt" validate:"required"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt" validate:"required"`
}

type CollaborationById struct {
	ID string `bson:"_id" json:"id" validate:"required"`
}

type CollaborationByItemIdAndCollaborator struct {
	ItemID       string `bson:"item_id" json:"itemId" validate:"required"`
	Collaborator string `bson:"collaborator" json:"collaborator" validate:"required,lowercase,alphanum,max=25,min=5"`
}

type CollaborationStatusRequest struct {
	Status string `bson:"status" json:"status" validate:"required,oneof=rejected sent accepted"`
}
