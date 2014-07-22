package net

import (
	"appengine"
	"appengine/urlfetch"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpEntry struct {
	ContentType string
	Body        []byte
}

func FetchUrl(context appengine.Context, url string) (httpEntry *HttpEntry, err error) {
	httpRequest, err := http.NewRequest("GET", url, nil)
	transport := &urlfetch.Transport{
		Context:  context,
		Deadline: 60 * time.Second,
	}
	httpClient := http.Client{Transport: transport}
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		context.Errorf("Error fetching url: %v", err)
		return
	}
	defer resp.Body.Close()
	httpEntry = new(HttpEntry)
	httpEntry.ContentType = resp.Header.Get("Content-Type")
	httpEntry.Body, err = ioutil.ReadAll(resp.Body)
	return
}
