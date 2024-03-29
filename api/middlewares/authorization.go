package middlewares

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/utils"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsLoggedIn(service users.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bearerToken, exist := c.GetReqHeaders()["Authorization"]
		if !exist {
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_NOT_EXIST_ERROR_MESSAGE))
		}
		token := strings.Split(bearerToken, " ")
		if len(token) < 2 || token[0] != "Bearer" {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_INVALID_ERROR_MESSAGE))
		}
		tokenStr := token[1]
		claims, err := utils.ParseAccessToken(tokenStr)
		if err != nil {
			if utils.IsTokenExpired(err) {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_ACCESS_TOKEN_EXPIRED_ERROR_MESSAGE))
			}
			c.Status(http.StatusUnauthorized)
			return c.JSON(presenters.ErrorResponse(messages.AUTH_TOKEN_INVALID_ERROR_MESSAGE))
		}

		userUsername := &models.UserByUsernameRequest{
			Username: claims.Username,
		}
		currentUser, err := service.FetchUserByUsername(userUsername)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.USER_USERNAME_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}
		c.Locals("currentUser", currentUser)
		return c.Next()
	}
}

func TestMW() fiber.Handler {
	return func(c *fiber.Ctx) error {
		testId := c.Params("testId")
		c.Locals("testingID", testId)
		return c.Next()
	}
}

func TestMW2() fiber.Handler {
	return func(c *fiber.Ctx) error {
		testId := c.Params("testId2")
		c.Locals("testingID", testId)
		c.Locals("passedParam", 123)
		return c.Next()
	}
}

func TestMW3() fiber.Handler {
	return func(c *fiber.Ctx) error {
		testId := c.Params("testId")
		newTestId := testId + "abc"
		c.Locals("passedParam", 123)
		c.Locals("mangstap", newTestId)
		return c.Next()
	}
}

func IsAllowedToSendCollaboration(service interface{}, collaborationService collaborations.Service,
	roleService roles.Service, isCollaboratorAllowed bool, isAuthorAllowed bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		itemCategory := strings.Split(c.Path(), "/")[4]
		switch itemCategory {
		case "flashcard":
			flashcardCoverService := service.(flashcardcovers.Service)
			fcCoverId := &models.FlashcardCoverById{ID: c.Params("itemId")}
			fcCover, err := flashcardCoverService.FetchFlashcardCoverById(fcCoverId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			if fcCover.Author != currentUser.Username {
				if isCollaboratorAllowed {
					cItemIdAndCollaborator := &models.CollaborationByItemIdAndCollaborator{ItemID: fcCover.ID.Hex(), Collaborator: currentUser.Username}
					isCollaborator, err := collaborationService.CheckIsCollaborator(cItemIdAndCollaborator)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.AUTH_MAKE_COLLABORATION_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					if !isCollaborator {
						c.Status(http.StatusUnauthorized)
						return c.JSON(presenters.ErrorResponse(messages.AUTH_MAKE_COLLABORATION_ERROR_MESSAGE))
					}
					collaborator, err := collaborationService.FetchCollaborationByItemIdAndCollaborator(cItemIdAndCollaborator)
					if err != nil {
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}

					roleId := &models.RoleById{ID: collaborator.ID.Hex()}
					role, err := roleService.FetchRoleById(roleId)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ROLE_ID_INVALID_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					c.Locals("isAuthor", false)
					c.Locals("userPermissions", role.Permissions)
					return c.Next()
				}
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_MAKE_COLLABORATION_ERROR_MESSAGE))
			} else if !isAuthorAllowed {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_COLLABORATION_UNAUTHORIZED_ERROR_MESSAGE))
			}
		default:
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.MIDDLEWARE_ISAUTHOR_UNKNOWN_SERVICE_TYPE_ERROR_MESSAGE))
		}

		c.Locals("isAuthor", true)
		return c.Next()
	}
}

