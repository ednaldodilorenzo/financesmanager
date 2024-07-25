package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetRoutes(controller AuthController, deserializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/auth", func(router fiber.Router) {
		router.Post("/login", controller.SigninUser)
		router.Post("/signup", controller.SignUpUser)
		router.Get("/logout", deserializer.DeserializeUser, controller.LogoutUser)
	}
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type SignUpInput struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type SignInResponse struct {
	ID    uint64 `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

type AuthController interface {
	SigninUser(c *fiber.Ctx) error
	SignUpUser(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
}

type AuthControllerStruct struct {
	authService AuthService
}

func NewAuthController(authService AuthService) AuthController {
	return &AuthControllerStruct{
		authService: authService,
	}
}

//func BuildAuthController(service AuthService) *AuthController {
//	return &AuthController{
//		authService: service,
//	}
//}

func (a *AuthControllerStruct) SigninUser(c *fiber.Ctx) error {

	var payload *SignInInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := util.ValidateStruct(payload)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user, err := a.authService.ExecuteAuthentication(payload.Email, payload.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "Invalid email or Password"})
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)

	claims["sub"] = user.ID
	claims["exp"] = now.Add(515151500000000000).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte("Xuxa"))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	response := SignInResponse{
		ID:    user.ID,
		Name:  user.Name,
		Token: tokenString,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (a *AuthControllerStruct) SignUpUser(c *fiber.Ctx) error {
	var payload *SignUpInput

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	errors := util.ValidateStruct(payload)

	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": errors})
	}

	if payload.Password != payload.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Passwords do not match"})
	}

	newUser := &model.User{
		Name:     payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
	}

	if a.authService == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Could not save new user"})
	}

	if err := a.authService.RegisterUser(newUser); err != nil && strings.Contains(err.Error(), "duplicate key value violates unique") {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"status": "fail", "message": "User with that email already exists"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Something bad happened"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (a *AuthControllerStruct) LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
