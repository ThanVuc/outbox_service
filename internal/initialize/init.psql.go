package initialize

import (
	"context"
	"fmt"
	"outbox_service/global"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thanvuc/go-core-lib/log"
	"go.uber.org/zap"
)

func initPostgreSQL() {
	logger := global.Logger
	initPostgreSQLAuth(logger)
}

func initPostgreSQLAuth(logger log.Logger) {
	dsn := "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai"
	configs := global.Config.PostgresAuth
	var connectString = fmt.Sprintf(dsn, configs.Host, configs.User, configs.Password, configs.Database, configs.Port)
	ctx := context.Background()
	for {
		config, err := pgxpool.ParseConfig(connectString)
		if err != nil {
			logger.Error("Failed to parse PostgreSQL connection string", "", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}
		setPostgresConfig(config)
		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			logger.Error("Failed to create PostgreSQL connection pool", "", zap.Error(err))
			time.Sleep(5 * time.Second)
			continue
		}

		if err := pool.Ping(ctx); err != nil {
			logger.Error("Failed to ping PostgreSQL", "", zap.Error(err))
			pool.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		global.AuthPostgresPool = pool
		break
	}
	logger.Info("Auth PostgreSQL connection pool initialized successfully", "", zap.String("host", configs.Host), zap.Int("port", configs.Port))
}

func setPostgresConfig(config *pgxpool.Config) {
	postConfig := global.Config.PostgresAuth
	config.MaxConns = int32(postConfig.MaxOpenConns)
	config.MaxConnIdleTime = time.Duration(postConfig.ConnMaxIdleTime) * time.Second
	config.MaxConnLifetime = time.Duration(postConfig.MaxLifetime) * time.Second
}
