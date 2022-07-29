package flashcards

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error)
	FetchFlashcard(flashcardId *models.ReadFlashcardRequest) (*models.Flashcard, error)
}

type service struct {
	repository Repository
}

// InsertFlashcard implements Service
func (s *service) InsertFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error) {
	return s.repository.CreateFlashcard(flashcard)
}

func (s *service) FetchFlashcard(flashcardId *models.ReadFlashcardRequest) (*models.Flashcard, error) {
	return s.repository.ReadFlashcard(flashcardId)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
