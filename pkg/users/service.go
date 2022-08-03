package users

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertUser(user *models.User) (*models.User, error)
	CheckEmailIsUnique(userEmail *models.UserByEmailRequest) bool
	CheckUsernameIsUnique(userUsername *models.UserByUsernameRequest) bool
	FetchUserByEmail(userEmail *models.UserByEmailRequest) (*models.User, error)
	FetchUSerByUsername(userName *models.UserByUsernameRequest) (*models.User, error)
}

type service struct {
	repository Repository
}

func (s *service) InsertUser(user *models.User) (*models.User, error) {
	return s.repository.CreateUser(user)
}

func (s *service) CheckEmailIsUnique(userEmail *models.UserByEmailRequest) bool {
	user, _ := s.repository.GetUserByEmail(userEmail)
	return user == nil
}

func (s *service) CheckUsernameIsUnique(userUsername *models.UserByUsernameRequest) bool {
	user, _ := s.repository.GetUserByUsername(userUsername)
	return user == nil
}

func (s *service) FetchUserByEmail(userEmail *models.UserByEmailRequest) (*models.User, error) {
	user, err := s.repository.GetUserByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) FetchUSerByUsername(userUsername *models.UserByUsernameRequest) (*models.User, error) {
	user, err := s.repository.GetUserByUsername(userUsername)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
