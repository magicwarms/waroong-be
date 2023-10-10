package routers

import (
	"waroong-be/apps/user_types"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	// "waroong-be/apps/user_profiles"
	// "waroong-be/apps/user_types"
	"waroong-be/apps/users"
)

func Dispatch(DBConnection *gorm.DB, apiV1 fiber.Router) {
	// REGISTER ALL YOUR REPO, SERVICE, AND HANDLER HERE

	apiV1.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
	})

	// USER
	userRepo := users.NewRepo(DBConnection)
	userService := users.NewService(userRepo)
	users.NewUserHandler(apiV1.Group("/users"), userService)

	// // USER_TYPE
	user_types.NewRepo(DBConnection)

	// // USER_PROFILE
	// user_profiles.NewRepo(DBConnection)
}
