package mv2

import (
	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type UrlEntry struct {
	Id             string    `datastore:"-"`
	Target         string    `datastore:"target"`
	Clicks         int64     `datastore:"clicks"`
	TimeoutSeconds int64     `datastore:"timeoutSeconds"`
	DateCreated    time.Time `datastore:"dateCreated"`
	LastAccessed   time.Time `datastore:"lastAccessed",noindex`
	Operation      string    `datastore:"-"`
}

func LoadUrlEntry(c appengine.Context, id string) (*UrlEntry, error) {
	urlEntry := new(UrlEntry)
	context, err := appengine.Namespace(c, "core")
	if err != nil {
		return urlEntry, err
	}

	urlEntry.Id = id
	key := datastore.NewKey(context, "Url", urlEntry.Id, 0, nil)
	err = datastore.RunInTransaction(context, func(tc appengine.Context) error {
		err := datastore.Get(tc, key, urlEntry)
		if err != nil {
			return err
		}
		urlEntry.LastAccessed = time.Now()
		urlEntry.Clicks++
		_, err = datastore.Put(tc, key, urlEntry)
		return err
	}, nil)
	return urlEntry, err
}

func FetchUrl(url string, context appengine.Context) (body []byte, contentType string, err error) {
	httpRequest, err := http.NewRequest("GET", url, nil)
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
	contentType = resp.Header.Get("Content-Type")
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func init() {
	http.HandleFunc("/", userContentHandler)
}

func userContentHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	context := appengine.NewContext(r)
	urlEntry, error := LoadUrlEntry(context, path)

	if error != nil {
		fmt.Fprint(w, path, " not found: ", error)
		return
	}

	if strings.HasPrefix(urlEntry.Target, "c:") {
		url2Fetch := urlEntry.Target[2:]
		body, contentType, err := FetchUrl(url2Fetch, appengine.NewContext(r))
		if err != nil {
			fmt.Fprint(w, "error: ", err)
			return
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(body)
		return
	}

	http.Redirect(w, r, urlEntry.Target, 301)
}
