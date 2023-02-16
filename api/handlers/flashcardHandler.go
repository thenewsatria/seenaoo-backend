package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddFlashcard(flashcardService flashcards.Service, flashcardCoverService flashcardcovers.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mendapatkan informasi flashcard cover berdasarkan slug yang akan ditambahkan flashcard
		fCvrSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
		fcCvr, err := flashcardCoverService.FetchFlashcardCoverBySlug(fCvrSlug)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
		}

		// mendapatkan informasi mengenai user yang saat ini login
		currentUser := c.Locals("currentUser").(*models.User)

		// Membaca / memparsing informasi yang dikirimkan dengan form data
		formData, err := c.MultipartForm()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE))
		}

		// Membentuk flashcard baru dengan informasi dari form data
		flashcard := &models.Flashcard{
			FrontText:        formData.Value["frontText"][0],
			BackText:         formData.Value["backText"][0],
			Question:         formData.Value["question"][0],
			FlashCardCoverId: fcCvr.ID,
		}

		// menyimpan file yang di upload melalui form dengan key frontImage dan backImage
		frontImage := formData.File["frontImage"]
		backImage := formData.File["backImage"]

		// deklarasi path lokal untuk menyimpan file
		newFrontImagePath := ""
		newBackImagePath := ""

		// Mengecek validasi file yang diupload pada frontImage dan backImage jika terdapat file pada keduanya / tidak nil
		if frontImage != nil {
			// Cek extensi file dan ukurannya
			err := validator.ValidateFile(frontImage[0], 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		if backImage != nil {
			// Cek extensi file dan ukurannya
			err := validator.ValidateFile(backImage[0], 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		// Jika gambar frontImage tidak nil maka:
		if frontImage != nil {
			// Membuka file untuk melalukan pengkopian ke direktori baru
			userFrontImage := frontImage[0]
			fileFrontImage, err := userFrontImage.Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// tutup file frontImage diakhir
			defer fileFrontImage.Close()

			// membuat folder penyimpanan file jika  tidak ada
			err = os.MkdirAll("./public/flashcards/fronts", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for flashcard front image"))
			}

			// Membuat file baru untuk tujuan copy file frontImage
			newFrontImagePath = fmt.Sprintf("./public/flashcards/fronts/%s_%d%s", currentUser.Username, time.Now().UnixNano(),
				filepath.Ext(userFrontImage.Filename))
			newFile, err := os.Create(newFrontImagePath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for flashcard front image"))
			}

			// Menutup file baru diakhir
			defer newFile.Close()

			// Mengcopy file frontImage ke dalam file baru
			_, err = io.Copy(newFile, fileFrontImage)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying flashcard front image to directory"))
			}

			// Menyimpan path file baru yang dapat diakses melalui url static pada frontImagePath flashcard
			flashcard.FrontImagePath = fmt.Sprintf("http://%s/api/v1/static/flashcards/fronts/%s", c.Hostname(), filepath.Base(newFrontImagePath))
		}

		// Jika gambar backImage tidak nil maka:
		if backImage != nil {
			// Membuka file untuk melalukan pengkopian ke direktori baru
			userBackImage := backImage[0]
			fileBackImage, err := userBackImage.Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// tutup file backImage diakhir
			defer fileBackImage.Close()

			// membuat folder penyimpanan file jika  tidak ada
			err = os.MkdirAll("./public/flashcards/backs", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for flashcard back image"))
			}

			// Membuat file baru untuk tujuan copy file backImage
			newBackImagePath = fmt.Sprintf("./public/flashcards/backs/%s_%d%s", currentUser.Username, time.Now().UnixNano(),
				filepath.Ext(userBackImage.Filename))
			newFile, err := os.Create(newBackImagePath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for flashcard back image"))
			}

			// Menutup file baru diakhir
			defer newFile.Close()

			// Mengcopy file avatar ke dalam file baru
			_, err = io.Copy(newFile, fileBackImage)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying flashcard back image to directory"))
			}

			// Menyimpan path file baru yang dapat diakses melalui url static pada backImagePath
			flashcard.BackImagePath = fmt.Sprintf("http://%s/api/v1/static/flashcards/backs/%s", c.Hostname(), filepath.Base(newBackImagePath))
		}

		// Membuat flashcard baru dengan service flashcardService
		result, err, isValidationError := flashcardService.InsertFlashcard(flashcard)
		// jika error maka
		if err != nil {
			// Hapus gambar yang baru saja diupload jika terdapat gambar yang diupload (frontImage dan backImage)
			if frontImage != nil {
				errFile := os.Remove(newFrontImagePath)
				if errFile != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(errFile.Error()))
				}
			}
			if backImage != nil {
				errFile := os.Remove(newBackImagePath)
				if errFile != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(errFile.Error()))
				}
			}
			if isValidationError {
				translatedErrors := validator.TranslateError(err)
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.FailResponse(translatedErrors))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_INSERT_ERROR_MESSAGE))
		}

		// Kirimkan sukses response dengan presenter terkait
		c.Status(http.StatusCreated)
		return c.JSON(presenters.FlashcardSuccessResponse(result))
	}
}

