package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	binance "github.com/aiviaio/go-binance/v2"
)

func main() {
	// Create a new Binance client
	client := binance.NewClient("", "")

	// Retrieve the first five trading pairs
	exchangeInfo, err := client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving exchange info: %v", err)
	}

	// Collect the first five symbols
	symbols := make([]string, 0, 5)
	for i, symbol := range exchangeInfo.Symbols {
		if i >= 5 {
			break
		}
		symbols = append(symbols, symbol.Symbol)
	}

	// Create a channel to receive price data
	priceChannel := make(chan map[string]string, len(symbols))

	var wg sync.WaitGroup
	for _, symbol := range symbols {
		wg.Add(1)
		go fetchPrice(symbol, priceChannel, &wg)
	}

	// Close the channel when all goroutines are done
	go func() {
		wg.Wait()
		close(priceChannel)
	}()

	// Print out the results
	for price := range priceChannel {
		for symbol, value := range price {
			fmt.Printf("%s %s\n", symbol, value)
		}
	}
}

func fetchPrice(symbol string, priceChannel chan map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create a new Binance client
	client := binance.NewClient("", "")

	// Retrieve the latest price for the symbol
	price, err := client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		log.Printf("Error retrieving price for %s: %v", symbol, err)
		return
	}

	// Send the price to the channel
	priceData := map[string]string{
		symbol: price[0].Price,
	}
	priceChannel <- priceData
}
