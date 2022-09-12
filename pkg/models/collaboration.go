package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collaboration struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	Inviter          string             `bson:"inviter" json:"inviter"`
	Collaborator     string             `bson:"collaborator" json:"collaborator"`
	ItemID           primitive.ObjectID `bson:"item_id" json:"itemId"`
	ItemType         string             `bson:"item_type" json:"itemType"`
	Status           string             `bson:"status" json:"status"` //rejected, sent, accepted
	RoleAttachmentId primitive.ObjectID `bson:"roleAttachmentId" json:"roleAttachmentId"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt"`
}

type CollaborationById struct {
	ID string `bson:"_id" json:"id"`
}

type CollaborationByItemIdAndCollaborator struct {
	ItemID       string `bson:"item_id" json:"itemId"`
	Collaborator string `bson:"collaborator" json:"collaborator"`
}
