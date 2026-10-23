package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ednaldodilorenzo/financesmanager/taskexecutor/config"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

const (
	emailStream   = "email:send"
	consumerGroup = "taskexecutor"
	// Pending messages idle for longer than this are claimed again (retry after a failure or a crash).
	retryAfter = time.Minute
)

type EmailPayload struct {
	To      string
	Subject string
	Body    string
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

func handleEmailMessage(msg redis.XMessage) error {
	to, _ := msg.Values["to"].(string)
	subject, _ := msg.Values["subject"].(string)
	body, _ := msg.Values["body"].(string)
	email := EmailPayload{To: to, Subject: subject, Body: body}

	fmt.Printf("Sending email to %s with subject: %s\n", email.To, email.Subject)

	err := sendEmailSMTP(email.To, email.Subject, email.Body)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func processMessages(ctx context.Context, rdb *redis.Client, messages []redis.XMessage) {
	for _, msg := range messages {
		if _, ok := msg.Values["to"].(string); !ok {
			// Malformed message: acknowledge it so it is not retried forever.
			log.Printf("Discarding malformed message %s: %v", msg.ID, msg.Values)
		} else if err := handleEmailMessage(msg); err != nil {
			// Left pending; it will be claimed again after retryAfter.
			continue
		}

		if err := rdb.XAck(ctx, emailStream, consumerGroup, msg.ID).Err(); err != nil {
			log.Printf("Failed to ack message %s: %v", msg.ID, err)
		}
	}
}

func createConsumerGroup(ctx context.Context, rdb *redis.Client) error {
	// "$" means only messages added after the group is created are delivered.
	err := rdb.XGroupCreateMkStream(ctx, emailStream, consumerGroup, "$").Err()
	if err != nil && !strings.HasPrefix(err.Error(), "BUSYGROUP") {
		return err
	}
	return nil
}

func claimStaleMessages(ctx context.Context, rdb *redis.Client, consumer string) []redis.XMessage {
	messages, _, err := rdb.XAutoClaim(ctx, &redis.XAutoClaimArgs{
		Stream:   emailStream,
		Group:    consumerGroup,
		Consumer: consumer,
		MinIdle:  retryAfter,
		Start:    "0",
		Count:    10,
	}).Result()
	if err != nil {
		log.Printf("Failed to claim pending messages: %v", err)
		return nil
	}
	return messages
}

func main() {
	mbSettings, err := config.ReadMessageBrokerSettings()
	if err != nil {
		log.Fatalf("Failed to read message broker settings: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rdb := redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%d", mbSettings.Host, mbSettings.Port)})
	defer rdb.Close()

	if err := createConsumerGroup(ctx, rdb); err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}

	consumer, _ := os.Hostname()

	fmt.Println("Worker started...")
	for ctx.Err() == nil {
		processMessages(ctx, rdb, claimStaleMessages(ctx, rdb, consumer))

		streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumer,
			Streams:  []string{emailStream, ">"},
			Count:    10,
			Block:    5 * time.Second,
		}).Result()
		if errors.Is(err, redis.Nil) || ctx.Err() != nil {
			continue
		}
		if err != nil {
			log.Printf("Failed to read from stream: %v", err)
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range streams {
			processMessages(ctx, rdb, stream.Messages)
		}
	}

	fmt.Println("Worker stopped.")
}
