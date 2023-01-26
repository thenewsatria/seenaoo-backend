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
		// ambil informasi user yang telah login
		currentUser := c.Locals("currentUser").(*models.User)

		// ambil profil berdasarkan user saat ini
		profileOwner := &models.UserProfileByOwner{Owner: currentUser.Username}
		userProfile, err := userProfileService.FetchProfileByOwner(profileOwner)

		if err != nil {
			// cek apakah error karena document tidak ditemukan
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// Membaca form
		formData, err := c.MultipartForm()
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FORM_PARSER_ERROR_MESSAGE))
		}

		// membaca file yang diupload melalui form dengan key avatar dan banner
		avatar := formData.File["avatar"]
		banner := formData.File["banner"]

		// deklarasi path lokal untuk menyimpan file
		newAvatarPath := ""
		newBannerPath := ""

		// Cek apakah file dengan key avatar tidak kosong
		if avatar != nil {
			// Cek extensi file dan ukurannya
			err := validator.ValidateFile(avatar[0], 2*1024*2014, []string{"image/png", "image/jpeg"})
			if err != nil {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// Membuka file untuk melalukan pengkopian ke direktori baru
			userAvatar := avatar[0]
			fileAvatar, err := userAvatar.Open()

			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// tutup file avatar diakhir
			defer fileAvatar.Close()

			// membuat folder penyimpanan file jika  tidak ada
			err = os.MkdirAll("./public/avatars", os.ModePerm)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating directory for avatars"))
			}

			// Membuat file baru untuk tujuan copy file avatar
			newAvatarPath = fmt.Sprintf("./public/avatars/%s_%d%s", currentUser.Username, time.Now().UnixNano(), filepath.Ext(userAvatar.Filename))
			newFile, err := os.Create(newAvatarPath)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error creating new file for avatars"))
			}

			// Menutup file baru diakhir
			defer newFile.Close()

			// hapus avatar lama jika bukan default avatar default
			if filepath.Base(userProfile.AvatarImagePath) != "default-avatar.png" {
				err := os.Remove(fmt.Sprintf("./public/avatars/%s", filepath.Base(userProfile.AvatarImagePath)))
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}

			// Mengcopy file avatar ke dalam file baru
			_, err = io.Copy(newFile, fileAvatar)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying avatar to directory"))
			}

			// Menyimpan path file baru yang dapat diakses melalui url static pada avatarImagePath
			userProfile.AvatarImagePath = fmt.Sprintf("http://%s/api/v1/static/avatars/%s", c.Hostname(), filepath.Base(newAvatarPath))
		}

		// Lakukan proses yang sama dengan banner
		if banner != nil {

			err = validator.ValidateFile(banner[0], 3*1024*1024, []string{"image/png", "image/jpeg"})
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

			// hapus banner lama jika bukan default
			if filepath.Base(userProfile.BannerImagePath) != "default-banner.jpg" {
				err := os.Remove(fmt.Sprintf("./public/banners/%s", filepath.Base(userProfile.BannerImagePath)))
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}

			_, err = io.Copy(newFile, fileBanner)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse("Error copying avatar to directory"))
			}

			userProfile.BannerImagePath = fmt.Sprintf("http://%s/api/v1/static/banners/%s", c.Hostname(), filepath.Base(newBannerPath))
		}

		// ubah display name dan biography pada userProfile sesuai input dari form
		userProfile.DisplayName = formData.Value["displayName"][0]
		userProfile.Biography = formData.Value["biography"][0]

		// lakukan update profile sesuai dengan service userprofiles
		updatedUserProfile, err, isValidationErr := userProfileService.UpdateProfile(userProfile)
		if err != nil {
			// Jika terdapat error maka hapus file yang sebelumnya telah diupload untuk mencegah duplikasi file
			// dilakukan pada banner dan juga avatar
			if avatar != nil {
				err = os.Remove(newAvatarPath)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}
			if banner != nil {
				err = os.Remove(newBannerPath)
				if err != nil {
					c.Status(http.StatusInternalServerError)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
			}
			// Cek apakah merupakan error validasi
			if isValidationErr {
				c.Status(http.StatusBadRequest)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		// jika tidak terdapat error maka lanjutkan mengirim response sukses
		c.Status(http.StatusOK)
		return c.JSON(presenters.UserProfileDetailSuccessResponse(updatedUserProfile, currentUser))
	}
}

func DeleteProfileBanner(userProfileService userprofiles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil user saat ini
		currentUser := c.Locals("currentUser").(*models.User)
		// 2. Ambil profile milik user
		upOwner := &models.UserProfileByOwner{Owner: currentUser.Username}
		currentUserProfile, err := userProfileService.FetchProfileByOwner(upOwner)

		// Cek apakah terdapat error pada pengambilan profile
		if err != nil {
			// cek apakah error karena document tidak ditemukan
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// 3. Seleksi kondisi nilai banner profile
		// - Jika merupakan default maka abaikan penghapusan
		// - Jika bukdan merupakan default maka hapus file tersebut dari direktori
		if filepath.Base(currentUserProfile.BannerImagePath) != "default-banner.jpg" {
			toDeleteFileName := fmt.Sprintf("./public/banners/%s", filepath.Base(currentUserProfile.BannerImagePath))
			err := os.Remove(toDeleteFileName)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// 4. Update userprofile dengan memasukan url static dari default banner
			defaultBannerPath := fmt.Sprintf("http://%s/api/v1/static/defaults/%s", c.Hostname(), "default-banner.jpg")
			currentUserProfile.BannerImagePath = defaultBannerPath

			// 5. Update profile dengan userprofiles service
			updatedUserProfile, err, isValidationError := userProfileService.UpdateProfile(currentUserProfile)
			if err != nil {
				if isValidationError {
					c.Status(http.StatusBadRequest)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_UPDATE_ERROR_MESSAGE))
			}

			// 6. Kirimkan response berupa user dan detil userprofile sebagai response
			c.Status(http.StatusOK)
			return c.JSON(presenters.UserProfileDetailSuccessResponse(updatedUserProfile, currentUser))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.UserProfileDetailSuccessResponse(currentUserProfile, currentUser))
	}
}

