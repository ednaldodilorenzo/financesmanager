package category

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller generic.GenericController[*model.Category], desializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/categories", func(router fiber.Router) {
		router.Get("/", desializer.DeserializeUser, controller.GetAll)
		router.Get("/:id", desializer.DeserializeUser, controller.GetOne)
		router.Post("/", desializer.DeserializeUser, controller.Post)
		router.Patch("/:id", desializer.DeserializeUser, controller.Patch)
	}
}
