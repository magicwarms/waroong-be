package middlewares

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"waroong-be/apps/constants"
	"waroong-be/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func validateToken(clientToken string, ctx *fiber.Ctx) (*jwt.Token, error) {
	if clientToken == "" {
		return nil, errors.New("token authorization can't empty")
	}

	extractedToken := strings.Split(clientToken, "Bearer ")
	if len(extractedToken) != 2 {
		return nil, errors.New("incorrect format authorization")
	}
	clientToken = strings.TrimSpace(extractedToken[1])

	token, err := jwt.Parse(clientToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("unexpected signing method: ", token.Header["alg"])
			return nil, errors.New("unexpected signing authorization method ")
		}
		return []byte(config.GoDotEnvVariable("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func CheckSuperadminAuthorization(ctx *fiber.Ctx) error {
	clientToken := ctx.Get("Authorization")

	token, errValidateToken := validateToken(clientToken, ctx)

	if errValidateToken != nil {
		return config.ErrorResponse(errValidateToken, ctx)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// get the payloads
		id := claims["id"]
		userTypeId := claims["userTypeId"]
		expiresAt := claims["expiresAt"]

		// Parse the input string as a float64
		parseFloat, errParseFloat := strconv.ParseFloat(fmt.Sprint(expiresAt), 64)
		if errParseFloat != nil {
			fmt.Println("errParseFloat:", errParseFloat)
			return config.ErrorResponse(errParseFloat, ctx)
		}

		isExpires := checkAuthorizationExpired(parseFloat, ctx)
		if isExpires {
			return config.ErrorResponse(errors.New("token has expired"), ctx)
		}

		if fmt.Sprint(userTypeId) != fmt.Sprint(constants.SUPERADMIN_USER_ROLE) {
			return config.ErrorResponse(errors.New("you don't have enough permission"), ctx)
		}

		// set the value into header
		ctx.Set("userId", fmt.Sprint(id))
		ctx.Set("userTypeId", fmt.Sprint(userTypeId))
	}

	return ctx.Next()
}

func CheckAuthorization(ctx *fiber.Ctx) error {
	clientToken := ctx.Get("Authorization")
	token, errValidateToken := validateToken(clientToken, ctx)

	if errValidateToken != nil {
		return config.ErrorResponse(errValidateToken, ctx)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// get the payloads
		id := claims["id"]
		userTypeId := claims["userTypeId"]
		expiresAt := claims["expiresAt"]

		// Parse the input string as a float64
		parseFloat, errParseFloat := strconv.ParseFloat(fmt.Sprint(expiresAt), 64)
		if errParseFloat != nil {
			fmt.Println("errParseFloat:", errParseFloat)
			return config.ErrorResponse(errParseFloat, ctx)
		}

		isExpires := checkAuthorizationExpired(parseFloat, ctx)
		if isExpires {
			return config.ErrorResponse(errors.New("token has expired"), ctx)
		}
		// set the value into header
		ctx.Set("userId", fmt.Sprint(id))
		ctx.Set("userTypeId", fmt.Sprint(userTypeId))
	}

	return ctx.Next()
}

func checkAuthorizationExpired(parseFloat float64, ctx *fiber.Ctx) bool {
	// Convert float64 to integer
	timestamp := int64(math.Round(parseFloat))
	timeNow := time.Now().Local().Unix()

	return timeNow > timestamp
}
