package services

import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"time"
	"sync"
	"strings"
	"reflect"
	"os"
	"github.com/tradesim/model"
)

var btc_endpoint = "http://api.coindesk.com/v1/bpi/currentprice.json"
var stock_endpoint = "https://sandbox.tradier.com/v1/markets/quotes?symbols="
var mutex = &sync.Mutex{}
var CurrentRates map[string]model.Rate

func InitRates() {
	CurrentRates = make(map[string]model.Rate)
	go MaintainUpdatedRates()
}

func MaintainUpdatedRates() {
	for {
		time.Sleep(time.Second * 2)
		symbols := make([]string, 0)
		for symbol, _ := range CurrentRates {
			symbols = append(symbols, symbol)
		}
		fmt.Printf("\nGetting symbols: %s", strings.Join(symbols, ","))
		UpdateRates(symbols)
	}
}

func GetBitcoinRate() []model.Rate {
	var rate model.Rate
	var rates []model.Rate
	rates = make([]model.Rate, 0)

	rate.Price = GetBitcoinPriceUSD()
	rate.Symbol = "BTC"
	rate.Name = "Bitcoin"
	rate.Change = 0.0
	rate.Bid = rate.Price
	rate.Ask = rate.Price

	rates = append(rates, rate)
	return rates
}

func GetBitcoinPriceUSD() float64 {
	response, err := http.Get(btc_endpoint)
	var data map[string]interface{}
	var usd map[string]interface{}
	if err != nil {
		log.Fatal(err)
	} else {
		var objmap map[string]interface{}
		respBody, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(respBody))
		err := json.Unmarshal(respBody, &objmap)
		if err != nil {
			panic(err)
		}
		data = objmap["bpi"].(map[string]interface{})
		usd = data["USD"].(map[string]interface{})
		fmt.Println(usd["rate"])
	}
	retVal := usd["rate_float"].(float64)
	retVal /= 1
	return retVal
}

func GetStockResponse(symbols []string) map[string]interface{} {
	response := MakeStockRequest(symbols)
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		fmt.Printf("\n\nError - Stock Response:\n%s\n\n", respBody)
		panic(err)
	}
	
	return data
}

func MakeStockRequest(symbols []string) *http.Response {
	symbolsString := strings.Join(symbols, ",")
	req, err := http.NewRequest("GET", stock_endpoint + symbolsString, nil)
	req.Header.Add("Authorization", "Bearer " + os.Getenv("STOCK_API_TOKEN"))
	req.Header.Add("Accept", "application/json")
	Wait()
	client := &http.Client{}
	fmt.Printf("\nRequest Url: %s\n", req.URL)
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return response
}

func Wait() {
	mutex.Lock()
	time.Sleep(250 * time.Millisecond)
	mutex.Unlock()
}

func RetrieveRates(symbols []string) []model.Rate {
	var rates []model.Rate
	rates = make([]model.Rate, 0)

	stocks := GetStockResponse(symbols)
	quotes := stocks["quotes"]
	quote := quotes.(map[string]interface{})["quote"]
	if quote == nil {
		return rates
	}
	if reflect.TypeOf(quote).Kind() == reflect.Map {
		tempRate := MapToRate(quote.(map[string]interface{}))
		rates = append(rates, tempRate)
	} else {
		for _, quoteObj := range quote.([]interface{}) {
			tempRate := MapToRate(quoteObj.(map[string]interface{}))
			rates = append(rates, tempRate)
		}
	}

	for _, rate := range rates {
		CurrentRates[rate.Symbol] = rate
	}

	return rates
}

func MapToRate(fromMap map[string]interface{}) model.Rate {
	var toRate model.Rate
	last := fromMap["last"]
	description := fromMap["description"]
	change := fromMap["change"]
	bid := fromMap["bid"]
	ask := fromMap["ask"]
	symbol := fromMap["symbol"]
	toRate.Symbol = symbol.(string)
	toRate.Name = description.(string)
	toRate.Price = last.(float64)
	toRate.Change = change.(float64)
	toRate.Bid = bid.(float64)
	toRate.Ask = ask.(float64)
	return toRate
}

func UpdateRates(symbols []string) {
	if len(symbols) == 0 {
		return
	}
	rates := RetrieveRates(symbols)
	for _, rate := range rates {
		if rate.Symbol == "BTC" {
			CurrentRates[rate.Symbol] = GetBitcoinRate()[0]
		} else if rate.Symbol == "USD" {
			var usd model.Rate
			usd.Name = "United States Dollar"
			usd.Price = 1.0
			usd.Bid = 1.0
			usd.Ask = 1.0
			usd.Change = 0.0
			usd.Symbol = "USD"
			CurrentRates[rate.Symbol] = usd
		} else {
			CurrentRates[rate.Symbol] = rate
		}
	}
}





