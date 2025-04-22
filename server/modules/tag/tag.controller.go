package tag

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller TagController, deserializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/v1/tags", func(router fiber.Router) {
		router.Get("/", deserializer.DeserializeUser, controller.GetAll)
	}
}

type TagController interface {
	GetAll(c *fiber.Ctx) error
}

type tagController struct {
	service TagService
}

func NewTagController(service TagService) TagController {
	return &tagController{
		service,
	}
}

func (t *tagController) GetAll(c *fiber.Ctx) error {

	filter := c.Query("filter")
	loggedUser := c.Locals("user").(model.User)
	items, err := t.service.FindAllWithFilter(filter, int(loggedUser.ID))

	if err != nil {
		return err
	}

	return util.SendData(c, "success", &items, fiber.StatusOK)
}