func GetFlashcard(flashcardService flashcards.Service, flashcardHintService flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		result, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		hints, err := flashcardHintService.PopulateFlashcard(flashcardId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_POPULATE_HINTS_ERROR_MESSAGE))
		}
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardDetailSuccessResponse(result, hints))
	}
}

func UpdateFlashcard(flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// dapatkan flashcard sesuai dengan flashcard Id pada ulr parameter
		flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		flashcard, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// dapatkan informasi user yang login saat ini
		currentUser := c.Locals("currentUser").(*models.User)

		// parsing body form data
		formData, err := c.MultipartForm()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_BODY_PARSER_ERROR_MESSAGE))
		}

		// Simpan file yang diupload melalui form dengan key frontImage path dan backImage
		frontImage := formData.File["frontImage"]
		backImage := formData.File["backImage"]

		// deklarasi path lokal untuk menyimpan file
		newFrontImagePath := ""
		newBackImagePath := ""

		// Simpan lokasi file saat ini untuk tujuan penghapusan jika proses update selesai
		currentFrontImagePath := flashcard.FrontImagePath
		currentBackImagePath := flashcard.BackImagePath

		// jika terdapat file yang diupload maka cek validitas filenya
		// pengecekan validasi dipisahkan dengan proses pengkopian ke lokal karena jika validasi file kedua (backImage) gagal
		// namun validasi file pertama sukses (frontImage) maka, file frontImage akan tetap masuk ke lokal walaupun flashcard
		// gagal tercipta, sehingga perlu dipisahkan.
		if frontImage != nil {
			// Cek extensi file dan ukurannya
			err := validator.ValidateFile(frontImage[0], 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		if backImage != nil {
			// Cek extensi file dan ukurannya
			err := validator.ValidateFile(backImage[0], 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		// jika file yang diupload ada maka:
		if frontImage != nil {
			// Membuka file untuk melalukan pengkopian ke direktori baru
			userFrontImage := frontImage[0]
			fileFrontImage, err := userFrontImage.Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// tutup file frontImage diakhir
			defer fileFrontImage.Close()

			// membuat folder penyimpanan file jika  tidak ada
			err = os.MkdirAll("./public/flashcards/fronts", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for flashcard front image"))
			}

			// Membuat file baru untuk tujuan copy file frontImage
			newFrontImagePath = fmt.Sprintf("./public/flashcards/fronts/%s_%d%s", currentUser.Username, time.Now().UnixNano(),
				filepath.Ext(userFrontImage.Filename))
			newFile, err := os.Create(newFrontImagePath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for flashcard front image"))
			}

			// Menutup file baru diakhir
			defer newFile.Close()

			// Mengcopy file frontImage ke dalam file baru
			_, err = io.Copy(newFile, fileFrontImage)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying flashcard front image to directory"))
			}

			// Menyimpan path file baru yang dapat diakses melalui url static pada frontImagePath
			flashcard.FrontImagePath = fmt.Sprintf("http://%s/api/v1/static/flashcards/fronts/%s", c.Hostname(), filepath.Base(newFrontImagePath))
		}

		if backImage != nil {
			// Membuka file untuk melalukan pengkopian ke direktori baru
			userBackImage := backImage[0]
			fileBackImage, err := userBackImage.Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// tutup file backImage diakhir
			defer fileBackImage.Close()

			// membuat folder penyimpanan file jika  tidak ada
			err = os.MkdirAll("./public/flashcards/backs", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for flashcard back image"))
			}

			// Membuat file baru untuk tujuan copy file backImage
			newBackImagePath = fmt.Sprintf("./public/flashcards/backs/%s_%d%s", currentUser.Username, time.Now().UnixNano(),
				filepath.Ext(userBackImage.Filename))
			newFile, err := os.Create(newBackImagePath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for flashcard back image"))
			}

			// Menutup file baru diakhir
			defer newFile.Close()

			// Mengcopy file backImage ke dalam file baru
			_, err = io.Copy(newFile, fileBackImage)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying flashcard back image to directory"))
			}

			// Menyimpan path file baru yang dapat diakses melalui url static pada backImagePath
			flashcard.BackImagePath = fmt.Sprintf("http://%s/api/v1/static/flashcards/backs/%s", c.Hostname(), filepath.Base(newBackImagePath))
		}

		// Merubah value dari atribut flashcard sesuai masukan dari form data
		flashcard.FrontText = formData.Value["frontText"][0]
		flashcard.BackText = formData.Value["backText"][0]
		flashcard.Question = formData.Value["question"][0]

		// Melakukan proses update flashcard
		updatedFlashcard, err, isValidationError := flashcardService.UpdateFlashcard(flashcard)
		// jika terdapat error
		if err != nil {
			// Hapus file yang terupload di lokal jika terdapat file yang diupload
			if frontImage != nil {
				errFile := os.Remove(newFrontImagePath)
				if errFile != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(errFile.Error()))
				}
			}
			if backImage != nil {
				errFile := os.Remove(newBackImagePath)
				if errFile != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(errFile.Error()))
				}
			}
			if isValidationError {
				translatedErrors := validator.TranslateError(err)
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.FailResponse(translatedErrors))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		// jika proses update sukses baru dilakukan penghapusan terhadap file sebelumnya (jika bukan merupakan file default dan
		// terdapat file yang diupload baik pada frontImage maupaun backImage)
		if currentFrontImagePath != "" && frontImage != nil {
			err := os.Remove(fmt.Sprintf("./public/flashcards/fronts/%s", filepath.Base(currentFrontImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		if currentBackImagePath != "" && backImage != nil {
			err := os.Remove(fmt.Sprintf("./public/flashcards/backs/%s", filepath.Base(currentBackImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		// Mengirimkan sukses response
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardSuccessResponse(updatedFlashcard))
	}
}

func DeleteFlashcard(flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Memperoleh flashcard melalui flashcardId dari parameter url
		flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		flashcard, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// Menghapus flashcard yang didapatkan dengan flashcardService
		deletedFlashcard, err := flashcardService.RemoveFlashcard(flashcard)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		// dari flashcard yang sukses dihapus dilakukan penghapusan pada gambar frontImage yang tersimpan secara lokal
		// sesuai dengan path file dari atribut flashcard yang terhapus
		if deletedFlashcard.FrontImagePath != "" {
			err := os.Remove(fmt.Sprintf("./public/flashcards/fronts/%s", filepath.Base(deletedFlashcard.FrontImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}
		if deletedFlashcard.BackImagePath != "" {
			err := os.Remove(fmt.Sprintf("./public/flashcards/backs/%s", filepath.Base(deletedFlashcard.BackImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		// Kirimkan response sukses
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardSuccessResponse(deletedFlashcard))
	}
}

func PurgeFlashcard(flashcardService flashcards.Service, flashcardHintService flashcardhints.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Dapatkan flashcard dengan flashcard id sesuai dengan parameter flashcardId pada URL
		flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		flashcard, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// lakukan penghapusan seluruh flashcard hints yang terkait dengan flashcard id
		_, err = flashcardHintService.RemoveFlashcardHintsByFlashcardId(flashcardId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		// proses penghapusan flashcard dengan flashcardService
		deletedFc, err := flashcardService.RemoveFlashcard(flashcard)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		// dari flashcard yang sukses dihapus dilakukan penghapusan pada gambar frontImage yang tersimpan secara lokal
		// sesuai dengan path file dari atribut flashcard yang terhapus
		if deletedFc.FrontImagePath != "" {
			err := os.Remove(fmt.Sprintf("./public/flashcards/fronts/%s", filepath.Base(deletedFc.FrontImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}
		if deletedFc.BackImagePath != "" {
			err := os.Remove(fmt.Sprintf("./public/flashcards/backs/%s", filepath.Base(deletedFc.BackImagePath)))
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		// Kriimkan response sukses
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardSuccessResponse(deletedFc))
	}
}

func DeleteFlashcardFrontImage(flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Dapatkan flashcard dengan flashcard id sesuai dengan parameter flashcardId pada URL
		flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		flashcard, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// simpan front image path saat ini
		currentFlashcardFrontImage := flashcard.FrontImagePath

		// set atribut front image path dengan string kosong
		flashcard.FrontImagePath = ""

		// lakukan update pada flashcard
		updatedFlashcard, err, isValidationError := flashcardService.UpdateFlashcard(flashcard)
		if err != nil {
			if isValidationError {
				translatedErrors := validator.TranslateError(err)
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.FailResponse(translatedErrors))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		// setelah proses update selesai dan atribut flashcardFrontImage sebelumnya terdapat path file static, maka
		// dilakukan penghapusan file tersebut
		if currentFlashcardFrontImage != "" {
			toDeleteFileName := fmt.Sprintf("./public/flashcards/fronts/%s", filepath.Base(currentFlashcardFrontImage))
			err := os.Remove(toDeleteFileName)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		// Kirimkan proses response sukses dengan flashcard yang telah diupdate
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardSuccessResponse(updatedFlashcard))
	}
}

func DeleteFlashcardBackImage(flashcardService flashcards.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Dapatkan flashcard dengan flashcard id sesuai dengan parameter flashcardId pada URL
		flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
		flashcard, err := flashcardService.FetchFlashcard(flashcardId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
		}
		// simpan front image path saat ini
		currentFlashcardBackImage := flashcard.BackImagePath

		// set atribut front image path dengan string kosong
		flashcard.BackImagePath = ""

		// lakukan update pada flashcard
		updatedFlashcard, err, isValidationError := flashcardService.UpdateFlashcard(flashcard)
		if err != nil {
			if isValidationError {
				translatedErrors := validator.TranslateError(err)
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.FailResponse(translatedErrors))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		// setelah proses update selesai dan atribut flashcardFrontImage sebelumnya terdapat path file static, maka
		// dilakukan penghapusan file tersebut
		if currentFlashcardBackImage != "" {
			toDeleteFileName := fmt.Sprintf("./public/flashcards/backs/%s", filepath.Base(currentFlashcardBackImage))
			err := os.Remove(toDeleteFileName)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
		}

		// Kirimkan proses response sukses dengan flashcard yang telah diupdate
		c.Status(http.StatusOK)
		return c.JSON(presenters.FlashcardSuccessResponse(updatedFlashcard))
	}
}
