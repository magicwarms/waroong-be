package users

import (
	"errors"
	"strconv"
	"waroong-be/apps/middlewares"
	"waroong-be/apps/users/entity"
	"waroong-be/apps/users/interfaces"
	"waroong-be/config"

	"github.com/gofiber/fiber/v2"
)

func NewUserHandler(user fiber.Router, userService interfaces.UserService) {
	user.Post("/store", AddNewUser(userService))
	user.Get("/superadmin/all", middlewares.CheckSuperadminAuthorization, GetAllSuperadminUsers(userService))
	user.Post("/login", Login(userService))
	user.Get("/profile", middlewares.CheckAuthorization, GetUserProfile(userService))
	user.Patch("/superadmin/change_password/user", middlewares.CheckSuperadminAuthorization, ChangeSuperadminPassword(userService))
	user.Patch("/update", middlewares.CheckAuthorization, UpdateUser(userService))
	// user.Delete("/delete", middlewares.CheckAuthorization, DeleteUser(userService))
	user.Post("/forgot_password", ForgotPasswordUser(userService))
	user.Patch("/change_forgot_password", ChangeForgotPasswordUser(userService))
}

// AddNewUser is store user superadmin data into database
func AddNewUser(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userDTO := &entity.AddUserRequestDTO{
			Email:      c.FormValue("email"),
			Password:   c.FormValue("password"),
			FirstName:  c.FormValue("firstname"),
			LastName:   c.FormValue("lastname"),
			Phone:      c.FormValue("phone"),
			UserTypeId: c.FormValue("user_type_id"),
		}

		if err := c.BodyParser(userDTO); err != nil {
			return config.ErrorResponse(err, c)
		}

		validationErr := config.ValidateFields(*userDTO)
		if validationErr != nil {
			return config.ValidateResponse(validationErr, c)
		}

		err := userService.StoreUser(userDTO)
		if err != nil {
			return config.ErrorResponse(err, c)
		}
		return config.AppResponse(nil, c)
	}
}

// UpdateUser is update user data into database
func UpdateUser(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userUpdateDTO := &entity.UpdateUserRequestDTO{
			ID:        c.GetRespHeader("userId"),
			FirstName: c.FormValue("firstname"),
			LastName:  c.FormValue("lastname"),
			Phone:     c.FormValue("phone"),
		}

		if err := c.BodyParser(userUpdateDTO); err != nil {
			return config.ErrorResponse(err, c)
		}

		validationErr := config.ValidateFields(*userUpdateDTO)
		if validationErr != nil {
			return config.ValidateResponse(validationErr, c)
		}

		err := userService.UpdateUser(userUpdateDTO)
		if err != nil {
			return config.ErrorResponse(err, c)
		}
		return config.AppResponse(nil, c)
	}
}

// GetAllSuperadminUsers is to get all users data
func GetAllSuperadminUsers(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		getAllUsers, err := userService.FindAllSuperadminUsers()
		if err != nil {
			return config.ErrorResponse(err, c)
		}
		return config.AppResponse(getAllUsers, c)
	}
}

// Login is to get all users data
func Login(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userLoginDTO := &entity.UserLoginRequestDTO{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}

		if err := c.BodyParser(userLoginDTO); err != nil {
			return config.ErrorResponse(err, c)
		}

		validationErr := config.ValidateFields(*userLoginDTO)
		if validationErr != nil {
			return config.ValidateResponse(validationErr, c)
		}

		userLogin, err := userService.LoginUser(userLoginDTO)

		if err != nil {
			return config.ErrorResponse(err, c)
		}

		return config.AppResponse(userLogin, c)
	}
}

// GetUserProfile is to get spesific user data by user ID
func GetUserProfile(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get userId from token payload
		userId := c.GetRespHeader("userId")

		parseUserId, errParseUserId := strconv.ParseUint(userId, 10, 64)
		if errParseUserId != nil {
			return errors.New(errParseUserId.Error())
		}

		userIdDTO := &entity.UserGetIDRequestDTO{
			ID: userId,
		}

		validationErr := config.ValidateFields(*userIdDTO)
		if validationErr != nil {
			return config.ValidateResponse(validationErr, c)
		}

		user, err := userService.GetUserById(parseUserId)
		if err != nil {
			return config.ErrorResponse(err, c)
		}
		if user.Email == "" {
			user = nil
		}
		return config.AppResponse(&user, c)
	}
}

