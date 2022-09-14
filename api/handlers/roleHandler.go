package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeNewRole(roleService roles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		var roleReq = &models.RoleRequest{}
		if err := c.BodyParser(roleReq); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_BODY_PARSER_ERROR_MESSAGE))
		}

		currentTimeStr := fmt.Sprintf("%v", time.Now().Unix())
		slug := roleReq.Name + "-" + currentTimeStr

		var role = &models.Role{
			Owner:       currentUser.Username,
			Name:        roleReq.Name,
			Slug:        slug,
			Description: roleReq.Description,
			Permissions: roleReq.Permissions,
		}

		result, err := roleService.InsertRole(role)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_INSERT_ERROR_MESSAGE))
		}
		c.Status(http.StatusOK)
		return c.JSON(presenters.RoleSuccessResponse(result))
	}
}

func GetRole(roleService roles.Service, userService users.Service, permissionService permissions.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleSlug := &models.RoleBySlug{Slug: c.Params("roleSlug")}
		role, err := roleService.FetchRoleBySlug(roleSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.ROLE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		userUname := &models.UserByUsernameRequest{Username: role.Owner}
		user, err := userService.FetchUserByUsername(userUname)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.USER_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		var permissions = []models.Permission{}
		for _, permit := range role.Permissions {
			permId := &models.PermissionById{ID: permit.Hex()}
			permission, err := permissionService.FetchPermissionById(permId)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					c.Status(http.StatusNotFound)
					return c.JSON(presenters.ErrorResponse(messages.PERMISSION_NOT_FOUND_ERROR_MESSAGE))
				}
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.PERMISSION_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
			permissions = append(permissions, *permission)
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.RoleDetailSuccessResponse(role, user, &permissions))
	}
}

func UpdateRole(roleService roles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleSlug := &models.RoleBySlug{Slug: c.Params("roleSlug")}
		role, err := roleService.FetchRoleBySlug(roleSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.ROLE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}
		var updateReq = &models.RoleRequest{}
		if err := c.BodyParser(updateReq); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_BODY_PARSER_ERROR_MESSAGE))
		}

		currentTimeStr := fmt.Sprintf("%v", time.Now().Unix())
		newSlug := updateReq.Name + "-" + currentTimeStr

		updateBody := &models.Role{
			ID:          role.ID,
			Owner:       role.Owner,
			Name:        updateReq.Name,
			Slug:        newSlug,
			Description: updateReq.Description,
			Permissions: updateReq.Permissions,
			CreatedAt:   role.CreatedAt,
		}

		updatedRole, err := roleService.UpdateRole(updateBody)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_UPDATE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.RoleSuccessResponse(updatedRole))
	}
}

func DeleteRole(roleService roles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleSlug := &models.RoleBySlug{Slug: c.Params("roleSlug")}
		role, err := roleService.FetchRoleBySlug(roleSlug)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.ROLE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		deletedRole, err := roleService.DeleteRole(role)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_DELETE_ERROR_MESSAGE))
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.RoleSuccessResponse(deletedRole))
	}
}

func GetMyRoles(roleService roles.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser := c.Locals("currentUser").(*models.User)
		roleOwner := &models.RoleByOwner{Owner: currentUser.Username}
		roles, err := roleService.FetchRolesByOwner(roleOwner)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.Status(http.StatusNotFound)
				return c.JSON(presenters.ErrorResponse(messages.ROLE_NOT_FOUND_ERROR_MESSAGE))
			}
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.ROLE_FAIL_TO_FETCH_ERROR_MESSAGE))
		}

		c.Status(http.StatusNotFound)
		return c.JSON(presenters.RolesSuccessResponse(roles))
	}
}