func IsAuthorized(serviceName string, service interface{}, parentService interface{}, isCollaboratorAllowed bool, isAuthorAllowed bool,
	collaborationService collaborations.Service, roleService roles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		switch serviceName {
		case "COLLABORATION":
			collabService := service.(collaborations.Service)
			collaborationId := &models.CollaborationById{ID: c.Params("collaborationId")}
			collab, err := collabService.FetchCollaboration(collaborationId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			if collab.Inviter != currentUser.Username {
				if isCollaboratorAllowed {
					if collab.Collaborator != currentUser.Username {
						c.Status(http.StatusUnauthorized)
						return c.JSON(presenters.ErrorResponse(messages.AUTH_COLLABORATION_UNAUTHORIZED_ERROR_MESSAGE))
					}
					return c.Next()
				}
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_COLLABORATION_UNAUTHORIZED_ERROR_MESSAGE))
			} else if !isAuthorAllowed {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_COLLABORATION_UNAUTHORIZED_ERROR_MESSAGE))
			}
		case "FLASHCARD_COVER":
			flashcardCoverService := service.(flashcardcovers.Service)
			flashcardCoverSlug := &models.FlashcardCoverBySlug{Slug: c.Params("flashcardCoverSlug")}
			fcCover, err := flashcardCoverService.FetchFlashcardCoverBySlug(flashcardCoverSlug)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			if fcCover.Author != currentUser.Username {
				if isCollaboratorAllowed {
					cItemIdAndCollaborator := &models.CollaborationByItemIdAndCollaborator{ItemID: fcCover.ID.Hex(), Collaborator: currentUser.Username}
					isCollaborator, err := collaborationService.CheckIsCollaborator(cItemIdAndCollaborator)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_COVER_UNATHORIZED_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					if !isCollaborator {
						c.Status(http.StatusUnauthorized)
						return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_COVER_UNATHORIZED_ERROR_MESSAGE))
					}
					collaborator, err := collaborationService.FetchCollaborationByItemIdAndCollaborator(cItemIdAndCollaborator)
					if err != nil {
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}

					roleId := &models.RoleById{ID: collaborator.RoleId.Hex()}
					role, err := roleService.FetchRoleById(roleId)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ROLE_ID_INVALID_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					c.Locals("isAuthor", false)
					c.Locals("userPermissions", role.Permissions)
					return c.Next()
				}
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_COVER_UNATHORIZED_ERROR_MESSAGE))
			} else if !isAuthorAllowed {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_COLLABORATION_UNAUTHORIZED_ERROR_MESSAGE))
			}

		case "FLASHCARD_HINT":
			flashcardHintService := service.(flashcardhints.Service)
			fcHintId := &models.FlashcardHintByIdRequest{ID: c.Params("flashcardHintId")}
			fcHint, err := flashcardHintService.FetchFlashcardHint(fcHintId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_HINT_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			flashcardCoverService := parentService.(flashcardcovers.Service)
			fcCoverId := &models.FlashcardCoverById{ID: fcHint.FlashcardCoverId.Hex()}
			fcCover, err := flashcardCoverService.FetchFlashcardCoverById(fcCoverId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			if fcCover.Author != currentUser.Username {
				if isCollaboratorAllowed {
					cItemIdAndCollaborator := &models.CollaborationByItemIdAndCollaborator{ItemID: fcCover.ID.Hex(), Collaborator: currentUser.Username}
					isCollaborator, err := collaborationService.CheckIsCollaborator(cItemIdAndCollaborator)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_HINT_UNATHORIZED_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					if !isCollaborator {
						c.Status(http.StatusUnauthorized)
						return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_HINT_UNATHORIZED_ERROR_MESSAGE))
					}
					collaborator, err := collaborationService.FetchCollaborationByItemIdAndCollaborator(cItemIdAndCollaborator)
					if err != nil {
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}

					roleId := &models.RoleById{ID: collaborator.RoleId.Hex()}
					role, err := roleService.FetchRoleById(roleId)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ROLE_ID_INVALID_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					c.Locals("isAuthor", false)
					c.Locals("userPermissions", role.Permissions)
					return c.Next()
				}
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_HINT_UNATHORIZED_ERROR_MESSAGE))
			} else if !isAuthorAllowed {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_COLLABORATION_UNAUTHORIZED_ERROR_MESSAGE))
			}
		case "FLASHCARD":
			flashcardService := service.(flashcards.Service)
			flashcardId := &models.FlashcardByIdRequest{ID: c.Params("flashcardId")}
			fc, err := flashcardService.FetchFlashcard(flashcardId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			flashcardCoverService := parentService.(flashcardcovers.Service)
			fcCoverId := &models.FlashcardCoverById{ID: fc.FlashCardCoverId.Hex()}
			fcCover, err := flashcardCoverService.FetchFlashcardCoverById(fcCoverId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.FLASHCARD_COVER_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			if fcCover.Author != currentUser.Username {
				if isCollaboratorAllowed {
					cItemIdAndCollaborator := &models.CollaborationByItemIdAndCollaborator{ItemID: fcCover.ID.Hex(), Collaborator: currentUser.Username}
					isCollaborator, err := collaborationService.CheckIsCollaborator(cItemIdAndCollaborator)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_UNAUTHORIZED_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					if !isCollaborator {
						c.Status(http.StatusUnauthorized)
						return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_UNAUTHORIZED_ERROR_MESSAGE))
					}
					collaborator, err := collaborationService.FetchCollaborationByItemIdAndCollaborator(cItemIdAndCollaborator)
					if err != nil {
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_FAIL_TO_FETCH_ERROR_MESSAGE))
					}

					roleId := &models.RoleById{ID: collaborator.RoleId.Hex()}
					role, err := roleService.FetchRoleById(roleId)
					if err != nil {
						if err == mongo.ErrNoDocuments {
							c.Status(http.StatusNotFound)
							return c.JSON(presenters.ErrorResponse(messages.COLLABORATION_ROLE_ID_INVALID_ERROR_MESSAGE))
						}
						c.Status(http.StatusInternalServerError)
						return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
					}
					c.Locals("isAuthor", false)
					c.Locals("userPermissions", role.Permissions)
					return c.Next()
				}
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_FLASHCARD_UNAUTHORIZED_ERROR_MESSAGE))
			}
		case "ROLE":
			roleService := service.(roles.Service)
			roleSlug := &models.RoleBySlug{Slug: c.Params("roleSlug")}
			rl, err := roleService.FetchRoleBySlug(roleSlug)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.ROLE_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
			}

			if rl.Owner != currentUser.Username {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_ROLE_UNAUTHORIZED_ERROR_MESSAGE))
			}
		default:
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.MIDDLEWARE_ISAUTHOR_UNKNOWN_SERVICE_TYPE_ERROR_MESSAGE))
		}

		c.Locals("isAuthor", true)
		return c.Next()
	}
}

