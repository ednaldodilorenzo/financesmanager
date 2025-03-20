package transaction

import (
	"strconv"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
)

func GetRoutes(controller TransactionController, deserializer *middleware.Deserializer) (string, func(router fiber.Router)) {
	return "/transactions", func(router fiber.Router) {
		router.Get("/", deserializer.DeserializeUser, controller.GetAllWithRelationships)
		router.Get("/:id", deserializer.DeserializeUser, controller.GetOne)
		router.Post("/", deserializer.DeserializeUser, controller.Post)
		router.Post("/list", deserializer.DeserializeUser, controller.PostAll)
		router.Post("/upload", deserializer.DeserializeUser, controller.UploadBatchFile)
		router.Patch("/:id", deserializer.DeserializeUser, controller.Patch)
		router.Delete("/:id", deserializer.DeserializeUser, controller.Delete)
	}
}

type TransactionUploadSchema struct {
	CategoryId      *uint32   `json:"categoryId"`
	AccountId       uint32    `json:"accountId" validate:"required"`
	Description     string    `json:"description" validate:"required"`
	Value           int32     `json:"value" validate:"required"`
	PaymentDate     time.Time `json:"paymentDate" validate:"required"`
	TransactionDate time.Time `json:"transactionDate" validate:"required"`
	Duplicated      bool      `json:"duplicated" validate:"required"`
	Detail          *string   `json:"detail"`
}

type TransactionController interface {
	generic.GenericController[*model.Transaction]
	GetAllWithRelationships(c *fiber.Ctx) error
	UploadBatchFile(c *fiber.Ctx) error
}

type TransactionControllerStruct struct {
	generic.GenericController[*model.Transaction]
	service TransactionService
}

func NewTransactionController(service TransactionService, controller generic.GenericController[*model.Transaction]) TransactionController {
	return &TransactionControllerStruct{
		controller,
		service,
	}
}

func (cc *TransactionControllerStruct) GetOne(c *fiber.Ctx) error {
	itemId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	loggedUser := c.Locals("user").(model.User)

	item, err := cc.service.FindById(itemId, int(loggedUser.ID))

	if err != nil {
		return err
	}

	return util.SendData(c, "success", item, int(fiber.StatusOK))
}

func (cc *TransactionControllerStruct) GetAllWithRelationships(c *fiber.Ctx) error {

	var monthParam, yearParam *int = nil, nil

	if month := c.Query("month"); month != "" {
		if value, err := strconv.Atoi(month); err != nil {
			return err
		} else {
			monthParam = &value
		}
	}

	if year := c.Query("year"); year != "" {
		if value, err := strconv.Atoi(year); err != nil {
			return err
		} else {
			yearParam = &value
		}
	}

	loggedUser := c.Locals("user").(model.User)

	items, err := cc.service.FindAllRelated(monthParam, yearParam, int(loggedUser.ID))

	if err != nil {
		return err
	}

	return util.SendData(c, "success", &items, int(fiber.StatusOK))
}

func (cc *TransactionControllerStruct) UploadBatchFile(c *fiber.Ctx) error {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	// Open the file
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	accountIDStr := c.FormValue("accountId")
	var accountId int
	if accountIDStr != "" {
		accountId, err = strconv.Atoi(accountIDStr)
		if err != nil {
			return err
		}
	}

	var date time.Time
	dateStr := c.FormValue("paymentDate")
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return err
		}
	}

	fileType := c.FormValue("fileType")
	loggedUser := c.Locals("user").(model.User)

	transactions, err := cc.service.PrepareFileImport(fileReader, uint32(accountId), &date, fileType, int(loggedUser.ID))
	if err != nil {
		return err
	}

	return util.SendData(c, "success", &transactions, int(fiber.StatusOK))
}
