package account

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller generic.GenericController[*model.Account], desirializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/accounts", func(router fiber.Router) {
		router.Get("/", desirializer.DeserializeUser, controller.GetAll)
		router.Get("/:id", desirializer.DeserializeUser, controller.GetOne)
		router.Post("/", desirializer.DeserializeUser, controller.Post)
		router.Patch("/:id", desirializer.DeserializeUser, controller.Patch)
	}
}
