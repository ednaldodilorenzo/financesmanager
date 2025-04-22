package planning

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
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

type planningController struct {
	service PlanningService
}

func NewPlanningController(service PlanningService) PlanningController {
	return &planningController{
		service: service,
	}
}

func (p *planningController) GetPlanning(c *fiber.Ctx) error {
	month := c.QueryInt("month")
	year := c.QueryInt("year")

	loggedUser := c.Locals("user").(model.User)
	items, err := p.service.FindByMonthAndYear(month, year, int(loggedUser.ID))

	if err != nil {
		return err
	}

	return util.SendData(c, "success", &items, fiber.StatusOK)
}
