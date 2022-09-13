package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
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
			return c.JSON(presenters.ErrorResponse(messages.ROLE_COVER_FAIL_TO_INSERT_ERROR_MESSAGE))
		}
		c.Status(http.StatusOK)
		return c.JSON(presenters.)
	}
}
