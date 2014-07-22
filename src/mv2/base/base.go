package base

import (
	"mv2/api"
	"mv2/api/net"
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/quickadd", quickAdd)
	http.HandleFunc("/", handleUrl)
}

func quickAdd(w http.ResponseWriter, r *http.Request) {
	cUrl := r.URL
	query := cUrl.Query()
	t := query.Get("target")
	
	if !(strings.HasPrefix(t, "http://") || strings.HasPrefix(t, "https://")) {
		fmt.Fprint(w, "{}")
		return
	}
	
	itemEntry := &api.ItemEntry{
		Target: t,
		TimeoutSeconds: 2700000,
		DateCreated: time.Now().UTC(),
		LastAccessed: time.Now().UTC(),
	}
	
	context := appengine.NewContext(r)
	context, err := appengine.Namespace(context, "core")
	if err != nil {
		fmt.Fprint(w, err)
	}
	itemEntry, err = api.AddNewItem(context, itemEntry)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, itemEntry.Id)
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
		fmt.Fprintf(w, "Something wrong, please refresh the page.")
		return
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
