package rpc

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"walletapi/config"
	"walletapi/log"
)

// Height - struct for getheight response
type Height struct {
	ID      string       `json:"id"`
	JSONRPC string       `json:"jsonrpc"`
	Result  HeightResult `json:"result"`
}

// HeightResult - struct for the actual result of get height
type HeightResult struct {
	Height  uint64 `json:"height,attr"`
	Balance uint64 `json:"balance,attr"`
}

// GetHeight gets the current block height from the local wallet
func GetHeight() (height uint64, oerr error) {
	var jsonStr = []byte(`{"jsonrpc": "2.0", "id": "0", "method": "getheight"}`)
	jsonData, _ := postRPC(jsonStr)

	var h Height
	// Now try and unmarshall the response into the struct
	e := json.Unmarshal(jsonData, &h)
	if e != nil {
		log.Msgf(0, "ERROR: %v\n", e)
		return
	}
	return h.Result.Height, nil
}

// GetBalance gets the current balance from the local wallet
func GetBalance() (balance float64, oerr error) {

	var jsonStr = []byte(`{"jsonrpc": "2.0", "id": "0", "method": "getbalance"}`)
	jsonData, _ := postRPC(jsonStr)

	var h Height
	// Now try and unmarshall the response into the struct
	e := json.Unmarshal(jsonData, &h)
	if e != nil {
		log.Msgf(0, "ERROR: %v\n", e)
		oerr = e
		return
	}
	return float64(h.Result.Balance) / 100, nil
}

func postRPC(jsonStr []byte) (data []byte, oerr error) {
	req, err := http.NewRequest("POST", config.RPCServerURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		oerr = err
		return
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)
	return jsonData, err
}
