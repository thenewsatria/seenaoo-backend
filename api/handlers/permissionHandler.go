package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/presenters"
	"github.com/thenewsatria/seenaoo-backend/pkg/models"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/variables/messages"
)

func GetAvailablePermissions(permissionService permissions.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		permissionByItemCategories := map[string]*[]models.Permission{}
		itemCategories, err := permissionService.FetchPermissionsDistinctItemCategory()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenters.ErrorResponse(messages.PERMISSION_FAIL_TO_FETCH_DISTINCT_ITEM_CATEGORY))
		}

		for _, category := range *itemCategories {
			itemCategory := &models.PermissionByItemCategory{ItemCategory: category}
			permissions, err := permissionService.FetchPermissionsByItemCategory(itemCategory)
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return c.JSON(presenters.ErrorResponse(messages.PERMISSION_FAIL_TO_FETCH_ERROR_MESSAGE))
			}
			permissionByItemCategories[category] = permissions
		}

		c.Status(http.StatusOK)
		return c.JSON(presenters.PermissionsGroupedByItemCateogory(permissionByItemCategories))
	}
}
