package main

import (
	"context"
	"crypto-aggregator-service/config"
	"crypto-aggregator-service/internal/adapters"
	httpAPI "crypto-aggregator-service/internal/adapters/httpapi"
	"crypto-aggregator-service/internal/adapters/webclients"
	"crypto-aggregator-service/internal/repositories"
	"crypto-aggregator-service/internal/services"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
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

	cleanLayout := configs.App.ToDomain()
	layoutStore := repositories.NewLayoutStore(cleanLayout)
	vendorsMap := configs.App.GetVendorMap()

	// Providers

	httpClient := webclients.NewClient(3 * time.Second)

	/*providers := []services.QuoteProvider{
		repositories.NewBitsoCryptoProvider(httpClient),
	}*/

	clients := map[string]repositories.CryptoClient{
		"bitso": repositories.NewBitsoCryptoProvider(httpClient),
		"mock":  &adapters.MockClient{},
	}
	// AggregatorSVC

	//aggSvc := services.NewAggSvc(layoutLoader, providers, 2*time.Second)

	// Poller
	poller := services.NewPoller(layoutStore, clients, vendorsMap, logger)
	ctx, cancel := context.WithCancel(context.Background())

	// Start polling loop in a goroutine
	go poller.Start(ctx, time.Duration(configs.Server.RefreshInterval))

	// HttpServer
	httpServer := httpAPI.NewHTTPServer(logger, configs.Server)

	httpAPI.NewHealthController(httpServer)

	httpAPI.NewAggregatorController(httpServer, nil, poller)

	//httpServer.Start()

	// Graceful Shutdown Channel
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Server Started", zap.String("port", strconv.Itoa(configs.Server.Port)))
		httpServer.Start()
	}()

	// Wait for signal
	<-done
	logger.Info("Server Stopped")

	// Cleanup
	cancel() // Stop the poller

	_, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	logger.Info("Server Exited Properly")

}
