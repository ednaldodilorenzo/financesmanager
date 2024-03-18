package routes

import (
	"github.com/ednaldo-dilorenzo/iappointment/modules/account"
	"github.com/ednaldo-dilorenzo/iappointment/modules/auth"
	"github.com/ednaldo-dilorenzo/iappointment/modules/category"
	"github.com/ednaldo-dilorenzo/iappointment/modules/transaction"
	"github.com/gofiber/fiber/v2"
)

func SetRoutes(api *fiber.Router) {
	(*api).Route(auth.GetRoutes())
	(*api).Route(category.GetRoutes())
	(*api).Route(account.GetRoutes())
	(*api).Route(transaction.GetRoutes())
}
