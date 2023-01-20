package flashcardhints

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error, bool)
	FetchFlashcardHint(flashcardHintId *models.FlashcardHintByIdRequest) (*models.FlashcardHint, error)
	PopulateFlashcard(flashcardId *models.FlashcardByIdRequest) (*[]models.FlashcardHint, error)
	UpdateFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error, bool)
	RemoveFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error)
	RemoveFlashcardHintsByFlashcardId(flashcardId *models.FlashcardByIdRequest) (int64, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) InsertFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error, bool) {
	return s.repository.CreateFlashcardHint(flashcardHint)
}

func (s *service) FetchFlashcardHint(flashcardHintId *models.FlashcardHintByIdRequest) (*models.FlashcardHint, error) {
	return s.repository.ReadFlashcardHint(flashcardHintId)
}

func (s *service) PopulateFlashcard(flashcardId *models.FlashcardByIdRequest) (*[]models.FlashcardHint, error) {
	return s.repository.ReadFlashcardHintsByFlashcardId(flashcardId)
}

func (s *service) UpdateFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error, bool) {
	return s.repository.UpdateFlashcardHint(flashcardHint)
}

func (s *service) RemoveFlashcardHint(flashcardHint *models.FlashcardHint) (*models.FlashcardHint, error) {
	return s.repository.DeleteFlashcardHint(flashcardHint)
}

func (s *service) RemoveFlashcardHintsByFlashcardId(flashcardId *models.FlashcardByIdRequest) (int64, error) {
	return s.repository.DeleteFlashcardHintsByFlashcardId(flashcardId)
}
