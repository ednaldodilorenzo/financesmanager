package auth

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/dto"
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
		router.Get("/google/callback", controller.GoogleRegistrationCallback)
		router.Get("/google/login", controller.SigninGoogle)
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
	GoogleRegistrationCallback(c *fiber.Ctx) error
	SigninGoogle(c *fiber.Ctx) error
}

type authController struct {
	authService AuthService
	settings    *config.Settings
}

var (
	oauthStateString = "random-state" // use a secure random string in production
)

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

	secureCookie := false

	if a.settings.AppSettings.Environment == config.ENVIRONMENT_PROD {
		secureCookie = true
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    user.Token,
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: true,         // Prevents JavaScript access
		Secure:   secureCookie, // Use true in production (HTTPS required)
		SameSite: "Strict",
	})

	return util.SendData(c, "success", &user, fiber.StatusOK)
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

func (a *authController) GoogleRegistrationCallback(c *fiber.Ctx) error {
	state := c.Query("state")

	if state != oauthStateString {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid OAuth state")
	}

	code := c.Query("code")
	token, err := a.settings.GoogleOAuthSettings.Exchange(context.Background(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Code exchange failed: " + err.Error())
	}

	client := a.settings.GoogleOAuthSettings.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed getting user info: " + err.Error())
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to read user info: " + err.Error())
	}

	var userInfo map[string]any
	if err := fiber.New().Config().JSONDecoder(bodyBytes, &userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse user info")
	}

	email := userInfo["email"].(string)
	oauthUser := dto.UserOAuthRegistrationRegistration{
		Name:  userInfo["name"].(string),
		Email: email,
	}

	currentUser, err := a.authService.RegisterUserOAuthUser(oauthUser)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token")
	}

	secureCookie := false

	if a.settings.AppSettings.Environment == config.ENVIRONMENT_PROD {
		secureCookie = true
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    currentUser.Token,
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: true,         // Prevents JavaScript access
		Secure:   secureCookie, // Use true in production (HTTPS required)
		SameSite: "Strict",
	})

	// Redirect to frontend with token (as query or fragment)
	return c.Redirect(fmt.Sprintf("%s/oauth-login", a.settings.AppSettings.Url))
}

func (a *authController) SigninGoogle(c *fiber.Ctx) error {
	url := a.settings.GoogleOAuthSettings.AuthCodeURL(oauthStateString)
	return c.Redirect(url)
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
