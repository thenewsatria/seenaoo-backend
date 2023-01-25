package tests

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertTest(test *models.Test) (*models.Test, error)
}

type service struct {
	repository Repository
}

// InsertTest implements Service
func (s *service) InsertTest(test *models.Test) (*models.Test, error) {
	panic("unimplemented")
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
