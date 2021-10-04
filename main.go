package main

import (
	"log"
	"silverspase/vwap_engine/internal/vwap"

	"go.uber.org/zap"
)

const (
	websocketEndpoint = "wss://ws-feed.exchange.coinbase.com"
	origin            = "http://localhost/"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) { // flushes buffer, if any
		err := logger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(logger)

	service, err := vwap.NewVwapService(websocketEndpoint, "", origin, logger)
	if err != nil {
		logger.Fatal("failed to init service", zap.Error(err))
	}

	// here you can define trading pairs to process
	service.PollAndProcessData(vwap.SubscribeMessage{
		Type: "subscribe",
		ProductIds: []vwap.TradingPair{
			vwap.BTC_USD,
			vwap.ETH_USD,
			vwap.ETH_BTC,
		},
		Channels: []string{"matches"},
	})
}
