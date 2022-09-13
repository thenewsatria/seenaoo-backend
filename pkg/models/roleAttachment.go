package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CollaborationAttachment struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Attacher  string             `bson:"attacher" json:"attacher"`
	Attached  string             `bson:"attached" json:"attached"`
	Role      string             `bson:"role" json:"role"`
	ItemID    primitive.ObjectID `bson:"item_id" json:"itemId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

type CollaborationAttachmentById struct {
	ID string `bson:"_id" json:"id"`
}
