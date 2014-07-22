package api

import (
	"appengine"
	"appengine/datastore"
	"time"
)

type ItemEntry struct {
	Id             string    `datastore:"-"`
	Target         string    `datastore:"target",noindex`
	Clicks         int64     `datastore:"clicks"`
	TimeoutSeconds int64     `datastore:"timeoutSeconds"`
	DateCreated    time.Time `datastore:"dateCreated"`
	LastAccessed   time.Time `datastore:"lastAccessed",noindex`
	Owner          string    `datastore:"owner"`
}

func (this *ItemEntry) GetKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Url", this.Id, 0, nil)
}

func (this *ItemEntry) LoadFromDatastore(c appengine.Context) error {
	err := datastore.Get(c, this.GetKey(c), this)
	return err
}

func (this *ItemEntry) SaveToDatastore(c appengine.Context) error {
	key := this.GetKey(c)
	_, err := datastore.Put(c, key, this)
	return err
}
