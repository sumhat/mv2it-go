package translate

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
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
	query := tUrl.Query()
	query.Set("key", "AIzaSyDnDbRY3cM8K2I9GNycPFRuUqsX_u9vH1g")
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
	if len(q) == 0 || len(destLang) == 0 {
		fmt.Fprint(w, "{}")
		return
	}

	v, err := FetchTranslations(appengine.NewContext(r), q, srcLang, destLang)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(v)
}