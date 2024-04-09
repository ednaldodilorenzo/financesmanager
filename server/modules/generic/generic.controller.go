package generic

import (
	"fmt"
	"strconv"

	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
)

type GenericController[V model.IUserDependent] interface {
	GetAll(c *fiber.Ctx) error
	GetOne(c *fiber.Ctx) error
	Post(c *fiber.Ctx) error
	Patch(c *fiber.Ctx) error
}

type GenericControllerStruct[V model.IUserDependent] struct {
	GenericService[V]
}

func NewGenericController[V model.IUserDependent](service GenericService[V]) GenericController[V] {
	return &GenericControllerStruct[V]{
		service,
	}
}

func (cc *GenericControllerStruct[V]) GetAll(c *fiber.Ctx) error {
	items, err := cc.FindAll()

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(items), "items": items})
}

func (cc *GenericControllerStruct[V]) GetOne(c *fiber.Ctx) error {

	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	item, err := cc.FindById(itemId)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "item": item})
}

func (cc *GenericControllerStruct[V]) Post(c *fiber.Ctx) error {
	var payload V

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if errors := util.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	loggedUser := c.Locals("user").(model.User)
	payload.SetUserID(loggedUser.ID)

	fmt.Print(loggedUser)

	if err := cc.Create(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (cc *GenericControllerStruct[V]) Patch(c *fiber.Ctx) error {
	var payload *V

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	if err = cc.Update(itemId, payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
