package routes

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thenewsatria/seenaoo-backend/api/middlewares"
	"github.com/thenewsatria/seenaoo-backend/database"
	"github.com/thenewsatria/seenaoo-backend/pkg/collaborations"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardcovers"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcardhints"
	"github.com/thenewsatria/seenaoo-backend/pkg/flashcards"
	"github.com/thenewsatria/seenaoo-backend/pkg/permissions"
	"github.com/thenewsatria/seenaoo-backend/pkg/refreshtokens"
	"github.com/thenewsatria/seenaoo-backend/pkg/roles"
	"github.com/thenewsatria/seenaoo-backend/pkg/tags"
	"github.com/thenewsatria/seenaoo-backend/pkg/tests"
	"github.com/thenewsatria/seenaoo-backend/pkg/userprofiles"
	"github.com/thenewsatria/seenaoo-backend/pkg/users"
)

func Router(app *fiber.App) {
	var flashcardCollection = database.UseDB().Collection(flashcards.CollectionName)
	var flashcardRepo = flashcards.NewRepo(flashcardCollection)
	var flashcardService = flashcards.NewService(flashcardRepo)

	var flashcardHintCollection = database.UseDB().Collection(flashcardhints.CollectionName)
	var flashcardHintRepo = flashcardhints.NewRepo(flashcardHintCollection)
	var flashcardHintService = flashcardhints.NewService(flashcardHintRepo)

	var userCollection = database.UseDB().Collection(users.CollectionName)
	var userRepo = users.NewRepo(userCollection)
	var userService = users.NewService(userRepo)

	var userProfileCollection = database.UseDB().Collection(userprofiles.CollectionName)
	var userProfileRepo = userprofiles.NewRepo(userProfileCollection)
	var userProfileService = userprofiles.NewService(userProfileRepo)

	var refreshTokenCollection = database.UseDB().Collection(refreshtokens.CollectionName)
	var refreshTokenRepo = refreshtokens.NewRepo(refreshTokenCollection)
	var refreshTokenService = refreshtokens.NewService(refreshTokenRepo)

	var collaborationCollection = database.UseDB().Collection(collaborations.CollectionName)
	var collaborationRepo = collaborations.NewRepo(collaborationCollection)
	var collaborationService = collaborations.NewService(collaborationRepo)

	var flashcardCoverCollection = database.UseDB().Collection(flashcardcovers.CollectionName)
	var flashcardCoverRepo = flashcardcovers.NewRepo(flashcardCoverCollection)
	var flashcardCoverService = flashcardcovers.NewService(flashcardCoverRepo)

	var tagCollection = database.UseDB().Collection(tags.CollectionName)
	var tagRepo = tags.NewRepo(tagCollection)
	var tagService = tags.NewService(tagRepo)

	var roleCollection = database.UseDB().Collection(roles.CollectionName)
	var roleRepo = roles.NewRepo(roleCollection)
	var roleService = roles.NewService(roleRepo)

	var permissionCollection = database.UseDB().Collection(permissions.CollectionName)
	var permissionRepo = permissions.NewRepo(permissionCollection)
	var permissionService = permissions.NewService(permissionRepo)

	var testCollection = database.UseDB().Collection(tests.CollectionName)
	var testRepo = tests.NewRepo(testCollection)
	var testService = tests.NewService(testRepo)

	api := app.Group("/api")
	apiV1 := api.Group("/v1")

	apiV1.Get("/", func(c *fiber.Ctx) error {
		c.Status(http.StatusOK)
		return c.JSON(fiber.Map{
			"success": true,
		})
	})

	apiV1.Static("/static", "./public")

	flashcardCoverRouter(apiV1, flashcardCoverService, flashcardService, flashcardHintService, tagService,
		userService, userProfileService, collaborationService, roleService, permissionService)
	flashcardRouter(apiV1, flashcardService, flashcardHintService, flashcardCoverService,
		userService, collaborationService, roleService, permissionService)
	flashcardHintRouter(apiV1, flashcardHintService, flashcardService, flashcardCoverService,
		userService, collaborationService, roleService, permissionService)
	authenticationRouter(apiV1, userService, refreshTokenService, userProfileService)
	roleRouter(apiV1, roleService, userService, permissionService)
	collaborationRouter(apiV1, collaborationService, userService, userProfileService,
		flashcardCoverService, roleService, permissionService)
	tagRouter(apiV1, tagService, flashcardCoverService)
	permissionRouter(apiV1, permissionService)
	testRouter(apiV1, testService)

	userRouter(apiV1, userProfileService, userService)

	testing := apiV1.Group("/testing")
	testing.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"passed_test_id": "sip",
		})
	})
	testing.Post("/", middlewares.TestMW2(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"passed_test_id": "sip2",
			"data":           c.Locals("passedParam"),
		})
	})
	testing.Delete("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"passed_test_id": "sip2",
			"data":           c.Locals("passedParam"),
		})
	})
	testing.Use("/:testId", middlewares.TestMW())
	testing.Get("/:testId", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"passed_test_id": c.Locals("testingID"),
		})
	})
	testing.Post("/:testId", middlewares.TestMW3(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"passed_test_id": c.Locals("testingID"),
			"passedParam":    c.Locals("passedParam"),
			"newParam":       c.Locals("mangstap"),
		})
	})
	testing.Delete("/:testId", func(c *fiber.Ctx) error {
		fmt.Print(c.Path())
		return c.JSON(fiber.Map{
			"passed_test_id": c.Locals("testingID"),
			"passedParam":    c.Locals("passedParam"),
			"newParam":       c.Locals("mangstap"),
		})
	})
	testing.Delete("/purge/:testId2", middlewares.TestMW2(), middlewares.TestMW3(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"passed_test_id": c.Locals("testingID"),
			"passedParam":    c.Locals("passedParam"),
			"newParam":       c.Locals("mangstap"),
		})
	})
	testing.Delete("/:testId/blabla", middlewares.TestMW3(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"passed_test_id": c.Locals("testingID"),
			"passedParam":    c.Locals("passedParam"),
			"newParam":       c.Locals("mangstap"),
		})
	})
}
