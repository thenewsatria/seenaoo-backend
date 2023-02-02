package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/tags"
	"github.com/thenewsatria/seenaoo-backend/pkg/userprofiles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddFlashcardCover(flashcardCoverService flashcardcovers.Service, tagService tags.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// dapatkan user yang dioper melalui middleware
		currentUser := c.Locals("currentUser").(*models.User)

		// baca form data
		formData, err := c.MultipartForm()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE))
		}

		currentTimeStr := fmt.Sprintf("%v", time.Now().Unix())

		// buat entitas flashcardCover
		fcCover := &models.FlashcardCover{
			Slug:        slug.Make(formData.Value["title"][0]) + "-" + currentTimeStr,
			Title:       formData.Value["title"][0],
			Description: formData.Value["description"][0],
			Author:      currentUser.Username,
		}

		// baca form data dengan tipe file dengan key coverImg
		fcImg := formData.File["coverImage"]

		// inisialisasi lokasi penyimpanan gambar cover
		newCvrImgPath := ""

		// jika terdapat file gambar cover maka
		if fcImg != nil {
			// Cek extensi file dan ukurannya
			err := validator.ValidateFile(fcImg[0], 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// Membuka file untuk melalukan pengkopian ke direktori baru
			fileCvrImg, err := fcImg[0].Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// tutup file cover diakhir
			defer fileCvrImg.Close()

			// membuat folder penyimpanan file jika  tidak ada sesuai lokasi
			err = os.MkdirAll("./public/flashcardcovers", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for flashcard cover image"))
			}

			// Membuat file baru untuk tujuan copy file flashcard cover
			newCvrImgPath = fmt.Sprintf("./public/flashcardcovers/%s_%d%s", currentUser.Username, time.Now().UnixNano(), filepath.Ext(fcImg[0].Filename))
			newFile, err := os.Create(newCvrImgPath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for flashcard cover images"))
			}

			// Menutup file baru diakhir
			defer newFile.Close()

			// Mengcopy file cover image ke dalam file baru
			_, err = io.Copy(newFile, fileCvrImg)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying flascard cover image to directory"))
			}

			// Menyimpan path file baru yang dapat diakses melalui url static pada imagePath
			fcCover.ImagePath = fmt.Sprintf("http://%s/api/v1/static/flashcardcovers/%s", c.Hostname(), filepath.Base(newCvrImgPath))
		}

		tagIds := []primitive.ObjectID{}

		// Proses pembacaan tag, untuk setiap nilai pada key tags pada form
		// maka buat baru jika tag tidak ada, jika ada masukan id pada tagIds

		for _, tagString := range formData.Value["tags"] {
			// Check if tagString is empty string, if yes ignore
			if tagString != "" {
				tagName := &models.TagByName{TagName: tagString}
				existedTag, err := tagService.FetchTagByName(tagName)
				if err != nil {
					if err == mongo.ErrNoDocuments { //jika tag tidak ada maka buat baru
						tag := &models.Tag{TagName: tagString}
						newTag, err, isValidationError := tagService.InsertTag(tag)
						if err != nil {
							if isValidationError {
								c.Status(http.StatusBadRequest)
								return c.JSON(presenters.ErrorResponse(err.Error()))
							}
							c.Status(http.StatusInternalServerError)
							return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_INSERT_ERROR_MESSAGE))
						}
						tagIds = append(tagIds, newTag.ID)
						continue
					}
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_FETCH_ERROR_MESSAGE))
				}
				tagIds = append(tagIds, existedTag.ID)
			}
		}

		fcCover.Tags = tagIds

		// Proses update
		insertedFcCover, err, isValidationError := flashcardCoverService.InsertFlashcardCover(fcCover)
		if err != nil {
			// hapus file yang baru saja dimasukan jika gagal
			if fcImg != nil {
				errFile := os.Remove(newCvrImgPath)
				if errFile != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(errFile.Error()))
				}
			}
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		return c.JSON(presenters.FlashcardCoverSuccessResponse(insertedFcCover))
	}
}