func DeleteProfileAvatar(userProfileService userprofiles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil user saat ini
		currentUser := c.Locals("currentUser").(*models.User)
		// 2. Ambil profile milik user
		upOwner := &models.UserProfileByOwner{Owner: currentUser.Username}
		currentUserProfile, err := userProfileService.FetchProfileByOwner(upOwner)

		// Cek apakah terdapat error pada pengambilan profile
		if err != nil {
			// cek apakah error karena document tidak ditemukan
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// 3. Seleksi kondisi nilai avatar profile
		// - Jika merupakan default maka abaikan penghapusan
		// - Jika bukan merupakan default maka hapus file tersebut dari direktori
		if filepath.Base(currentUserProfile.AvatarImagePath) != "default-avatar.png" {
			toDeleteFileName := fmt.Sprintf("./public/avatars/%s", filepath.Base(currentUserProfile.AvatarImagePath))
			err := os.Remove(toDeleteFileName)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(err.Error()))
			}

			// 4. Update userprofile dengan memasukan url static dari default avatar
			defaultAvatarPath := fmt.Sprintf("http://%s/api/v1/static/defaults/%s", c.Hostname(), "default-avatar.png")
			currentUserProfile.AvatarImagePath = defaultAvatarPath

			// 5. Update profile dengan userprofiles service
			updatedUserProfile, err, isValidationError := userProfileService.UpdateProfile(currentUserProfile)
			if err != nil {
				if isValidationError {
					c.Status(http.StatusBadRequest)
					return c.JSON(presenters.ErrorResponse(err.Error()))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_UPDATE_ERROR_MESSAGE))
			}

			// 6. Kirimkan response berupa user dan detil userprofile sebagai response
			c.Status(http.StatusOK)
			return c.JSON(presenters.UserProfileDetailSuccessResponse(updatedUserProfile, currentUser))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.UserProfileDetailSuccessResponse(currentUserProfile, currentUser))
	}
}

func GetUserProfile(userProfileService userprofiles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil user saat ini
		currentUser := c.Locals("currentUser").(*models.User)
		// 2. Ambil profile milik user
		upOwner := &models.UserProfileByOwner{Owner: currentUser.Username}
		currentUserProfile, err := userProfileService.FetchProfileByOwner(upOwner)

		// Cek apakah terdapat error pada pengambilan profile
		if err != nil {
			// cek apakah error karena document tidak ditemukan
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_PROFILE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		// 3. Kirimkan response berupa user dan detil userprofile sebagai response
		c.Status(http.StatusOK)
		return c.JSON(presenters.UserProfileDetailSuccessResponse(currentUserProfile, currentUser))
	}
}
