package planning

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller PlanningController, deserializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/v1/plannings", func(router fiber.Router) {
		router.Get("/", deserializer.DeserializeUser, controller.GetPlanning)
	}
}

type PlanningController interface {
	GetPlanning(c *fiber.Ctx) error
}

type PlanningControllerStruct struct {
	service PlanningService
}

func NewPlanningController(service PlanningService) PlanningController {
	return &PlanningControllerStruct{
		service: service,
	}
}

func (p *PlanningControllerStruct) GetPlanning(c *fiber.Ctx) error {
	month := c.QueryInt("month")
	year := c.QueryInt("year")

	items, err := p.service.FindByMonthAndYear(month, year)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(items), "items": items})
}
