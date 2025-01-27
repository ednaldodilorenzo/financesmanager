package budget

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller BudgetController, desirializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/v1/budgets", func(router fiber.Router) {
		router.Get("/", desirializer.DeserializeUser, controller.GetAllByYear)
		router.Get("/:id", desirializer.DeserializeUser, controller.GetOne)
		router.Post("/", desirializer.DeserializeUser, controller.Post)
		router.Patch("/:id", desirializer.DeserializeUser, controller.Patch)
		router.Delete("/:id", desirializer.DeserializeUser, controller.Delete)
	}
}

type BudgetController interface {
	generic.GenericController[*model.Budget]
	GetAllByYear(c *fiber.Ctx) error
}

type BudgetControllerStruct struct {
	generic.GenericController[*model.Budget]
	service BudgetService
}

func NewBudgetController(controller generic.GenericController[*model.Budget], service BudgetService) BudgetController {
	return &BudgetControllerStruct{
		controller,
		service,
	}
}

func (b *BudgetControllerStruct) GetAllByYear(c *fiber.Ctx) error {
	year := c.QueryInt("year")

	items, err := b.service.FindAllByYear(year)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(items), "items": items})
}
