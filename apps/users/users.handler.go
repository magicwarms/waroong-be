package users

import (
	"waroong-be/apps/middlewares"
	"waroong-be/apps/users/entity"
	"waroong-be/apps/users/interfaces"
	"waroong-be/config"

	"github.com/gofiber/fiber/v2"
)

// TO-DO
// add middleware later
func NewUserHandler(user fiber.Router, userService interfaces.UserService) {
	user.Post("/store", AddNewUser(userService))
	user.Get("/superadmin/all", middlewares.CheckAuthorization, GetAllSuperadminUsers(userService))
	user.Post("/login", Login(userService))
	// bank.Get("/get/:id", middlewares.CheckAuthorization, GetBankById(userService))
	// bank.Patch("/update", middlewares.CheckAuthorization, UpdateBank(userService))
	// user.Delete("/delete", middlewares.CheckAuthorization, DeleteUser(userService))
	// user.Post("/customer/store", middlewares.CheckAuthorization, AddNewCustomer(userService))
}

// AddNewUser is store user superadmin data into database
func AddNewUser(userService interfaces.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userDTO := &entity.UserRequestDTO{
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
		config.PrettyPrint(err)
		if err != nil {
			return config.ErrorResponse(err, c)
		}

		return config.AppResponse(userLogin, c)
	}
}

// // AddNewUser is store user data into database
// func AddNewCustomer(userService interfaces.UserService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		userDTO := &entity.UserRequestDTO{
// 			Email:     c.FormValue("email"),
// 			Password:  c.FormValue("password"),
// 			FirstName: c.FormValue("firstname"),
// 			LastName:  c.FormValue("lastname"),
// 		}

// 		if err := c.BodyParser(userDTO); err != nil {
// 			return config.ErrorResponse(err, c)
// 		}

// 		validationErr := config.ValidateFields(*userDTO)
// 		if validationErr != nil {
// 			return config.ValidateResponse(validationErr, c)
// 		}

// 		err := userService.StoreUser(userDTO)
// 		if err != nil {
// 			return config.ErrorResponse(err, c)
// 		}
// 		return config.AppResponse(nil, c)
// 	}
// }

// // GetBankById is to get spesific bank data by bank ID
// func GetBankById(bankService interfaces.BankService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// get bankId from header value
// 		var id string = c.Params("id")

// 		bankDTO := &entity.BankGetIDRequestDTO{
// 			ID: id,
// 		}

// 		validationErr := config.ValidateFields(*bankDTO)
// 		if validationErr != nil {
// 			return config.ValidateResponse(validationErr, c)
// 		}

// 		getBank, err := bankService.GetBankById(id)
// 		if err != nil {
// 			return config.ErrorResponse(err, c)
// 		}
// 		if getBank.ID == "" {
// 			getBank = nil
// 		}
// 		return config.AppResponse(&getBank, c)
// 	}
// }

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
