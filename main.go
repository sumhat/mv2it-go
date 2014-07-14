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
	urlEntry.Id = id
	return urlEntry, err
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
		httpRequest, err := http.NewRequest("GET", url2Fetch, nil)
		transport := &urlfetch.Transport{
			Context:  appengine.NewContext(r),
			Deadline: 60 * time.Second,
		}
		httpClient := http.Client{Transport: transport}
		resp, err := httpClient.Do(httpRequest)
		if err != nil {
			//http.Redirect(w, r, url2Fetch, 302)
			fmt.Fprint(w, "error: ", err)
			return
		}
		defer resp.Body.Close()
		contentType := resp.Header["Content-Type"]
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprint(w, "error: ", err)
			return
		}
		w.Header()["Content-Type"] = contentType
		w.Write(body)
		return
	}

	http.Redirect(w, r, urlEntry.Target, 301)
}
