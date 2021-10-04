package vwap

type TradingPair string

const (
	BTC_USD TradingPair = "BTC-USD"
	ETH_USD TradingPair = "ETH-USD"
	ETH_BTC TradingPair = "ETH-BTC"
)

//SubscribeMessage is used as a param to subscribe to websocket api
type SubscribeMessage struct {
	Type       string        `json:"type"`
	ProductIds []TradingPair `json:"product_ids"`
	Channels   []string      `json:"channels"`
}

//MatchResponse is the response type from "channel" websocket
type MatchResponse struct {
	Size      float64     `json:"size,string"`
	Price     float64     `json:"price,string"`
	ProductID TradingPair `json:"product_id"`
}
