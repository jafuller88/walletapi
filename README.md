
## Wallet API
This API can be used to retrieve crypto coin local wallet balances and calculate the current total value.
It is based on Monero however the examples used are Electroneum based (Electroneum is a fork of Monero)

The data returned can be in either JSON or HTML format.

## Usage
 -logLevel int - Log Level (0=FATAL,1=ERROR,2=INFO,3=DEBUG  
 -rpc string - RPC URL and Port (default "http://127.0.0.1:3000/json_rpc")  
 -port string - Web API Server Port (default "8080")  
 -rt int - API response type (0=JSON,1=HTML)  
 -api string - Price API URL (default "https://min-api.cryptocompare.com/data/price?")  
 -tsym string - To Symbol/s (default "BTC,USD,EUR")  
 -fsym string - From Symbol (default "ETN")  

## Running the wallet daemon services
For the API to run you will need to have the main Wallet Daemon running and synchronized with the network. Example startup command:
```
./electroneumd --data-dir=/path/to/lmdb/directory
```
Once the Daemon is running and synchronized you will need to start the Wallet RPC Server. Example startup command:
```
./electroneum-wallet-rpc --password xxxx --restricted-rpc --rpc-bind-port 3000 --disable-rpc-login --wallet-file /path/to/wallet/file/walletfile.etn
```
You can then build and start the Wallet API
```
go build && ./walletapi
```

## Security - IMPORTANT!
Ideally you would run this on your local network only.
However if you do choose to expose this API externally then please ensure the following security measures are taken to avoid any trouble :-)

1. Ensure your wallet private keys are not kept on the server running this API or the server running your wallet deamon services
2. Run this API on a different server to where you are hosting the Wallet RPC server
3. Ensure that you set the --restricted-rpc flag when running the RPC server. this restricts to view-only commands
4. Lockdown access to the API server and wallet servers by port on your network

## Wallet RPC Calls
GET Height
```
curl -X POST http://127.0.0.1:3000/json_rpc -d '{"jsonrpc":"2.0","id":"0","method":"getheight"}' -H 'Content-Type: application/json'
```
Response
```
{
  "id": "0",
  "jsonrpc": "2.0",
  "result": {
	"height": 994310
  }
}
```
GET Balance
```
curl -X POST http://127.0.0.1:3000/json_rpc -d '{"jsonrpc:"2.0","id":"0","method":"getbalance"}' -H 'Content-Type: application/json'
```
Response
```
{
  "id": "0",
  "jsonrpc": "2.0",
  "result": {
    "balance": 138834,
    "unlocked_balance": 138834
}
```

## Crypto Coin Price API Calls
API for Prices - https://www.cryptocompare.com/api/#introduction

Example:

Request : https://min-api.cryptocompare.com/data/price?fsym=ETN&tsyms=BTC,USD,EUR

Response: {"BTC":0.00000827,"USD":0.08392,"EUR":0.06694}