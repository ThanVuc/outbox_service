package initialize

import (
	"outbox_service/global"
	"sync"

	"github.com/thanvuc/go-core-lib/log"
)

func initConfigAndResources() error {
	loadConfig()
	initLogger()
	initPostgreSQL()
	initEventBus()

	return nil
}

func gracefulShutdown(wg *sync.WaitGroup, logger log.Logger) {

	if global.AuthPostgresPool != nil {
		global.AuthPostgresPool.Close()
		logger.Info("PostgreSQL connection pool closed", "")
	}

	if global.EventBusConnector != nil {
		wg.Add(1)
		global.EventBusConnector.Close(wg)
		logger.Info("EventBus connection closed", "")
	}
	wg.Wait()
	global.Logger.Info("Application shutdown gracefully", "")
}
