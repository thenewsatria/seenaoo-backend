package routes

import (
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
)

var flashcardCollection = database.UseDB().Collection(flashcards.CollectionName)
var flashcardRepo = flashcards.NewRepo(flashcardCollection)
var flashcardService = flashcards.NewService(flashcardRepo)

var flashcardHintCollection = database.UseDB().Collection(flashcardhints.CollectionName)
var flashcardHintRepo = flashcardhints.NewRepo(flashcardHintCollection)
var flashcardHintService = flashcardhints.NewService(flashcardHintRepo)
