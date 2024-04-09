package transaction

import (
	"strconv"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller TransactionController) (string, func(router fiber.Router)) {
	return "/transactions", func(router fiber.Router) {
		router.Get("/", middleware.DeserializeUser, controller.GetAllWithRelationships)
		router.Get("/:id", middleware.DeserializeUser, controller.GetOne)
		router.Post("/", middleware.DeserializeUser, controller.Post)
		router.Patch("/:id", middleware.DeserializeUser, controller.Patch)
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

type TransactionController interface {
	generic.GenericController[*model.Transaction]
	GetAllWithRelationships(c *fiber.Ctx) error
}

type TransactionControllerStruct struct {
	generic.GenericController[*model.Transaction]
	TransactionService
}

func NewTransactionController(service TransactionService, controller generic.GenericController[*model.Transaction]) TransactionController {
	return &TransactionControllerStruct{
		controller,
		service,
	}
}

func (cc *TransactionControllerStruct) GetAllWithRelationships(c *fiber.Ctx) error {

	var monthParam, yearParam *int = nil, nil

	if month := c.Query("month"); month != "" {
		if value, err := strconv.Atoi(month); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "invalid parameter month value"})
		} else {
			monthParam = &value
		}
	}

	if year := c.Query("year"); year != "" {
		if value, err := strconv.Atoi(year); err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"status": "error", "message": "invalid parameter year value"})
		} else {
			yearParam = &value
		}
	}

	items, err := cc.FindAllRelated(monthParam, yearParam)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(items), "items": items})
}