func GetFlashcardCover(flashcardCoverService flashcardcovers.Service, tagService tags.Service,
	userService users.Service, userProfileService userprofiles.Service, flashcardService flashcards.Service) fiber.Handler {

	return func(c *fiber.Ctx) error {
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		tagDetails := []models.Tag{}

		for _, fcTag := range fcCover.Tags {
			tagId := &models.TagById{ID: fcTag.Hex()}
			tag, err := tagService.FetchTagById(tagId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.TAG_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
			tagDetails = append(tagDetails, *tag)
		}

		userUname := &models.UserByUsernameRequest{Username: fcCover.Author}
		author, err := userService.FetchUserByUsername(userUname)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		fcCvrId := &models.FlashcardCoverById{ID: fcCover.ID.Hex()}
		flashcards, err := flashcardService.PopulateFlashcardCover(fcCvrId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_POPULATE_FLASHCARDS_ERROR_MESSAGE))
		}

		upOwner := &models.UserProfileByOwner{Owner: author.Username}
		userProfile, err := userProfileService.FetchProfileByOwner(upOwner)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverDetailSuccessResponse(fcCover, &tagDetails, flashcards, author, userProfile))
	}
}

func UpdateFlashcardCover(flashcardCoverService flashcardcovers.Service, tagService tags.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Cari flashcard berdasarkan slug
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// dapatkan user saat ini melalui middleware
		currentUser := c.Locals("currentUser").(*models.User)

		// Baca form data
		formData, err := c.MultipartForm()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE))
		}

		// Update flashcard cover
		newSlug := slug.Make(formData.Value["title"][0]) + "-" + fmt.Sprintf("%v", time.Now().Unix())
		fcCover.Slug = newSlug
		fcCover.Title = formData.Value["title"][0]
		fcCover.Description = formData.Value["description"][0]

		// membaca file yang diupload melalui form dengan key coverImage
		fcImg := formData.File["coverImage"]

		// deklarasi path lokal untuk menyimpan file
		newFcCvrImgPath := ""

		// variabel untuk tujuan penghapusan jika proses update pada database sukses
		currentCvrImgPath := fcCover.ImagePath

		// jika file yang diupload tidak ksong
		if fcImg != nil {
			// Cek extensi file dan ukurannya
			err := validator.ValidateFile(fcImg[0], 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// Membuka file untuk melalukan pengkopian ke direktori baru
			fileCvrImg, err := fcImg[0].Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// tutup file untuk pengkopian diakhir
			defer fileCvrImg.Close()

			// membuat folder penyimpanan file jika  tidak ada
			err = os.MkdirAll("./public/flashcardcovers", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for flashcard cover image"))
			}

			// Membuat file baru untuk tujuan copy file cover image
			newFcCvrImgPath = fmt.Sprintf("./public/flashcardcovers/%s_%d%s", currentUser.Username, time.Now().UnixNano(), filepath.Ext(fcCover.ImagePath))
			newFile, err := os.Create(newFcCvrImgPath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for flashcard cover image"))
			}

			// Menutup file tujuan pengkopian diakhir
			defer newFile.Close()

			// Mengcopy file coverImg yang diupload melalui form ke dalam file baru
			_, err = io.Copy(newFile, fileCvrImg)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("rror copying flascard cover image to directory"))
			}

			// Menyimpan path file baru yang dapat diakses melalui url static pada ImagePath
			fcCover.ImagePath = fmt.Sprintf("http://%s/api/v1/static/flashcardcovers/%s", c.Hostname(), filepath.Base(newFcCvrImgPath))
		}

		tagIds := []primitive.ObjectID{}

		// Proses pembacaan tag, untuk setiap nilai pada key tags pada form
		// maka buat baru jika tag tidak ada, jika ada masukan id pada tagIds

		for _, tagString := range formData.Value["tags"] {
			if tagString != "" {
				tagName := &models.TagByName{TagName: tagString}
				existedTag, err := tagService.FetchTagByName(tagName)
				if err != nil {
					if err == mongo.ErrNoDocuments { //jika tag tidak ada maka buat baru
						tag := &models.Tag{TagName: tagString}
						newTag, err, isValidationError := tagService.InsertTag(tag)
						if err != nil {
							if isValidationError {
								c.Status(http.StatusBadRequest)
								return c.JSON(presenters.ErrorResponse(err.Error()))
							}
							c.Status(http.StatusInternalServerError)
							return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_INSERT_ERROR_MESSAGE))
						}
						tagIds = append(tagIds, newTag.ID)
						continue
					}
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(messages.TAG_FAIL_TO_FETCH_ERROR_MESSAGE))
				}
				tagIds = append(tagIds, existedTag.ID)
			}
		}

		fcCover.Tags = tagIds

		// Melakukan update pada database
		updatedFcCover, err, isValidationError := flashcardCoverService.UpdateFlashcardCover(fcCover)
		if err != nil {

			// hapus file yang baru saja diupload jika error terjadi
			if fcImg != nil {
				errFile := os.Remove(newFcCvrImgPath)
				if errFile != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(errFile.Error()))
				}
			}
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		// Hapus gambar flashcard cover yang lama jika proses update dengan upload gambar baru berhasil disimpan
		if currentCvrImgPath != "" && fcImg != nil {
			err := os.Remove(fmt.Sprintf("./public/flashcardcovers/%s", filepath.Base(currentCvrImgPath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverSuccessResponse(updatedFcCover))
	}
}

func DeleteFlashcardCover(flashcardCoverService flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Memperoleh flashcardCover melalui slug
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
		}

		// Menghapus flashcardCover
		deletedFcCover, err := flashcardCoverService.RemoveFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		// Menghapus gambar flashcardCover yang telah sukses dihapus
		if deletedFcCover.ImagePath != "" {
			err := os.Remove(fmt.Sprintf("./public/flashcardcovers/%s", filepath.Base(fcCover.ImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverSuccessResponse(deletedFcCover))
	}
}

func PurgeFlashcardCover(flashcardCoverService flashcardcovers.Service, flashcardService flashcards.Service, flashcardHintService flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//dapatkan flashcard cover melalui slug
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
		}

		//Get flashcards by flashcard cover id so we able to delete each flashcard's hints
		fcCoverId := &models.FlashcardCoverById{ID: fcCover.ID.Hex()}
		flashcards, err := flashcardService.PopulateFlashcardCover(fcCoverId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_POPULATE_FLASHCARDS_ERROR_MESSAGE))
		}

		//Loop each flashcard to delete each flashcard's hints
		for _, flashcard := range *flashcards {
			flashcardId := &models.FlashcardByIdRequest{ID: flashcard.ID.Hex()}
			_, err := flashcardHintService.RemoveFlashcardHintsByFlashcardId(flashcardId)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_DELETE_ERROR_MESSAGE))
			}
		}

		//Delete all flashcard with the same flashcard cover id
		_, err = flashcardService.RemoveFlashcardsByFlashcardCoverId(fcCoverId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		//Loop each flashcard to delete each flashcard's frontImage and backImage
		for _, flashcard := range *flashcards {
			if flashcard.FrontImagePath != "" {
				toDeleteFileName := fmt.Sprintf("./public/flashcards/fronts/%s", filepath.Base(flashcard.FrontImagePath))
				if err := os.Remove(toDeleteFileName); err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}
			if flashcard.BackImagePath != "" {
				toDeleteFileName := fmt.Sprintf("./public/flashcards/backs/%s", filepath.Base(flashcard.BackImagePath))
				if err := os.Remove(toDeleteFileName); err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}
		}

		// Delete The flashcard cover
		deletedFcCover, err := flashcardCoverService.RemoveFlashcardCover(fcCover)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		// Dilakukan penghapusan pada gambar cover setelah flashcard cover berhasil dihapus
		if deletedFcCover.ImagePath != "" {
			err := os.Remove(fmt.Sprintf("./public/flashcardcovers/%s", filepath.Base(fcCover.ImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverSuccessResponse(deletedFcCover))
	}
}

func DeleteFlashcardCoverImage(flashcardCoverService flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil flashcard cover sesuai dengan slug
		fcCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(fcCoverSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// Menyimpan path gambar sebelumnya untuk melakukan penghapusan jika penghapusan pada database sukses
		currentImgPath := fcCover.ImagePath
		fcCover.ImagePath = ""

		// Update flashcard cover dengan flashcardCover service
		updatedFlashcardCover, err, isValidationError := flashcardCoverService.UpdateFlashcardCover(fcCover)
		if err != nil {
			if isValidationError {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		// Jika cover yang ada bukan merupakan default maka hapus gambar yang ada pada lokal
		if currentImgPath != "" {
			toDeleteFileName := fmt.Sprintf("./public/flashcardcovers/%s", filepath.Base(currentImgPath))
			err = os.Remove(toDeleteFileName)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}
		//Kirimkan response sukses dengan presenter terkait
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardCoverSuccessResponse(updatedFlashcardCover))
	}
}
