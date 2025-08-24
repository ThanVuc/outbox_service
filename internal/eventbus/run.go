package eventbus

import (
	"context"
	"outbox_service/internal/eventbus/consumer"
)

func Run(ctx context.Context) {
	authDLQ := consumer.NewAuthDLQConsumer()
	go authDLQ.ConsumeSyncDLQ(ctx)
}
