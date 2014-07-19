package base

import (
	"api"
	"api/net"
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/", handleUrl)
}

func handleUrl(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	if len(path) == 0 {
		http.Redirect(w, r, "http://leonax.net", 302)
		return
	}
	context := appengine.NewContext(r)
	namedContext, err := appengine.Namespace(context, "core")
	if err != nil {
		context.Errorf("Error switching to core namespace: %s", err)
	}
	urlEntry := &api.ItemEntry{Id: path}
	err = datastore.RunInTransaction(namedContext, func (tc appengine.Context) error {
		err := urlEntry.LoadFromDatastore(tc)
		if err != nil {
			return err
		}
		urlEntry.LastAccessed = time.Now().UTC()
		urlEntry.Clicks++
		return urlEntry.SaveToDatastore(tc)
	}, nil)

	if err != nil {
		fmt.Fprintf(w, "%s not found: %s", path, err)
		return
	}

	dispatchUrl(w, r, context, urlEntry.Target)
}

func dispatchUrl(w http.ResponseWriter, r *http.Request, c appengine.Context, url string) {
	if strings.HasPrefix(url, "c:") {
		actualUrl := url[2:]
		httpEntry, err := net.FetchUrl(c, actualUrl)
		if err != nil {
			fmt.Fprintf(w, "Error fetching Url (%s): %s", actualUrl, err)
			return
		}
		w.Header().Set("Content-Type", httpEntry.ContentType)
		w.Write(httpEntry.Body)
		return
	}
	
	http.Redirect(w, r, url, 301)
}
