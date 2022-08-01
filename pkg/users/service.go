package users

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertUser(user *models.User) (*models.User, error)
}

type service struct {
	repository Repository
}

func (s *service) InsertUser(user *models.User) (*models.User, error) {
	return s.repository.CreateUser(user)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
