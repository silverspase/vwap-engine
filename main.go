package main

import (
	"silverspase/vwap_engine/internal/vwap"

	"go.uber.org/zap"
)

const (
	websocketEndpoint = "wss://ws-feed.exchange.coinbase.com"
	origin            = "http://localhost/"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	service, err := vwap.NewVwapService(websocketEndpoint, "", origin, logger)
	if err != nil {
		logger.Fatal("failed to init service", zap.Error(err))
	}

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
