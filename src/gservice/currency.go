package gservice

import (
	"api/net"
	"appengine"
	_ "appengine/datastore"
	"appengine/memcache"
	"encoding/json"
	"net/http"
	"net/url"
	_ "strconv"
	"strings"
	"time"
	"unicode"
)

const (
	currency_cache_prefix = "currency_exchange_"
)

func init() {
	http.HandleFunc("/gservice/currency", handleCurrency)
}

type CurrencyValue struct {
	Rate int64     `datastore:"rate,noindex" json:"rate"`
	Date time.Time `datastore:"date,noindex" json:"date"`
}

type CurrencyEntry struct {
	From string
	To   string
	Values       []CurrencyValue
}

func stripCurrencyData(html string) int64 {
	const startTag = "<span class=bld>"
	const endTag = "</span>"
	startIndex := strings.Index(html, startTag)
	if startIndex < 0 {
		return 0
	}
	startIndex += len(startTag)
	endIndex := strings.Index(html, endTag)
	if endIndex < 0 {
		return 0
	}
	html = html[startIndex:endIndex]
	value := int64(0)
	for _, v := range html {
		if unicode.IsDigit(v) {
			value *= 10
			value += int64(v - '0')
		} else if v == ' ' {
			break
		}
	}
	return value
}

func fetchCurrencyFromGFinance(c appengine.Context, fromCurrency string, toCurrency string) (value CurrencyValue, err error) {
	tUrl, err := url.Parse("https://www.google.com/finance/converter")
	if err != nil {
		return
	}
	query := tUrl.Query()
	query.Set("a", "1")
	query.Set("from", fromCurrency)
	query.Set("to", toCurrency)
	tUrl.RawQuery = query.Encode()

	httpEntry, err := net.FetchUrl(c, tUrl.String())
	html := string(httpEntry.Body)

	value.Rate = stripCurrencyData(html)
	value.Date = time.Now().UTC()
	return
}

func fetchLastCurrency(c appengine.Context, fromCurrency string, toCurrency string) CurrencyValue {
	cachedKey := currency_cache_prefix + fromCurrency + "_" + toCurrency
	cachedItem, err := memcache.Get(c, cachedKey)
	if err == nil {
		value := CurrencyValue{}
		err := json.Unmarshal(cachedItem.Value, &value)
		if err == nil {
			return value
		}
	}
	value, err := fetchCurrencyFromGFinance(c, fromCurrency, toCurrency)
	if err == nil {
		jsonValue, err := json.Marshal(value)
		if err == nil {
			memcache.Set(c, &memcache.Item{Key: cachedKey, Value: jsonValue, Expiration: 12*time.Hour})
		}
	}
	return value
}

func handleCurrency(w http.ResponseWriter, r *http.Request) {
	cUrl := r.URL
	query := cUrl.Query()
	fromCurrency := query.Get("from")
	toCurrency := query.Get("to")
	//strDays := query.Get("days")
	//numDays := 1
	//if len(strDays) > 0 {
	//    numDays, err := strconv.ParseInt(strDays, 10, 0)
	//    if err != nil {
	//        numDays = 1
	//    }
	//}
	currency := fetchLastCurrency(appengine.NewContext(r), fromCurrency, toCurrency)
	currencyEntry := &CurrencyEntry{From: fromCurrency, To: toCurrency, Values: []CurrencyValue{currency}}
	jsonValue, _ := json.Marshal(currencyEntry)
	w.Header().Set("Content-Type", "application/json")
    w.Write(jsonValue)
}
