package account

import (
	"strconv"

	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes() (string, func(router fiber.Router)) {
	accountController := NewAccountController()
	return "/accounts", func(router fiber.Router) {
		router.Get("/", middleware.DeserializeUser, accountController.GetAll)
		router.Get("/:id", middleware.DeserializeUser, accountController.GetOne)
		router.Post("/", middleware.DeserializeUser, accountController.Post)
		router.Patch("/:id", middleware.DeserializeUser, accountController.Patch)
	}
}

type AccountSchema struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
}

type AccountUpdateSchema struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AccountController struct {
	generic.GenericService[model.Account]
}

func NewAccountController() *AccountController {
	return &AccountController{
		generic.NewGenericService[model.Account](),
	}
}

func (cc *AccountController) GetAll(c *fiber.Ctx) error {
	accounts, err := cc.FindAll()

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(accounts), "items": accounts})
}

func (cc *AccountController) GetOne(c *fiber.Ctx) error {

	accountId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	account, err := cc.FindById(accountId)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "item": account})
}

func (cc *AccountController) Post(c *fiber.Ctx) error {
	var payload *AccountSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if errors := util.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	loggedUser := c.Locals("user").(model.User)

	newAccount := model.Account{
		Name:   payload.Name,
		Type:   payload.Type,
		IdUser: loggedUser.ID,
	}

	if err := cc.Create(&newAccount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (cc *AccountController) Patch(c *fiber.Ctx) error {
	var payload *AccountUpdateSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})

	if payload.Name != "" {
		updates["name"] = payload.Name
	}

	if payload.Type != "" {
		updates["type"] = payload.Type
	}

	accountId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	if err = cc.Update(accountId, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
