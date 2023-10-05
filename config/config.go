package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var validate = validator.New()

// GoDotEnvVariable at godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		panic("Error loading .env file")
	}
	return os.Getenv(key)
}

// PrettyPrint is make easier to print data result to console after querying data to db
func PrettyPrint(i interface{}) string {
	results, _ := json.MarshalIndent(i, "", "\t")
	return string(results)
}

// AppResponse is for response config show to Frontend side
func AppResponse(data interface{}, ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&DefaultResponse{
		Success: true,
		Data:    data,
		Errors:  "",
	})
}

// ErrorResponse is for default error response
func ErrorResponse(err error, ctx *fiber.Ctx) error {
	errMessage := err.Error()

	errType := fiber.StatusInternalServerError
	isSuccess := false

	if strings.Contains(errMessage, "existed") {
		errType = fiber.StatusConflict
	}

	if strings.Contains(errMessage, "signature") || strings.Contains(errMessage, "authorization") {
		errType = fiber.StatusUnauthorized
	}

	if strings.Contains(errMessage, "not found") {
		errType = fiber.StatusOK
		isSuccess = true
	}

	if strings.Contains(errMessage, "Unprocessable Entity") {
		errType = fiber.StatusUnprocessableEntity
	}

	if errType == fiber.StatusInternalServerError {
		errMessage = "Something went wrong"
	}

	return ctx.Status(errType).JSON(&DefaultResponse{
		Success: isSuccess,
		Data:    nil,
		Errors:  errMessage,
	})
}

// ValidateResponse is for response error
func ValidateResponse(responses []*ResponseField, ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(&ValidationDefaultResponse{
		Success: false,
		Data:    nil,
		Errors:  responses,
	})
}

// ValidateFields is for response validation error
func ValidateFields(model interface{}) []*ResponseField {
	var validationErrors []*ResponseField
	err := validate.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {

			// fmt.Println(err.Namespace())
			// fmt.Println(err.Field())
			// fmt.Println(err.StructNamespace())
			// fmt.Println(err.StructField())
			// fmt.Println(err.Tag())
			// fmt.Println(err.ActualTag())
			// fmt.Println(err.Kind())
			// fmt.Println(err.Type())
			// fmt.Println(err.Value())
			// fmt.Println(err.Param())
			// fmt.Println()

			var element ResponseField
			element.Field = err.Field()
			element.Message = messageForTag(err.Tag(), err.Param(), err.Field())
			element.Value = err.Value().(string)
			validationErrors = append(validationErrors, &element)
		}
	}

	return validationErrors
}

func messageForTag(tag, value, field string) string {
	switch tag {
	case "required":
		return field + " wajib dimasukkan"
	case "email":
		return field + " tidak valid"
	case "min":
		return field + " minimal " + value + " karakter"
	case "max":
		return field + " maksimal " + value + " karakter"
	case "lowercase":
		return field + " harus kecil semua"
	case "e164":
		return "Format " + field + " tidak valid"
	case "uuid4":
		return field + " UUID tidak valid"
	case "latitude":
		return field + value + " tidak valid"
	case "longitude":
		return field + value + " tidak valid"
	case "numeric":
		return field + value + " harus berupa angka"
	case "alpha":
		return field + value + " harus berupa huruf saja"
	case "alphanum":
		return field + " harus berupa angka dan huruf"
	}

	// default error
	return strings.ToUpper("WARNING!! pesan validasi belum di convert. type: " + tag)
}
