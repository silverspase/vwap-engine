package vwap

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gosuri/uilive"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

const defaultSlidingWindowLimit = 200

type Service interface {
	PollAndProcessData(subscribeMessage SubscribeMessage)
}

type service struct {
	logger                   *zap.Logger
	ws                       *websocket.Conn
	data                     map[TradingPair]*pairData // all fun is here
	subscribedPairs          []TradingPair             // is used to store subscribed pairs
	slidingWindowLimitByPair map[TradingPair]int       // if you need a custom sliding window, put value here
}

func NewVwapService(endpoint, protocol, origin string, logger *zap.Logger) (Service, error) {
	ws, err := websocket.Dial(endpoint, protocol, origin)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to establish WS connection: %v", err))
	}

	writer := uilive.New()
	writer.Start()

	return &service{
		ws:                       ws,
		logger:                   logger,
		data:                     make(map[TradingPair]*pairData),
		slidingWindowLimitByPair: make(map[TradingPair]int),
	}, nil
}

//PollAndProcessData retrieves data from websocket and calls processing function
func (s *service) PollAndProcessData(subscribeMessage SubscribeMessage) {
	s.subscribedPairs = subscribeMessage.ProductIds

	bytes, err := json.Marshal(subscribeMessage)
	if err != nil {
		s.logger.Fatal("failed to marshal json, check the param message", zap.Error(err))
	}

	if _, err := s.ws.Write(bytes); err != nil {
		s.logger.Fatal("failed to send message to ws", zap.Error(err))
	}

	var msg = make([]byte, 512)
	var n int
	for {
		if n, err = s.ws.Read(msg); err != nil {
			s.logger.Error("failed to read ws response", zap.Error(err))
			continue
		}
		s.processData(msg[:n])
		s.printData()
	}
}

//processData unmarshalls input and calls vwap update
func (s *service) processData(rawData []byte) {
	var singleTrade MatchResponse
	err := json.Unmarshal(rawData, &singleTrade)
	if err != nil {
		s.logger.Error("json parsing failed", zap.Error(err))
		return
	}

	//if we don't have given key - create it and prefill with value
	_, ok := s.data[singleTrade.ProductID]
	if !ok {
		validPair := false
		for _, v := range s.subscribedPairs {
			if v == singleTrade.ProductID { //continue only with pairs defined in subscribe message
				slidingLimit, ok := s.slidingWindowLimitByPair[singleTrade.ProductID]
				if !ok { // if there is no predefined window limit, pick up the default one
					slidingLimit = defaultSlidingWindowLimit
				}
				s.data[singleTrade.ProductID] = &pairData{
					logger:             s.logger,
					slidingWindowLimit: slidingLimit,
				}
				validPair = true
			}
		}
		// in raw data we are getting untyped messages with not valid pair, we should skip them
		if !validPair {
			return
		}
	}

	s.data[singleTrade.ProductID].updateVWAP(singleTrade.Price, singleTrade.Size)
}

//printData outputs Volume Weighted Adjusted Price into console
func (s *service) printData() {
	fmt.Println(strings.Repeat("# ", 15))
	for _, pair := range s.subscribedPairs {
		if s.data[pair] != nil {
			printData := ""
			if len(s.data[pair].slidingWindowData) < s.data[pair].slidingWindowLimit {
				printData += fmt.Sprintf("collecting data [%v%%] ", len(s.data[pair].slidingWindowData)/(s.data[pair].slidingWindowLimit/100))
			}
			printData += fmt.Sprintf("%s: vwap %.2f; cumulativeVolume %.2f %s \n",
				pair, s.data[pair].volumeWeightedAveragePrice, s.data[pair].cumulativeData.volume, strings.Split(string(pair), "-")[0])

			fmt.Print(printData)
		}
	}
	time.Sleep(time.Millisecond * 500) // increases readability
}
