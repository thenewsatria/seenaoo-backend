package models

type Authentication struct {
	Credential string `bson:"credential" json:"credential"`
	Password   string `bson:"password" json:"password"`
}
