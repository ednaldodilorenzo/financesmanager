package category

import (
	"strconv"

	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes() (string, func(router fiber.Router)) {
	categoryController := NewCategoryController()
	return "/categories", func(router fiber.Router) {
		router.Get("/", middleware.DeserializeUser, categoryController.GetAll)
		router.Get("/:id", middleware.DeserializeUser, categoryController.GetOne)
		router.Post("/", middleware.DeserializeUser, categoryController.Post)
		router.Patch("/:id", middleware.DeserializeUser, categoryController.Patch)
	}
}

type CategorySchema struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
}

type CategoryUpdateSchema struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type CategoryController struct {
	generic.GenericService[model.Category]
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		generic.NewGenericService[model.Category](),
	}
}

func (cc *CategoryController) GetAll(c *fiber.Ctx) error {
	categories, err := cc.FindAll()

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(categories), "items": categories})
}

func (cc *CategoryController) GetOne(c *fiber.Ctx) error {

	categoryId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	category, err := cc.FindById(categoryId)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "item": category})
}

func (cc *CategoryController) Post(c *fiber.Ctx) error {
	var payload *CategorySchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if errors := util.ValidateStruct(payload); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	loggedUser := c.Locals("user").(model.User)

	newCategory := model.Category{
		Name:   payload.Name,
		Type:   payload.Type,
		IdUser: loggedUser.ID,
	}

	if err := cc.Create(&newCategory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success"})
}

func (cc *CategoryController) Patch(c *fiber.Ctx) error {
	var payload *CategoryUpdateSchema

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	updates := make(map[string]interface{})

	if payload.Name != "" {
		updates["name"] = payload.Name
	}

	if payload.Type != "" {
		updates["type"] = payload.Type
	}

	categoryId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	if err = cc.Update(categoryId, updates); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "fail"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
