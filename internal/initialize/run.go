package initialize

import (
	"context"
	"os"
	"os/signal"
	"outbox_service/global"
	"outbox_service/internal/eventbus"
	"outbox_service/internal/pool"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

func Run() {
	if err := initConfigAndResources(); err != nil {
		global.Logger.Error("Failed to initialize configs and resources", "", zap.Error(err))
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	global.Logger.Info("Service is starting...", "")

	eventbus.Run(ctx)
	pool.Run(ctx)

	<-stop
	global.Logger.Info("Shutdown signal received, shutting down...", "")

	cancel()

	gracefulShutdown(wg, global.Logger)
}
