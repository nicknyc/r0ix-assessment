package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/keypair"
)

type Tickers struct {
	Data []Ticker `json:"data"`
}

type Ticker struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type TickerDetails struct {
	Data TickerDetail `json:"data"`
}

type TickerDetail struct {
	Qoutes QouteJson `json:"quotes"`
}

type QouteJson struct {
	USD USDJson `json:"USD"`
}

type USDJson struct {
	PC float32 `json:"percent_change_1h"`
}

var tickerSymbolDict = map[string]Ticker{}

//TODO: handle errors properly
func getTickerList() {
	url := "https://api.coinmarketcap.com/v2/listings/"
	res, _ := http.Get(url)
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	var tickers Tickers
	err := json.Unmarshal(content, &tickers)
	if err != nil {
		fmt.Println(err)
	} else {
		for _, v := range tickers.Data {
			tickerSymbolDict[v.Symbol] = v
		}
	}
}

func getTickerChange(id int) float32 {
	url := fmt.Sprint("https://api.coinmarketcap.com/v2/ticker/", id, "/")
	res, _ := http.Get(url)
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	var detail TickerDetails
	err := json.Unmarshal(content, &detail)
	if err != nil {
		fmt.Println(err)
	}
	return detail.Data.Qoutes.USD.PC
}

func compareHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if p1, p2 := q.Get("ticker_symbol_1"), q.Get("ticker_symbol_2"); tickerSymbolDict[p1].Id > 0 && tickerSymbolDict[p2].Id > 0 {
		if t1, t2 := tickerSymbolDict[p1], tickerSymbolDict[p2]; getTickerChange(t1.Id) > getTickerChange(t2.Id) {
			fmt.Fprintf(w, t1.Name)
		} else {
			fmt.Fprintf(w, t2.Name)
		}
	} else {
		http.Error(w, "Bad request", 400)
	}
}

func keyGenerateHandler(w http.ResponseWriter, r *http.Request) {
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, "Seed      : ", pair.Seed())
	fmt.Fprintln(w, "Public Key: ", pair.Address())
}

func getAccountHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("key")
	fmt.Println("create account ", address)

	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + address)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, string(body))
}

func accountDetailHandler(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("key")
	fmt.Println("Details of account ", address)
	account, err := horizon.DefaultTestNetClient.LoadAccount(address)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(w, "Balances for account:", address)

	for _, balance := range account.Balances {
		fmt.Fprintln(w, balance)
	}
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	source := q.Get("source_account")
	destination := q.Get("destination_account")
	amount := q.Get("amount")

	fmt.Println("Start transfering ", amount, " from ", source, " to ", destination)

	if _, err := horizon.DefaultTestNetClient.LoadAccount(destination); err != nil {
		fmt.Fprintln(w, "Error: ", err)
		panic(err)
	}

	tx, err := build.Transaction(
		build.TestNetwork,
		build.SourceAccount{source},
		build.AutoSequence{horizon.DefaultTestNetClient},
		build.Payment(
			build.Destination{destination},
			build.NativeAmount{amount},
		),
	)

	if err != nil {
		fmt.Fprintln(w, "Error: ", err)
		panic(err)
	}

	// Sign the transaction to prove you are actually the person sending it.
	txe, err := tx.Sign(source)
	if err != nil {
		fmt.Fprintln(w, "Error : ", err)
		panic(err)
	}

	txeB64, err := txe.Base64()
	if err != nil {
		fmt.Fprintln(w, "Error: ", err)
		panic(err)
	}

	// And finally, send it off to Stellar!
	resp, err := horizon.DefaultTestNetClient.SubmitTransaction(txeB64)
	if err != nil {
		fmt.Fprintln(w, "Error: ", err)
		panic(err)
	}

	fmt.Fprintln(w, "Successful Transaction:")
	fmt.Fprintln(w, "Ledger:", resp.Ledger)
	fmt.Fprintln(w, "Hash:", resp.Hash)
}

func main() {
	getTickerList()
	http.HandleFunc("/compare", compareHandler)
	http.HandleFunc("/keygen/", keyGenerateHandler)
	http.HandleFunc("/account", getAccountHandler)
	http.HandleFunc("/accountDetail", accountDetailHandler)
	http.HandleFunc("/transfer", transferHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
