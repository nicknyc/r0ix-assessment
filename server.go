package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Tickers struct {
	Data []Ticker `json:"data"`
}

type Ticker struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

var tickerSymbolDict = map[string]string{}

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
			tickerSymbolDict[v.Symbol] = fmt.Sprint(v.Id)
			fmt.Println(v.Symbol, v.Id)
		}
	}
}

func getTickerChange(id string) string {
	url := "https://api.coinmarketcap.com/v2/ticker/" + id + "/"
	fmt.Println(url)
	res, _ := http.Get(url)
	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	return string(content)
}

func compareHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if p1, p2 := q.Get("ticker_symbol_1"), q.Get("ticker_symbol_2"); tickerSymbolDict[p1] != "" && tickerSymbolDict[p2] != "" {
		fmt.Fprintf(w, getTickerChange(tickerSymbolDict[p1]))
		fmt.Fprintf(w, getTickerChange(tickerSymbolDict[p2]))
	} else {
		http.Error(w, "Bad request", 400)
	}
}

func main() {
	getTickerList()
	http.HandleFunc("/compare", compareHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
