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

func init() {
	http.HandleFunc("/", userContentHandler)
}

func userContentHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	context := appengine.NewContext(r)
	context, _ = appengine.Namespace(context, "core")
	key := datastore.NewKey(context, "Url", path, 0, nil)

	urlEntry := new(UrlEntry)
	error := datastore.RunInTransaction(context, func(c appengine.Context) error {
		err := datastore.Get(c, key, urlEntry)
		if err != nil {
			return err
		}
		urlEntry.LastAccessed = time.Now()
		urlEntry.Clicks++
		_, err = datastore.Put(c, key, urlEntry)
		return err
	}, nil)

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
