package consumer

import (
	"context"
	"outbox_service/global"
	"outbox_service/internal/helper"
	"outbox_service/proto/common"

	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type AuthDLQConsumer struct {
	logger      log.Logger
	queryHelper *helper.QueryHelper
}

func NewAuthDLQConsumer() *AuthDLQConsumer {
	return &AuthDLQConsumer{
		logger:      global.Logger,
		queryHelper: helper.NewQueryHelper(global.AuthPostgresPool),
	}
}

func (adq *AuthDLQConsumer) ConsumeSyncDLQ(ctx context.Context) {
	logger := global.Logger
	syncConsumerDLQ := eventbus.NewConsumer(
		global.EventBusConnector,
		eventbus.DLQSyncDatabaseExchange,
		eventbus.ExchangeTypeTopic,
		SyncAuthDLQ_RoutingKey,
		"sync_dlq",
		2,
	)

	err := syncConsumerDLQ.Consume(ctx, func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		requestId := d.Headers["request_id"].(string)
		logger.Warn("Received message from sync DLQ", requestId)
		outbox := &common.Outbox{}
		err := proto.Unmarshal(d.Body, outbox)
		if err != nil {
			logger.Error("Failed to unmarshal message from sync user DB queue", "", zap.Error(err))
			return rabbitmq.NackDiscard
		}

		if outbox.ProcessedAt != nil || outbox.Status != common.OutboxStatus_OUTBOX_STATUS_PENDING.String() {
			print("Skipping outbox with processed_at or non-pending status")
			adq.queryHelper.SetFailedOutbox(ctx, outbox)
		}

		return rabbitmq.Ack
	})

	if err != nil {
		logger.Error("Failed to consume messages from sync DLQ", "")
		return
	}
}
