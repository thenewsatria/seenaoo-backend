package seeds

import (
	"fmt"
	"time"

	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SeedPermissionsCollection() {
	fmt.Println("Seeding permissions collection ...")

	permissionCollection := database.UseDB().Collection("permissions")

	fmt.Println("Clear existing permissions in permisisons collection ...")
	_, err := permissionCollection.DeleteMany(database.GetDBContext(), bson.M{})
	if err != nil {
		panic("Fail to clear existing permissions collection ...")
	}
	fmt.Println("Existing permissions collection has been cleared ...")

	fmt.Println("Inserting permissions into permissions collection ...")
	permissions := []interface{}{
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "UPDATE_FLASHCARD",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "DELETE_FLASHCARD",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "PURGE_FLASHCARD",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "ADD_CARD",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "UPDATE_CARD",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "DELETE_CARD",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "PURGE_CARD",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "ADD_CARD_HINT",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "UPDATE_CARD_HINT",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "DELETE_CARD_HINT",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "CLEAR_CARD_HINT",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			ItemCategory: "FLASHCARD",
			Name:         "INVITE_COLLABORATOR",
			Description:  "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
	insertCount, err := permissionCollection.InsertMany(database.GetDBContext(), permissions)
	if err != nil {
		panic("Fail to seeds permissions database")
	}
	fmt.Println("Successfully seeding permissions collection with total of " + fmt.Sprint(len(insertCount.InsertedIDs)) + " permissions.")
}
