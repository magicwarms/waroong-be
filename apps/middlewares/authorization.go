package middlewares

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func CheckAuthorization(c *fiber.Ctx) error {
	fmt.Println("Authorization called")
	return c.Next()
}
