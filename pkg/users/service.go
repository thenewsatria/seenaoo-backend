package users

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertUser(user *models.User) (*models.User, error, bool)
	CheckEmailIsExist(userEmail *models.UserByEmailRequest) bool
	CheckUsernameIsExist(userUsername *models.UserByUsernameRequest) bool
	FetchUserByEmail(userEmail *models.UserByEmailRequest) (*models.User, error)
	FetchUserByUsername(userName *models.UserByUsernameRequest) (*models.User, error)
}

type service struct {
	repository Repository
}

func (s *service) InsertUser(user *models.User) (*models.User, error, bool) {
	return s.repository.CreateUser(user)
}

func (s *service) CheckEmailIsExist(userEmail *models.UserByEmailRequest) bool {
	user, _ := s.repository.ReadUserByEmail(userEmail)
	return user != nil
}

func (s *service) CheckUsernameIsExist(userUsername *models.UserByUsernameRequest) bool {
	user, _ := s.repository.ReadUserByUsername(userUsername)
	return user != nil
}

func (s *service) FetchUserByEmail(userEmail *models.UserByEmailRequest) (*models.User, error) {
	user, err := s.repository.ReadUserByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) FetchUserByUsername(userUsername *models.UserByUsernameRequest) (*models.User, error) {
	user, err := s.repository.ReadUserByUsername(userUsername)
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
