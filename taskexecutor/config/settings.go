package config

import (
	"os"
	"strconv"
)

type MailSettings struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MessageBrokerSettings struct {
	Host string
	Port int
}

var settings *MailSettings

var brokerSettings *MessageBrokerSettings

func ReadMailSettings() (*MailSettings, error) {
	mailPort, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	if err != nil {
		return nil, err
	}

	if settings == nil {
		settings = &MailSettings{
			Host:     os.Getenv("EMAIL_HOST"),
			Port:     mailPort,
			Username: os.Getenv("EMAIL_USERNAME"),
			Password: os.Getenv("EMAIL_PASSWORD"),
		}
	}

	return settings, nil
}

func ReadMessageBrokerSettings() (*MessageBrokerSettings, error) {
	brokerPort, err := strconv.Atoi(os.Getenv("MB_PORT"))

	if err != nil {
		return nil, err
	}

	if brokerSettings == nil {
		brokerSettings = &MessageBrokerSettings{
			Host: os.Getenv("MB_HOST"),
			Port: brokerPort,
		}
	}

	return brokerSettings, nil
}
