package transaction

import (
	"errors"
	"strconv"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes() (string, func(router fiber.Router)) {
	transactionController := NewTransactionController()
	return "/transactions", func(router fiber.Router) {
		router.Get("/", middleware.DeserializeUser, transactionController.GetAll)
		router.Get("/:id", middleware.DeserializeUser, transactionController.GetOne)
		router.Post("/", middleware.DeserializeUser, transactionController.Post)
		router.Patch("/:id", middleware.DeserializeUser, transactionController.Patch)
	}
}

type TransactionSchema struct {
	IdCategory  uint64    `json:"categoryId" validate:"required"`
	IdAccount   int       `json:"accountId" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Value       float32   `json:"value" validate:"required"`
	Date        time.Time `json:"date" validate:"required"`
	IdInvoice   *string   `json:"invoiceId"`
}

type TransactionUpdateSchema struct {
	IdCategory  *uint64    `json:"categoryId"`
	IdAccount   *int       `json:"accountId"`
	Description *string    `json:"description"`
	Value       *float32   `json:"value"`
	Date        *time.Time `json:"date"`
	IdInvoice   *string    `json:"invoiceId"`
}

type TransactionController struct {
	TransactionService
}

func NewTransactionController() *TransactionController {
	return &TransactionController{
		generic.NewGenericService[model.Transaction](),
	}
}

func (cc *TransactionController) GetAll(c *fiber.Ctx) error {
	accounts, err := cc.FindAll()

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(accounts), "items": accounts})
}

func (cc *TransactionController) GetOne(c *fiber.Ctx) error {

	accountId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	account, err := cc.FindById(accountId)

	if err != nil {
		var errorNotFound *util.ItemNotFoundError
		if errors.As(err, &errorNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": errorNotFound.Message})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "item": account})
}

func (cc *TransactionController) Post(c *fiber.Ctx) error {
	var payload *TransactionSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if errors := util.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	loggedUser := c.Locals("user").(model.User)

	newTransaction := model.Transaction{
		IdUser:      loggedUser.ID,
		Description: payload.Description,
		Date:        payload.Date,
		IdCategory:  payload.IdCategory,
		IdAccount:   payload.IdAccount,
		IdInvoice:   payload.IdInvoice,
		Value:       payload.Value,
	}

	if err := cc.Create(&newTransaction); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (cc *TransactionController) Patch(c *fiber.Ctx) error {
	var payload *TransactionUpdateSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})

	if payload.Description != nil {
		updates["description"] = payload.Description
	}

	if payload.Date != nil {
		updates["date"] = payload.Date
	}

	if payload.Value != nil {
		updates["value"] = payload.Value
	}

	accountId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	if err = cc.Update(accountId, updates); err != nil {
		var errorNotFound *util.ItemNotFoundError
		if errors.As(err, &errorNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": errorNotFound.Message})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": err})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
