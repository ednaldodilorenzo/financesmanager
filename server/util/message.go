package util

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
)

type ApiResponse[T any] struct {
	Status   string   `json:"status"`
	Data     T        `json:"data"`
	Messages []string `json:"messages"`
}

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

type EmailTask struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailSender interface {
	SendEmail(to, subject, body string) error
	Config(settings *config.BrokerSettings)
}

type EmailSenderStruct struct {
	settings *config.BrokerSettings
}

func NewEmailSender() EmailSender {
	return &EmailSenderStruct{}
}

func (e *EmailSenderStruct) Config(settings *config.BrokerSettings) {
	e.settings = settings
}

func (e *EmailSenderStruct) SendEmail(to, subject, body string) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", e.settings.Host, e.settings.Port)})
	defer client.Close()

	email := EmailTask{To: to, Subject: subject, Body: body}
	payload, _ := json.Marshal(email)

	task := asynq.NewTask("email:send", payload)
	_, err := client.Enqueue(task, asynq.MaxRetry(5), asynq.Timeout(30*time.Second))

	return err
}

func SendData[T any](ctx *fiber.Ctx, status string, data *T, statusCode int) error {
	var responseData T

	if data != nil {
		responseData = *data
	}

	return ctx.Status(statusCode).JSON(ApiResponse[T]{
		Status: status,
		Data:   responseData,
	})
}
