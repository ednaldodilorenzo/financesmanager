package generic

import (
	"strconv"

	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
)

type GenericController[V model.IUserDependent] interface {
	GetAll(c *fiber.Ctx) error
	GetOne(c *fiber.Ctx) error
	Post(c *fiber.Ctx) error
	PostAll(c *fiber.Ctx) error
	Patch(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type genericController[V model.IUserDependent] struct {
	GenericService[V]
}

func NewGenericController[V model.IUserDependent](service GenericService[V]) GenericController[V] {
	return &genericController[V]{
		service,
	}
}

func (cc *genericController[V]) GetAll(c *fiber.Ctx) error {
	paginate := c.QueryBool("paginate", true)
	pageSize := c.QueryInt("pageSize", 10)

	pageNumber := c.QueryInt("page", 1)
	filter := c.Query("filter")

	loggedUser := c.Locals("user").(model.User)

	if paginate {
		result, err := cc.FindAllPaginatedAndFiltered(c.Context(), int(loggedUser.ID), int(pageSize), int(pageNumber), filter)

		if err != nil {
			return err
		}

		return util.SendData(c, "success", &result, int(fiber.StatusOK))
	} else {
		result, err := cc.FindAll(c.Context(), int(loggedUser.ID))

		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
		}
		return util.SendData(c, "success", &result, int(fiber.StatusOK))
	}

}

func (cc *genericController[V]) GetOne(c *fiber.Ctx) error {

	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	loggedUser := c.Locals("user").(model.User)

	item, err := cc.FindById(c.Context(), itemId, int(loggedUser.ID))

	if err != nil {
		return err
	}

	return util.SendData(c, "success", &item, int(fiber.StatusOK))
}

func (cc *genericController[V]) Delete(c *fiber.Ctx) error {

	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	loggedUser := c.Locals("user").(model.User)

	err = cc.DeleteRecord(c.Context(), itemId, int(loggedUser.ID))

	if err != nil {
		return err
	}

	return util.SendData[any](c, "success", nil, int(fiber.StatusOK))
}

func (cc *genericController[V]) Post(c *fiber.Ctx) error {
	var payload V

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	if errors := util.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	loggedUser, ok := c.Locals("user").(model.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "unauthorized"})
	}

	payload.SetUserID(loggedUser.ID)

	if err := cc.Create(c.Context(), payload); err != nil {
		return err
	}

	return util.SendData[any](c, "success", nil, int(fiber.StatusCreated))
}

func (cc *genericController[V]) PostAll(c *fiber.Ctx) error {
	var payload []V

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	loggedUser := c.Locals("user").(model.User)
	for _, value := range payload {
		if errors := util.ValidateStruct(value); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		} else {
			value.SetUserID(loggedUser.ID)
		}
	}

	if err := cc.CreateAll(c.Context(), payload); err != nil {
		return err
	}

	return util.SendData[any](c, "success", nil, int(fiber.StatusCreated))
}

func (cc *genericController[V]) Patch(c *fiber.Ctx) error {
	var payload V

	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	loggedUser := c.Locals("user").(model.User)

	if err = cc.Update(c.Context(), itemId, payload, int(loggedUser.ID)); err != nil {
		return err
	}

	return util.SendData[any](c, "success", nil, int(fiber.StatusOK))
}
