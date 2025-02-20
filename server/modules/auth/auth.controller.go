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
)

func GetRoutes(controller AuthController, deserializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/auth", func(router fiber.Router) {
		router.Post("/login", controller.SigninUser)
		router.Post("/signup", controller.SignUpUser)
		router.Get("/logout", deserializer.DeserializeUser, controller.LogoutUser)
		router.Post("/register", controller.StartRegistration)
		router.Post("/changePassword", deserializer.DeserializeUser, controller.ChangePassword)
		router.Post("/recoverPassword", controller.RecoverPassword)
		router.Post("/redefinePassword", controller.RedefinePassword)
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

type RedefinePasswordRequest struct {
	Token             string `json:"token" validate:"required"`
	Password          string `json:"password" validate:"required"`
	CofirmNewPassword string `json:"confirmPassword" validate:"required"`
}

type AuthController interface {
	SigninUser(c *fiber.Ctx) error
	SignUpUser(c *fiber.Ctx) error
	LogoutUser(c *fiber.Ctx) error
	StartRegistration(c *fiber.Ctx) error
	ChangePassword(c *fiber.Ctx) error
	RecoverPassword(c *fiber.Ctx) error
	RedefinePassword(c *fiber.Ctx) error
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

	payload, err := util.ValidateRequestPayload[SignInInput](c.BodyParser)

	if err != nil {
		return err
	}

	user, err := a.authService.ExecuteAuthentication(payload.Email, payload.Password)

	if err != nil {
		return err
	}

	expTime := 30 * time.Minute
	tokenString, err := util.GenerateToken(user.ID, &a.settings.AppSettings.JwtKey, &expTime)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("generating JWT Token failed: %v", err)})
	}

	response := SignInResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	secureCookie := false

	if a.settings.AppSettings.Environment == config.ENVIRONMENT_PROD {
		secureCookie = true
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    *tokenString,
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: true,         // Prevents JavaScript access
		Secure:   secureCookie, // Use true in production (HTTPS required)
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(response)
}

func (a *AuthControllerStruct) SignUpUser(c *fiber.Ctx) error {

	payload, err := util.ValidateRequestPayload[SignUpInput](c.BodyParser)

	if err != nil {
		return err
	}

	if err = a.authService.RegisterUserWithToken(payload); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (a *AuthControllerStruct) LogoutUser(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire immediately
		HTTPOnly: true,
		Secure:   true, // Use true in production
		SameSite: "Strict",
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}

func (a *AuthControllerStruct) StartRegistration(c *fiber.Ctx) error {

	payload, err := util.ValidateRequestPayload[StartRegistrationRequest](c.BodyParser)

	if err != nil {
		return err
	}

	if err = a.authService.StartRegistrationProcess(payload.Email); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (a *AuthControllerStruct) ChangePassword(c *fiber.Ctx) error {

	payload, err := util.ValidateRequestPayload[ChangePasswordRequest](c.BodyParser)

	if err != nil {
		return err
	}

	loggedUser := c.Locals("user").(model.User)

	if err := a.authService.ChangePassword(int(loggedUser.ID), payload); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (a *AuthControllerStruct) RecoverPassword(c *fiber.Ctx) error {
	var payload *StartRegistrationRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	validationErrors := util.ValidateStruct(payload)

	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": validationErrors})
	}

	if err := a.authService.StartRecoverPasswordProcess(payload.Email); err != nil {

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (a *AuthControllerStruct) RedefinePassword(c *fiber.Ctx) error {
	var payload *RedefinePasswordRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	validationErrors := util.ValidateStruct(payload)

	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": validationErrors})
	}

	if err := a.authService.RedefinePassword(payload); err != nil {
		var validationError *util.BusinessError
		var notFoundError *util.NotFoundError
		if errors.As(err, &validationError) {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": validationError.Message, "code": validationError.Code})
		} else if errors.As(err, &notFoundError) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "fail", "errors": "Usuário não encontrado"})
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"status": "fail", "errors": "failed to change password"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
