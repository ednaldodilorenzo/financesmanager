package config

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

type Broker struct {
	Client *asynq.Client
}

func NewBroker() *Broker {
	return &Broker{}
}

func (b *Broker) Connect(settings *BrokerSettings) {
	b.Client = asynq.NewClient(asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%s", settings.Host, settings.Port)})
}

func (b *Broker) Enqueue(taskName string, payload []byte) error {
	task := asynq.NewTask(taskName, payload)
	_, err := b.Client.Enqueue(task, asynq.MaxRetry(5), asynq.Timeout(30*time.Second))

	return err
}

func (b *Broker) Close() error {
	return b.Client.Close()
}
