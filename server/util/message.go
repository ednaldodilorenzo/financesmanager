package util

import (
	"encoding/json"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/gofiber/fiber/v2"
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
}

type EmailSenderStruct struct {
	broker *config.Broker
}

func NewEmailSender(broker *config.Broker) EmailSender {
	return &EmailSenderStruct{
		broker: broker,
	}
}

func (e *EmailSenderStruct) SendEmail(to, subject, body string) error {
	email := EmailTask{To: to, Subject: subject, Body: body}
	payload, _ := json.Marshal(email)

	err := e.broker.Enqueue("email:send", payload)

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
