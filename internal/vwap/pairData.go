package vwap

import (
	"go.uber.org/zap"
)

type pairData struct {
	logger                     *zap.Logger
	volumeWeightedAveragePrice float64
	slidingWindowLimit         int
	slidingWindowData          []cumulativeData
	cumulativeData
}

type cumulativeData struct {
	order  float64 // is cumulative sum of price * size within sliding window
	volume float64 // total volume within sliding window
}

func (p *pairData) updateVWAP(price, size float64) {
	p.updateTradingPairData(price, size)
	if p.cumulativeData.volume == 0 { // not likely, but let's prevent zero division error
		p.logger.Warn("cumulativeVolume is zero")
		return
	}

	p.volumeWeightedAveragePrice = p.cumulativeData.order / p.cumulativeData.volume
}

func (p *pairData) updateTradingPairData(price, size float64) {
	// if we reach the maximum limit in slidingWindowData:
	// 1. Remove the first entry from slidingWindowData
	// 2. Subtract removed data from cumulative order and volume
	if len(p.slidingWindowData) > p.slidingWindowLimit-1 {
		p.cumulativeData.order -= p.slidingWindowData[0].order
		p.cumulativeData.volume -= p.slidingWindowData[0].volume
		p.slidingWindowData = p.slidingWindowData[1:]
	}

	p.cumulativeData.order += price * size
	p.cumulativeData.volume += size
	p.slidingWindowData = append(p.slidingWindowData, cumulativeData{
		order:  price * size,
		volume: size,
	})
}
