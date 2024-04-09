package account

import (
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller generic.GenericController[*model.Account]) (string, func(router fiber.Router)) {
	return "/accounts", func(router fiber.Router) {
		router.Get("/", middleware.DeserializeUser, controller.GetAll)
		router.Get("/:id", middleware.DeserializeUser, controller.GetOne)
		router.Post("/", middleware.DeserializeUser, controller.Post)
		router.Patch("/:id", middleware.DeserializeUser, controller.Patch)
	}
}
