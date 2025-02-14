package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
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
		router.Post("/register", controller.StartRegistration)
		router.Post("/changePassword", deserializer.DeserializeUser, controller.ChangePassword)
	}
}

type SignInInput struct {
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password"  validate:"required"`
}

type SignUpInput struct {
	Name            string `json:"name" validate:"required"`
	Token           string `json:"token" validate:"required"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8"`
}

type SignInResponse struct {
	ID    uint64 `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Token string `json:"token,omitempty"`
}

type StartRegistrationRequest struct {
	Email string `json:"email" validate:"required"`
}

type ChangePasswordRequest struct {
	Password          string `json:"password" validate:"required"`
	NewPassword       string `json:"newPassword" validate:"required"`
	CofirmNewPassword string `json:"confirmNewPassword" validate:"required"`
}

type AuthController interface {
	SigninUser(c *fiber.Ctx) error
	SignUpUser(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
	StartRegistration(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
}

type AuthControllerStruct struct {
	authService AuthService
	settings    *config.Settings
}

func NewAuthController(authService AuthService, settings *config.Settings) AuthController {
	return &AuthControllerStruct{
		authService: authService,
		settings:    settings,
	}
}

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

	tokenString, err := tokenByte.SignedString([]byte(a.settings.AppSettings.JwtKey))

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

	validationErrors := util.ValidateStruct(payload)

	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": validationErrors})
	}

	if payload.Password != payload.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": "Passwords do not match"})
	}

	var businessError *util.BusinessError

	if err := a.authService.RegisterUserWithToken(payload); err != nil && errors.As(err, &businessError) {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": businessError.Message})
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

func (a *AuthControllerStruct) StartRegistration(c *fiber.Ctx) error {
	var payload *StartRegistrationRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	err := a.authService.StartRegistrationProcess(payload.Email)

	if err != nil {
		var businessError *util.BusinessError
		if errors.As(err, &businessError) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "message": businessError.Message})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": "Failed to start registration process"})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (a *AuthControllerStruct) ChangePassword(c *fiber.Ctx) error {
	var payload *ChangePasswordRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	validationErrors := util.ValidateStruct(payload)

	if validationErrors != nil || (payload.NewPassword != payload.CofirmNewPassword) || (payload.NewPassword == payload.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": validationErrors})
	}

	loggedUser := c.Locals("user").(model.User)

	if err := a.authService.ChangePassword(int(loggedUser.ID), payload.NewPassword); err != nil {
		var validationError *util.BusinessError
		if errors.As(err, &validationError) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": validationError.Message, "code": validationError.Code})
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": "failed to change password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
