package gservice

import (
	"appengine"
	"fmt"
	"mv2/api"
	"mv2/api/net"
	"net/http"
	"net/url"
)

func init() {
	http.HandleFunc("/gservice/translate", handleTranslate)
}

func getTranslationUrl(apiKey string, q string, srcLang string, destLang string) (string, error) {
	tUrl, err := url.Parse("https://www.googleapis.com/language/translate/v2")
	if err != nil {
		return "", err
	}
	query := tUrl.Query()
	query.Set("key", apiKey)
	query.Set("q", q)
	query.Set("source", srcLang)
	query.Set("target", destLang)
	tUrl.RawQuery = query.Encode()
	
	return tUrl.String(), nil
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
	apiKey, err := api.GetConfig(context, "google-api-key")
	if err != nil {
		return
	}
	tUrl, err:= getTranslationUrl(apiKey, q, srcLang, destLang)
	if err != nil {
		context.Errorf("Error creating translation url: %v", err)
		fmt.Fprint(w, err)
		return
	}
	v, err := net.FetchUrl(context, tUrl)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", v.ContentType)
	w.Write(v.Body)
}
