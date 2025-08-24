package pool

import (
	"context"
	"outbox_service/global"
	"outbox_service/internal/helper"

	"github.com/thanvuc/go-core-lib/eventbus"
)

func Run(ctx context.Context) {
	publisher := eventbus.NewPublisher(
		global.EventBusConnector,
		eventbus.SyncDatabaseExchange,
		eventbus.ExchangeTypeTopic,
		helper.IntPtr(3),
		helper.IntPtr(2000),
		true,
	)

	authPool := NewAuthPool(publisher)

	go authPool.OpenUserPoolWorker(ctx)
}
