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

type GenericControllerStruct[V model.IUserDependent] struct {
	GenericService[V]
}

func NewGenericController[V model.IUserDependent](service GenericService[V]) GenericController[V] {
	return &GenericControllerStruct[V]{
		service,
	}
}

func (cc *GenericControllerStruct[V]) GetAll(c *fiber.Ctx) error {
	paginate := c.QueryBool("paginate", true)
	pageSize := c.QueryInt("pageSize", 10)

	pageNumber := c.QueryInt("page", 1)
	filter := c.Query("filter")

	loggedUser := c.Locals("user").(model.User)

	var result interface{}
	var err interface{}

	if paginate {
		result, err = cc.FindAllPaginatedAndFiltered(int(loggedUser.ID), int(pageSize), int(pageNumber), filter)

		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
		}

		response := result.(*PaginatedResponse[V])

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "items": response.Items, "total": response.Total, "page": response.Page})
	} else {
		result, err = cc.FindAll(int(loggedUser.ID))

		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "items": result})
	}

}

func (cc *GenericControllerStruct[V]) GetOne(c *fiber.Ctx) error {

	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	loggedUser := c.Locals("user").(model.User)

	item, err := cc.FindById(itemId, int(loggedUser.ID))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "item": item})
}

func (cc *GenericControllerStruct[V]) Delete(c *fiber.Ctx) error {

	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	loggedUser := c.Locals("user").(model.User)

	err = cc.DeleteRecord(itemId, int(loggedUser.ID))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
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

	if err := cc.Create(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (cc *GenericControllerStruct[V]) PostAll(c *fiber.Ctx) error {
	var payload []V

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	loggedUser := c.Locals("user").(model.User)
	for _, value := range payload {
		if errors := util.ValidateStruct(value); errors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(errors)
		} else {
			value.SetUserID(loggedUser.ID)
		}
	}

	if err := cc.CreateAll(payload); err != nil {
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

	loggedUser := c.Locals("user").(model.User)

	if err = cc.Update(itemId, payload, int(loggedUser.ID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
