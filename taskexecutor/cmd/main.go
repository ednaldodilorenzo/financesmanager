package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ednaldodilorenzo/financesmanager/taskexecutor/config"
	"github.com/hibiken/asynq"
	"gopkg.in/gomail.v2"
)

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func sendEmailSMTP(to, subject, body string) error {
	mailSettings, err := config.ReadMailSettings()
	if err != nil {
		fmt.Printf("Failed to read mail settings.")
		return err
	}
	smtpHost := mailSettings.Host // Change to "smtp.office365.com" for Outlook
	smtpPort := mailSettings.Port
	smtpUser := mailSettings.Username
	smtpPass := mailSettings.Password // Use App Passwords for security

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Printf("Email sent to %s\n", to)
	return nil
}

func handleEmailTask(ctx context.Context, task *asynq.Task) error {
	var email EmailPayload
	if err := json.Unmarshal(task.Payload(), &email); err != nil {
		return err
	}

	fmt.Printf("Sending email to %s with subject: %s\n", email.To, email.Subject)

	err := sendEmailSMTP(email.To, email.Subject, email.Body)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func main() {
	mbSettings, err := config.ReadMessageBrokerSettings()
	if err != nil {
		fmt.Printf("Failed to read message broker settings.")
	}
	server := asynq.NewServer(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", mbSettings.Host, mbSettings.Port)}, asynq.Config{
		Concurrency: 10, // Number of concurrent workers
		Queues:      map[string]int{"default": 1},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc("email:send", handleEmailTask)

	fmt.Println("Worker started...")
	if err := server.Run(mux); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}
