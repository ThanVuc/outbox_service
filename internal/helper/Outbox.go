package helper

import (
	"context"
	"fmt"
	"outbox_service/global"
	"outbox_service/proto/common"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/thanvuc/go-core-lib/log"
	"go.uber.org/zap"
)

const openPoolQuery = `
	SELECT
		id,
		aggregate_type,
		aggregate_id,
		event_type,
		payload,
		status,
		occurred_at,
		processed_at,
		error_message,
		retry_count,
		request_id
	FROM outbox
	WHERE status = 1
	ORDER BY occurred_at DESC
`

const successedOutboxQuery = `
	UPDATE outbox
	SET status = 2, processed_at = NOW()
	WHERE id = $1
`

const failedOutboxQuery = `
	UPDATE outbox
	SET status = 3, processed_at = NOW(), error_message = $2, retry_count = retry_count + 1
	WHERE id = $1
`

type QueryHelper struct {
	psqlPool *pgxpool.Pool
	logger   log.Logger
}

func NewQueryHelper(psqlPool *pgxpool.Pool) *QueryHelper {
	return &QueryHelper{
		psqlPool: psqlPool,
		logger:   global.Logger,
	}
}

func (q *QueryHelper) SelectOutboxes(ctx context.Context) ([]*common.Outbox, error) {
	rows, err := q.psqlPool.Query(ctx, openPoolQuery)
	if err != nil {
		q.logger.Error("Failed to query outbox table", "")
		return nil, err
	}

	outboxes := make([]*common.Outbox, 0)
	for rows.Next() {
		var (
			id            string
			aggregateType string
			aggregateId   string
			eventType     string
			payload       []byte
			status        int16
			occurredAt    time.Time
			processedAt   *time.Time
			errorMessage  *string
			retryCount    int32
			requestId     string
		)

		err = rows.Scan(
			&id,
			&aggregateType,
			&aggregateId,
			&eventType,
			&payload,
			&status,
			&occurredAt,
			&processedAt,
			&errorMessage,
			&retryCount,
			&requestId,
		)
		if err != nil {
			q.logger.Error("Failed to scan outbox row", "", zap.Error(err))
			return nil, err
		}

		out := &common.Outbox{
			Id:            id,
			AggregateType: aggregateType,
			AggregateId:   aggregateId,
			EventType:     eventType,
			Payload:       payload,
			Status:        fmt.Sprint(status),
			OccurredAt:    occurredAt.Unix(),
			RetryCount:    retryCount,
			RequestId:     requestId,
		}

		if processedAt != nil {
			t := processedAt.Unix()
			out.ProcessedAt = &t
		}
		if errorMessage != nil {
			out.ErrorMessage = errorMessage
		}

		outboxes = append(outboxes, out)
	}
	return outboxes, nil
}

func (q *QueryHelper) SetProcessedOutbox(ctx context.Context, outbox *common.Outbox) error {
	_, err := q.psqlPool.Exec(ctx, successedOutboxQuery, outbox.Id)
	if err != nil {
		q.logger.Error("Failed to update outbox status", outbox.RequestId, zap.Error(err))
		return err
	}
	return nil
}

func (q *QueryHelper) SetFailedOutbox(ctx context.Context, outbox *common.Outbox) error {
	_, err := q.psqlPool.Exec(ctx, failedOutboxQuery, outbox.Id, outbox.ErrorMessage)
	if err != nil {
		q.logger.Error("Failed to update outbox status", outbox.RequestId, zap.Error(err))
		return err
	}
	return nil
}
