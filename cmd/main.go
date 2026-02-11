package main

import (
	"crypto-aggregator-service/config"
	httpAPI "crypto-aggregator-service/internal/adapters/httpapi"
)

func main() {
	// Logger
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	// Configs
	configs := config.LoadConfig(logger)

	logger.Infof("Application started, configs: %v", configs)

	//Layout

	layoutLoader := config.NewStaticLayoutLoader()
	layout, err := layoutLoader.Load(nil)

	if err != nil {
		logger.Fatalf("Failed to load layout. %v", err)
	}

	logger.Infof("Loaded layout: %v", layout)

	// Http Router
	httpServer := httpAPI.NewHTTPServer(logger, configs.Server)
	httpServer.Start()

}
