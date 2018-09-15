package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var tickerSymbolDict = map[string]string{
	"BTC": "1",
	"ETH": "1027",
	"XRP": "52",
	"LTC": "2",
	"BCH": "1831",
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
	http.HandleFunc("/compare", compareHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
