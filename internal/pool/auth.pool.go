package pool

import (
	"context"
	"outbox_service/global"
	"outbox_service/internal/eventbus/consumer"
	"outbox_service/internal/helper"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanvuc/go-core-lib/eventbus"
	"github.com/thanvuc/go-core-lib/log"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type AuthPool struct {
	AuthPostgresPool *pgxpool.Pool
	logger           log.Logger
	publisher        eventbus.Publisher
	queryHelper      *helper.QueryHelper
}

func NewAuthPool(
	publisher eventbus.Publisher,
) *AuthPool {
	return &AuthPool{
		logger:      global.Logger,
		publisher:   publisher,
		queryHelper: helper.NewQueryHelper(global.AuthPostgresPool),
	}
}

func (p *AuthPool) OpenUserPoolWorker(ctx context.Context) error {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Context cancelled, stopping worker", "")
			return nil
		case <-ticker.C:
			if err := p.handleOutboxes(ctx); err != nil {
				p.logger.Error("Error handling outboxes", err.Error())
				return err
			}
		}
	}
}

func (p *AuthPool) handleOutboxes(ctx context.Context) error {
	outboxes, err := p.queryHelper.SelectOutboxes(ctx)
	if err != nil {
		p.logger.Error("Failed to select outboxes", err.Error())
		return err
	}

	// Process the outboxes and publish events
	for _, outbox := range outboxes {
		if outbox == nil {
			p.logger.Warn("Received nil outbox, skipping", "nil_outbox")
			continue
		}

		outboxBytes, err := proto.Marshal(outbox)
		if err != nil {
			p.logger.Error("Failed to marshal outbox", outbox.RequestId)
			continue
		}

		err = p.publisher.SafetyPublish(
			ctx,
			outbox.RequestId,
			[]string{"sync.auth.user"},
			outboxBytes,
			nil,
			helper.ExchangeNamePtr(eventbus.DLQSyncDatabaseExchange),
			helper.StringPtr(consumer.SyncAuthDLQ_RoutingKey),
		)

		if err != nil {
			p.logger.Error("Failed to publish outbox event", outbox.RequestId, zap.Error(err))
			continue
		}

		p.queryHelper.SetProcessedOutbox(ctx, outbox)
		p.logger.Info("Successfully published outbox event", outbox.RequestId)
	}

	return nil
}
