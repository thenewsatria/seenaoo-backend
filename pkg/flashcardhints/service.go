package flashcardhints

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error)
}

type service struct {
	repository Repository
}

// InsertFlashcardHint implements Service
func (s *service) InsertFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error) {
	return s.repository.CreateFlashcardHint(flashcardHint)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
