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

func main() {
	getTickerList()
	http.HandleFunc("/compare", compareHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
