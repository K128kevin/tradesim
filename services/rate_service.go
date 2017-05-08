package services

import (
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

var btc_endpoint = "http://api.coindesk.com/v1/bpi/currentprice.json"

func GetBitcoinPriceUSD(symbol string) float64 {
	response, err := http.Get(btc_endpoint)
	var data map[string]interface{}
	var usd map[string]interface{}
	if err != nil {
		log.Fatal(err)
	} else {
		var objmap map[string]interface{}
		respBody, _ := ioutil.ReadAll(response.Body)
		err := json.Unmarshal(respBody, &objmap)
		if err != nil {
			panic(err)
		}
		data = objmap["bpi"].(map[string]interface{})
		usd = data["USD"].(map[string]interface{})
		fmt.Println(usd["rate"])
	}
	retVal := usd["rate_float"].(float64)
	if (symbol == "USD") {
		retVal /= 1
	}
	return retVal
}

func GetStockPriceUSD(symbol string) (float64, error) {
	url := "http://dev.markitondemand.com/MODApis/Api/v2/Quote/json?symbol=" + symbol
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respBody))
	var data map[string]interface{}
	err = json.Unmarshal(respBody, &data)
	if err != nil {
		panic(err)
	}
	price := data["LastPrice"]
	if price == nil || price.(float64) == 0.0 {
		return 0.0, fmt.Errorf("Provided symbol (" + symbol + ") is not available")
	}
	fmt.Printf("\n%d", price.(float64))
	return price.(float64), nil
}