func HavePermit(permissionService permissions.Service, mustHaveAll bool, allowedPermissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isAuthor := c.Locals("isAuthor").(bool); isAuthor {
			return c.Next()
		}
		userPermissions := c.Locals("userPermissions").([]primitive.ObjectID)
		userPermitStr := []string{}
		for _, permission := range userPermissions {
			permitId := &models.PermissionById{ID: permission.Hex()}
			permit, err := permissionService.FetchPermissionById(permitId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.ROLE_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
			permitString := permit.ItemCategory + "." + permit.Name
			userPermitStr = append(userPermitStr, permitString)
		}

		var trueAllowedPermit = make([]int, len(allowedPermissions))
		var allowedPermitCount = 0

		for i := 0; i < len(allowedPermissions); i++ {
			for _, userPermit := range userPermitStr {
				if userPermit == allowedPermissions[i] {
					if !mustHaveAll {
						return c.Next()
					}
					trueAllowedPermit[i] = 1
					break
				}
			}
			if mustHaveAll && trueAllowedPermit[i] != 1 {
				c.Status(http.StatusUnauthorized)
				return c.JSON(presenters.ErrorResponse(messages.AUTH_DONT_HAVE_SUITABLE_PERMISSION_ERROR_MESSAGE))
			}
		}
		for _, num := range trueAllowedPermit {
			allowedPermitCount += num
		}
		if allowedPermitCount == len(allowedPermissions) {
			return c.Next()
		}

		c.Status(http.StatusUnauthorized)
		return c.JSON(presenters.ErrorResponse(messages.AUTH_DONT_HAVE_SUITABLE_PERMISSION_ERROR_MESSAGE))
	}
}
