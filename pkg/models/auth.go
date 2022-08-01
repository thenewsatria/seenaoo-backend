package models

type Authentication struct {
	Credential string `json:"credential"`
	Password   string `json:"password"`
}
