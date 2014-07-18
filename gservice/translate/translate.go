package translate

import (
	"api"
	"appengine"
	"appengine/urlfetch"
	"fmt"
	"gservice"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func Init() {
	http.HandleFunc("/gservice/translate", handleTranslate)
}

func FetchTranslations(context appengine.Context, q string, srcLang string, destLang string) (data []byte, err error) {
	tUrl, err := url.Parse("https://www.googleapis.com/language/translate/v2")
	if err != nil {
		return
	}
	apiKey, err := api.GetConfig(context, "google_api_key")
	if err != nil {
		return
	}
	query := tUrl.Query()
	query.Set("key", apiKey)
	query.Set("q", q)
	query.Set("source", srcLang)
	query.Set("target", destLang)
	tUrl.RawQuery = query.Encode()

	httpRequest, err := http.NewRequest("GET", tUrl.String(), nil)
	transport := &urlfetch.Transport{
		Context:  context,
		Deadline: 60 * time.Second,
	}
	httpClient := http.Client{Transport: transport}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}

func handleTranslate(w http.ResponseWriter, r *http.Request) {
	cUrl := r.URL
	query := cUrl.Query()
	q := query.Get("q")
	srcLang := query.Get("s")
	destLang := query.Get("t")
	ck := query.Get("ck")
	if len(q) == 0 || len(destLang) == 0 || !gservice.IsValidClientKey(ck) {
		fmt.Fprint(w, "{}")
		return
	}

	context := appengine.NewContext(r)
	context.Infof("Translate: %s", q)
	v, err := FetchTranslations(context, q, srcLang, destLang)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(v)
}
