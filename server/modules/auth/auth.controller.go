package auth

import (
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

type authController struct {
	authService AuthService
	settings    *config.Settings
}

func NewAuthController(authService AuthService, settings *config.Settings) AuthController {
	return &authController{
		authService: authService,
		settings:    settings,
	}
}

func (a *authController) SigninUser(c *fiber.Ctx) error {

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
		return err
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

	return util.SendData(c, "success", &response, fiber.StatusOK)
}

func (a *authController) SignUpUser(c *fiber.Ctx) error {

	payload, err := util.ValidateRequestPayload[SignUpInput](c.BodyParser)

	if err != nil {
		return err
	}

	if err = a.authService.RegisterUserWithToken(payload); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (a *authController) LogoutUser(c *fiber.Ctx) error {
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

func (a *authController) StartRegistration(c *fiber.Ctx) error {

	payload, err := util.ValidateRequestPayload[StartRegistrationRequest](c.BodyParser)

	if err != nil {
		return err
	}

	if err = a.authService.StartRegistrationProcess(payload.Email); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func (a *authController) ChangePassword(c *fiber.Ctx) error {

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

func (a *authController) RecoverPassword(c *fiber.Ctx) error {
	payload, err := util.ValidateRequestPayload[StartRegistrationRequest](c.BodyParser)

	if err != nil {
		return err
	}

	if err := a.authService.StartRecoverPasswordProcess(payload.Email); err != nil {
		return err
	}

	return util.SendData[any](c, "success", nil, fiber.StatusOK)
}

func (a *authController) RedefinePassword(c *fiber.Ctx) error {
	var payload *RedefinePasswordRequest

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	validationErrors := util.ValidateStruct(payload)

	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "errors": validationErrors})
	}

	if err := a.authService.RedefinePassword(payload); err != nil {
		return err
	}

	return util.SendData[any](c, "success", nil, fiber.StatusOK)
}