func ChangeSuperadminPassword(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// get userId from token payload
		userId := c.GetRespHeader("userId")

		changeUserPasswordDTO := &entity.ChangePasswordUserDTO{
			UserID:   userId,
			Password: c.FormValue("password"),
		}

		if err := c.BodyParser(changeUserPasswordDTO); err != nil {
			return config.ErrorResponse(err, c)
		}

		validationErr := config.ValidateFields(*changeUserPasswordDTO)
		if validationErr != nil {
			return config.ValidateResponse(validationErr, c)
		}

		err := userService.UpdateSuperadminPassword(changeUserPasswordDTO)

		if err != nil {
			return config.ErrorResponse(err, c)
		}

		return config.AppResponse(nil, c)

	}
}

// ForgotPasswordUser is store user superadmin data into database
func ForgotPasswordUser(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		forgotPasswordDTO := &entity.ForgotPasswordRequestDTO{
			Email: c.FormValue("email"),
		}

		if err := c.BodyParser(forgotPasswordDTO); err != nil {
			return config.ErrorResponse(err, c)
		}

		validationErr := config.ValidateFields(*forgotPasswordDTO)
		if validationErr != nil {
			return config.ValidateResponse(validationErr, c)
		}

		forgotPassword, err := userService.ForgotPassword(forgotPasswordDTO)
		if err != nil {
			return config.ErrorResponse(err, c)
		}
		return config.AppResponse(forgotPassword, c)
	}
}

// ChangeForgotPasswordUser is to change the user password after forgot password
func ChangeForgotPasswordUser(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		changeForgotPasswordDTO := &entity.ChangeForgotPasswordRequestDTO{
			Password: c.FormValue("password"),
			Token:    c.FormValue("token"),
		}

		if err := c.BodyParser(changeForgotPasswordDTO); err != nil {
			return config.ErrorResponse(err, c)
		}

		validationErr := config.ValidateFields(*changeForgotPasswordDTO)
		if validationErr != nil {
			return config.ValidateResponse(validationErr, c)
		}

		forgotPassword, err := userService.ChangeForgotPassword(changeForgotPasswordDTO)
		if err != nil {
			return config.ErrorResponse(err, c)
		}
		return config.AppResponse(forgotPassword, c)
	}
}

// // UpdateBank is update user data into database
// func UpdateBank(bankService interfaces.BankService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		var licenseId string = c.GetRespHeader("LicenseId")
// 		isActive, err := strconv.ParseBool(c.FormValue("is_active"))
// 		if err != nil {
// 			return config.ErrorResponse(err, c)
// 		}
// 		bankDTO := &entity.BankUpdateRequestDTO{
// 			ID:            c.FormValue("id"),
// 			Name:          c.FormValue("name"),
// 			LicenseId:     licenseId,
// 			AccountName:   c.FormValue("account_name"),
// 			AccountNumber: c.FormValue("account_number"),
// 			IsActive:      &isActive,
// 		}

// 		if err := c.BodyParser(bankDTO); err != nil {
// 			return config.ErrorResponse(err, c)
// 		}

// 		validationErr := config.ValidateFields(*bankDTO)
// 		if validationErr != nil {
// 			return config.ValidateResponse(validationErr, c)
// 		}

// 		errUpdate := bankService.UpdateBank(bankDTO)
// 		if errUpdate != nil {
// 			return config.ErrorResponse(errUpdate, c)
// 		}
// 		return config.AppResponse(nil, c)
// 	}
// }

// // DeleteUser is delete branch data in database
// func DeleteUser(userService interfaces.UserService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		idStr := c.FormValue("id")
// 		id, err := strconv.ParseUint(idStr, 10, 32)
// 		if err != nil {
// 			return config.ErrorResponse(err, c)
// 		}

// 		userDTO := &entity.UserGetIDRequestDTO{
// 			ID: uint(id),
// 		}

// 		if err := c.BodyParser(userDTO); err != nil {
// 			return config.ErrorResponse(err, c)
// 		}

// 		validationErr := config.ValidateFields(*userDTO)
// 		if validationErr != nil {
// 			return config.ValidateResponse(validationErr, c)
// 		}

// 		deleteUser, err := userService.DeleteUser(uint(id))
// 		if err != nil {
// 			return config.ErrorResponse(err, c)
// 		}

// 		return config.AppResponse(deleteUser, c)
// 	}
// }
