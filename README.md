# Volume-Weighted Average Price Engine

Service which calculates Volume-Weighted Average Price for predefined trading pairs in real time

## Description

This is golang microservice which does next:
1. Retrieves data from coinbase api for predefined trading pairs(BTC-USD, ETH-USD, ETH-BTC).
2. Calculate Volume-Weighted Average Price within sliding window of last 200(or any other defined number)
3. Prints to console Volume-Weighted Average Price for each trading pair.

## Prerequisites
- Golang

## How to use

1. Clone the repo
2. Install dependencies with `go mod tidy`
3. Run in dev mode `go run main.go` or build it and execute as a binary `go build main.go && ./main`

you should see the following:
```
collecting data [8%] BTC-USD: vwap 47644.21; cumulativeVolume 0.70 BTC 
collecting data [15%] ETH-USD: vwap 3348.45; cumulativeVolume 5.81 ETH 
collecting data [0%] ETH-BTC: vwap 0.07; cumulativeVolume 0.05 ETH
```

## Features

- Picking up any trading pairs supported by coinbase api
- Defining sliding window size by each trading pair separately
- Interactive output, where user can see the sliding window's data load in percentage.
- Key-Value structured logging.

## Assumptions

- I decided to implement solution without goroutines considering RPC rate and for the simplicity sake.
If needed, I can implement the same solution using separate goroutine for each trading pair and channels for communication. 
- I covered with UT only part with calculation(the most sufficient part)
- I decided to stick with simple folder structure, but for complex projects I prefer to use `Clean Architecture` approach with additional layers(validator, normalizer etc)
- I could use some tool to get command line arguments or Makefile for storing commands, but decided to focus on the main task
- To make this service fault resistant, we can run it as systemd process or wrap it up in docker and run as a k8s pod.
