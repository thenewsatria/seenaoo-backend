package flashcards

import "github.com/thenewsatria/seenaoo-backend/pkg/models"

type FlashcardRepository interface {
	CreateFlashcard(flashcard *models.Flashcard) (*models.Flashcard, error)
}
