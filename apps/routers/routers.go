package routers

import (
	"waroong-be/apps/user_profiles"
	"waroong-be/apps/user_types"
	"waroong-be/apps/users"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Dispatch(DBConnection *gorm.DB, apiV1 fiber.Router) {
	// REGISTER ALL YOUR REPO, SERVICE, AND HANDLER HERE

	apiV1.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
	})

	// REPOSITORIES
	userRepo := users.NewRepo(DBConnection)
	userTypeRepo := user_types.NewRepo(DBConnection)
	user_profiles.NewRepo(DBConnection)

	// SERVICES
	userTypeService := user_types.NewService(userTypeRepo)
	userService := users.NewService(userRepo, userTypeService)

	// HANDLERS
	users.NewUserHandler(apiV1.Group("/users"), userService)
}
