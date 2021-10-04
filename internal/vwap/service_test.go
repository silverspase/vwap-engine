package vwap

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func Test_updateVWAP(t *testing.T) {
	s := service{
		logger:                   zap.NewNop(),
		data:                     make(map[TradingPair]*pairData),
		subscribedPairs:          []TradingPair{BTC_USD, ETH_USD, ETH_BTC},
		slidingWindowLimitByPair: map[TradingPair]int{BTC_USD: 2, ETH_USD: 2},
	}

	for _, entry := range getDataFromJson() {
		bytes, err := json.Marshal(entry)
		if err != nil {
			log.Fatal(err)
		}
		s.processData(bytes)
	}

	assert.Equal(t, 48794.72554883774, s.data[BTC_USD].volumeWeightedAveragePrice)
	assert.Equal(t, 0.01134789000000005, s.data[BTC_USD].cumulativeData.volume)
	assert.Equal(t, 553.7171781084028, s.data[BTC_USD].cumulativeData.order)
	// we have 3 entries for BTC_USD pair, sliding window is two, so it should contain two latest entries
	assert.Equal(t, []cumulativeData{
		{order: 490.42508498399997, volume: 0.01005065},
		{order: 63.2920931244, volume: 0.00129724},
	}, s.data[BTC_USD].slidingWindowData)

	assert.Equal(t, 3453.031432179072, s.data[ETH_USD].volumeWeightedAveragePrice)
	assert.Equal(t, 0.056942599999999996, s.data[ETH_USD].cumulativeData.volume)
	assert.Equal(t, 196.62458763, s.data[ETH_USD].cumulativeData.order)
	assert.Equal(t, []cumulativeData{
		{order: 100.07458091000001, volume: 0.028981},
		{order: 96.55000672, volume: 0.0279616},
	}, s.data[ETH_USD].slidingWindowData)
	//
	assert.Equal(t, 0.07078, s.data[ETH_BTC].volumeWeightedAveragePrice)
	assert.Equal(t, 0.00288811, s.data[ETH_BTC].cumulativeData.volume)
	assert.Equal(t, s.data[ETH_BTC].cumulativeData.order, 0.0002044204258)
	assert.Equal(t, []cumulativeData{{order: 0.0002044204258, volume: 0.00288811}}, s.data[ETH_BTC].slidingWindowData)

	// we didn't set it during service init, should be default value then
	assert.Equal(t, 200, s.data[ETH_BTC].slidingWindowLimit)
}

func getDataFromJson() []interface{} {
	// Open our jsonFile
	jsonFile, err := os.Open("testData.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result []interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
