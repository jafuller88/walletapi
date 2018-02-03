package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"walletapi/config"
	"walletapi/log"
	"walletapi/rpc"
)

const (
	version         = "1.0.0"
	apiVVV          = "%vfsym=%v&tsyms=%v"
	htmlTemplateVVV = "<font color=blue>%v</font>: <font color=yellow>%v%v</font></br>"
)

const (
	jsonResponse = iota
	htmlResponse
)

// CoinPrice struct for response from api call
type CoinPrice struct {
	BTC float64 `json:"BTC"`
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
}

// JSONResponse struct
type JSONResponse struct {
	Balance float64 `json:"balance"`
	Price   float64 `json:"price"`
	Total   float64 `json:"totalValue"`
}

func main() {
	log.Msgf(0, "Crypto Wallet API Version %v\n", version)

	// Setup the flags and config
	rpcServer := flag.String("rpc", "http://127.0.0.1:3000/json_rpc", "RPC URL and Port")
	apiPort := flag.String("port", "8080", "Web API Server Port")
	apiDomain := flag.String("api", "https://min-api.cryptocompare.com/data/price?", "Price API URL")
	fsym := flag.String("fsym", "ETN", "From Symbol")
	tsym := flag.String("tsym", "BTC,USD,EUR", "To Symbol/s")
	responseType := flag.Int("rt", 0, "API response type (0=JSON,1=HTML)")
	logLevel := flag.Int("logLevel", 0, "Log Level (0=FATAL,1=ERROR,2=INFO,3=DEBUG")

	flag.Parse()

	// Setup the running configuration from the set flags..
	config.RPCServerURL = *rpcServer
	config.APIDomain = *apiDomain
	config.APIPort = *apiPort
	config.FromSym = *fsym
	config.ToSym = *tsym
	config.ResponseType = *responseType
	config.APIURL = fmt.Sprintf(apiVVV, config.APIDomain, config.FromSym, config.ToSym)
	config.LogLevel = *logLevel

	// test connection to rpc by getting block height
	log.Msgf(0, "Connecting to RPC Server...\n")
	h, e := rpc.GetHeight()
	if e != nil {
		log.Msgf(0, "The RPC Service is not Running\n")
	} else {
		log.Msgf(0, "Block Height: %v\n", h)
		log.Msgf(0, "Listening for requests on port %v\n", config.APIPort)
		// setup http endpoint
		http.HandleFunc("/balance", getBalanceHTTPRequestHandler)           // set router
		err := http.ListenAndServe(fmt.Sprintf(":%v", config.APIPort), nil) // set listen port
		if err != nil {
			log.Msgf(0, err.Error())
		}
	}
}

func getBalanceHTTPRequestHandler(w http.ResponseWriter, re *http.Request) {
	// Get the balance from the RPC server..
	balance, err := getBalance()
	if err != nil {
		log.Msgf(0, "Get Balance error:[%v]\n", err.Error())
		// No point continuing so bail out..
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Now get coin value from 3rd Party API. Ignore any errors because we have the balance so still worth sending back out..
	log.Msgf(2, "Making API call...\n")
	coinPrice, _ := getCoinPrice()

	// Calculate the total balance, to 2dp
	totalBalance := balance * coinPrice.USD
	totalBalance = float64(int(totalBalance*100)) / 100

	switch config.ResponseType {
	case jsonResponse:
		w.Header().Set("Content-Type", "application/json")
		jsonString := JSONResponse{Balance: balance, Price: coinPrice.USD, Total: totalBalance}
		jsonBytes, err := json.Marshal(jsonString)
		if err != nil {
			log.Msgf(0, "ERROR with response json marshal")
			return
		}
		w.Write([]byte(jsonBytes))
	case htmlResponse:
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("<html><head><title>Wallet RPC API</title></head><body bgcolor=#696969>"))

		// write out the wallet balance
		w.Write([]byte(fmt.Sprintf(htmlTemplateVVV, "Balance", balance, "ETN")))

		if coinPrice.USD > 0 {
			// write out the USD price
			w.Write([]byte(fmt.Sprintf(htmlTemplateVVV, "Price", "$", coinPrice.USD)))

			// write out the total value
			w.Write([]byte(fmt.Sprintf(htmlTemplateVVV, "Total Value", "$", totalBalance)))
		}

		// Close off the HTML tags
		w.Write([]byte("</body></html>"))
	}

	return
}

func getBalance() (outBal float64, oerr error) {
	outBal, err := rpc.GetBalance()
	if err != nil {
		oerr = fmt.Errorf("Error retrieving Balance from RPC Server: [%v]", err)
		return
	}
	log.Msgf(3, "Balance: %v\n", outBal)
	return
}

func getCoinPrice() (coinPrice CoinPrice, oerr error) {
	response, err := http.Get(config.APIURL)
	//var coinPrice CoinPrice
	if err != nil {
		// Log the error and return
		log.Msgf(1, "API Call error:[%v]\n", err.Error())
		oerr = err
		return
	}

	log.Msgf(2, "API call complete\n")
	// unmarshal json response into struct
	jsonData, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(jsonData, &coinPrice)
	if err != nil {
		log.Msgf(0, "Error Unmarshaling the Response from API:[%v]\n", err.Error())
		oerr = err
		return
	}

	return coinPrice, nil
}
