package flashcardcovers

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type Service interface {
	InsertFlashcardCover(flashcardCover *models.FlashcardCover) (*models.FlashcardCover, error, bool)
	RemoveFlashcardCover(flashcardCover *models.FlashcardCover) (*models.FlashcardCover, error)
	FetchFlashcardCoverById(flashcardCoverId *models.FlashcardCoverById) (*models.FlashcardCover, error)
	FetchFlashcardCoverBySlug(flashcardCoverSlug *models.FlashcardCoverBySlug) (*models.FlashcardCover, error)
	UpdateFlashcardCover(flashcardCover *models.FlashcardCover) (*models.FlashcardCover, error, bool)
	FetchFlashcardCoversByTagId(tagId *models.TagById) (*[]models.FlashcardCover, error)
}

type service struct {
	repository Repository
}

func (s *service) InsertFlashcardCover(flashcardCover *models.FlashcardCover) (*models.FlashcardCover, error, bool) {
	return s.repository.CreateFlashcardCover(flashcardCover)
}

func (s *service) RemoveFlashcardCover(flashcardCover *models.FlashcardCover) (*models.FlashcardCover, error) {
	return s.repository.DeleteFlashcardCover(flashcardCover)
}

func (s *service) FetchFlashcardCoverById(flashcardCoverId *models.FlashcardCoverById) (*models.FlashcardCover, error) {
	return s.repository.ReadFlashcardCoverById(flashcardCoverId)
}

func (s *service) FetchFlashcardCoverBySlug(flashcardCoverSlug *models.FlashcardCoverBySlug) (*models.FlashcardCover, error) {
	return s.repository.ReadFlashcardCoverBySlug(flashcardCoverSlug)
}

func (s *service) UpdateFlashcardCover(flashcardCover *models.FlashcardCover) (*models.FlashcardCover, error, bool) {
	return s.repository.UpdateFlashcardCover(flashcardCover)
}

func (s *service) FetchFlashcardCoversByTagId(tagId *models.TagById) (*[]models.FlashcardCover, error) {
	return s.repository.ReadFlashcardCoversByTagId(tagId)
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
