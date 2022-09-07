package flashcards

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error)
	FetchFlashcard(flashcardId *models.FlashcardByIdRequest) (*models.Flashcard, error)
	UpdateFlashcard(f *models.Flashcard) (*models.Flashcard, error)
	RemoveFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error)
	PopulateFlashcardCover(flashcardCoverId *models.FlashcardCoverById) (*[]models.Flashcard, error)
}

type service struct {
	repository Repository
}

// InsertFlashcard implements Service
func (s *service) InsertFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error) {
	return s.repository.CreateFlashcard(flashcard)
}

func (s *service) FetchFlashcard(flashcardId *models.FlashcardByIdRequest) (*models.Flashcard, error) {
	return s.repository.ReadFlashcard(flashcardId)
}

func (s *service) UpdateFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error) {
	return s.repository.UpdateFlashcard(flashcard)
}

func (s *service) RemoveFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error) {
	return s.repository.DeleteFlashcard(flashcard)
}

func (s *service) PopulateFlashcardCover(flashcardCoverId *models.FlashcardCoverById) (*[]models.Flashcard, error) {
	return s.repository.ReadFlashcardsByFlashcardCoverId(flashcardCoverId)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
