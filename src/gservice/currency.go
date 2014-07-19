package gservice

import (
    "appengine"
    _ "appengine/datastore"
    "appengine/memcache"
    "appengine/urlfetch"
    _ "strconv"
    "strings"
    "unicode"
    "time"
    "net/http"
    "net/url"
    "io/ioutil"
)

const (
    currency_cache_prefix = "currency_exchange_"
)

func init() {
    http.HandleFunc("/gservice/currency", handleCurrency)
}

type CurrencyValue struct {
    rate int64 `datastore:"rate,noindex"`
    date time.Time `datastore:"date,noindex"`
}

type CurrencyEntry struct {
    fromCurrency string
    toCurrency string
    values *[]CurrencyValue
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
	html = html[startIndex : endIndex]
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

func fetchCurrencyFromGFinance(c appengine.Context, fromCurrency string, toCurrency string) (value *CurrencyValue, err error) {
    tUrl, err := url.Parse("https://www.google.com/finance/converter")
	if err != nil {
		return
	}
	query := tUrl.Query()
	query.Set("a", "1")
	query.Set("from", fromCurrency)
	query.Set("to", toCurrency)
	tUrl.RawQuery = query.Encode()

	httpRequest, err := http.NewRequest("GET", tUrl.String(), nil)
	transport := &urlfetch.Transport{
		Context:  c,
		Deadline: 60 * time.Second,
	}
	httpClient := http.Client{Transport: transport}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	html := string(data)
	
	value = new(CurrencyValue)
	value.rate = stripCurrencyData(html)
	value.date = time.Now().UTC()
	return
}
    

func fetchCurrency(c appengine.Context, fromCurrency string, toCurrency string, days int) {
    cachedKey := currency_cache_prefix + fromCurrency + "_" + toCurrency
    cachedItem, err := memcache.Get(c, cachedKey);
    if err == nil {
    }
    _ = cachedItem
}

func handleCurrency(w http.ResponseWriter, r *http.Request) {
	//cUrl := r.URL
	//query := cUrl.Query()
	//fromCurrency := query.Get("from")
	//toCurrency := query.Get("to")
	//strDays := query.Get("days")
	//numDays := 1
	//if len(strDays) > 0 {
    //    numDays, err := strconv.ParseInt(strDays, 10, 0)
	//    if err != nil {
	//        numDays = 1
	//    }
	//}
}