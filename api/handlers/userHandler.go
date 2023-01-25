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
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/userprofiles"
	"github.com/thenewsatria/seenaoo-backend/utils/validator"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func EditUserProfile(userProfileService userprofiles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		profileOwner := &models.UserProfileByOwner{Owner: currentUser.Username}
		userProfile, err := userProfileService.FetchProfileByOwner(profileOwner)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		formData, err := c.MultipartForm()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FORM_PARSER_ERROR_MESSAGE))
		}

		// cek apakah terdapat file banner dan avatar
		// jika ada validasi terlebih dahulu untuk tipe data dan juga maximum file size allowed
		avatar := formData.File["avatar"]
		banner := formData.File["banner"]

		host := os.Getenv("HOST")
		port := os.Getenv("PORT")

		newAvatarPath := ""
		newBannerPath := ""

		if avatar != nil {
			// hapus avatar lama jika bukan default
			if filepath.Base(userProfile.AvatarImagePath) != "default-avatar.png" {
				err := os.Remove(fmt.Sprintf("./public/avatars/%s", filepath.Base(userProfile.AvatarImagePath)))
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}

			err := validator.ValidateFiles(avatar, 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			userAvatar := avatar[0]
			fileAvatar, err := userAvatar.Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			defer fileAvatar.Close()

			err = os.MkdirAll("./public/avatars", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for avatars"))
			}

			newAvatarPath = fmt.Sprintf("./public/avatars/%s_%d%s", currentUser.Username, time.Now().UnixNano(), filepath.Ext(userAvatar.Filename))
			newFile, err := os.Create(newAvatarPath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for avatars"))
			}

			defer newFile.Close()

			_, err = io.Copy(newFile, fileAvatar)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying avatar to directory"))
			}

			userProfile.AvatarImagePath = fmt.Sprintf("http://%s:%s/api/v1/static/avatars/%s", host, port, filepath.Base(newAvatarPath))
		}

		if banner != nil {
			// hapus banner lama jika bukan default
			if filepath.Base(userProfile.BannerImagePath) != "default-banner.jpg" {
				err := os.Remove(fmt.Sprintf("./public/banners/%s", filepath.Base(userProfile.BannerImagePath)))
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}

			err = validator.ValidateFiles(banner, 3*1024*1024, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			userBanner := banner[0]
			fileBanner, err := userBanner.Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			defer fileBanner.Close()

			err = os.MkdirAll("./public/banners", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				c.JSON(presenters.ErrorResponse("Error creating directory for banners"))
			}

			newBannerPath := fmt.Sprintf("./public/banners/%s_%d%s", currentUser.Username, time.Now().UnixNano(), filepath.Ext(userBanner.Filename))
			newFile, err := os.Create(newBannerPath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for banners"))
			}

			defer newFile.Close()

			_, err = io.Copy(newFile, fileBanner)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying avatar to directory"))
			}

			userProfile.BannerImagePath = fmt.Sprintf("http://%s/api/v1/static/banners/%s", c.Hostname(), filepath.Base(newBannerPath))
		}

		// jika ada save ke dalam folder public

		// ubah variabel path avatar atau banner sesuai dengan path url
		userProfile.DisplayName = formData.Value["displayName"][0]
		userProfile.Biography = formData.Value["biography"][0]

		// lakukan update profile sesuai dengan service
		updatedUserProfile, err, isValidationErr := userProfileService.UpdateProfile(userProfile)
		if err != nil {
			if banner != nil {
				err = os.Remove(newBannerPath)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}
			if avatar != nil {
				err = os.Remove(newAvatarPath)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}
			if isValidationErr {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		// setelah servis selesai maka tahapan selanjutnya adalah dengan mengecek apakah terdapat error
		// jika ada error maka hapus file yang baru saja tersimpan
		// jika tidak ada maka lanjutkan mengirim response sukses

		c.Status(http.StatusOK)
		return c.JSON(presenters.UserDetailSuccessResponse(currentUser, updatedUserProfile))
	}
}
