package models

type LoginRequest struct {
	Credential string `bson:"credential" json:"credential"`
	Password   string `bson:"password" json:"password"`
}
