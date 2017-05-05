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