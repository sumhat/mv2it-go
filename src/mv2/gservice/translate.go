package gservice

import (
	"mv2/api"
	"mv2/api/net"
	"appengine"
	"fmt"
	"net/http"
	"net/url"
)

func init() {
	http.HandleFunc("/gservice/translate", handleTranslate)
}

func fetchTranslations(context appengine.Context, q string, srcLang string, destLang string) (data *net.HttpEntry, err error) {
	tUrl, err := url.Parse("https://www.googleapis.com/language/translate/v2")
	if err != nil {
		return
	}
	apiKey, err := api.GetConfig(context, "google-api-key")
	if err != nil {
		return
	}
	query := tUrl.Query()
	query.Set("key", apiKey)
	query.Set("q", q)
	query.Set("source", srcLang)
	query.Set("target", destLang)
	tUrl.RawQuery = query.Encode()
	
	data, err = net.FetchUrl(context, tUrl.String())
	return
}

func handleTranslate(w http.ResponseWriter, r *http.Request) {
	cUrl := r.URL
	query := cUrl.Query()
	q := query.Get("q")
	srcLang := query.Get("s")
	destLang := query.Get("t")
	ck := query.Get("ck")
	if len(q) == 0 || len(destLang) == 0 || !IsValidClientKey(ck) {
		fmt.Fprint(w, "{}")
		return
	}

	context := appengine.NewContext(r)
	context.Infof("Translate: %s", q)
	v, err := fetchTranslations(context, q, srcLang, destLang)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", v.ContentType)
	w.Write(v.Body)
}
