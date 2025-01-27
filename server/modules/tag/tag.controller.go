package tag

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
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

type TagControllerStruct struct {
	service TagService
}

func NewTagController(service TagService) TagController {
	return &TagControllerStruct{
		service,
	}
}

func (t *TagControllerStruct) GetAll(c *fiber.Ctx) error {

	filter := c.Query("filter")

	items, err := t.service.FindAllWithFilter(filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(items), "items": items})
}
