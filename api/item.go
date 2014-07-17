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
	Owner	string `datastore:"owner"`
}

func LoadItemEntry(c appengine.Context, id string) (*ItemEntry, error) {
	itemEntry := new(ItemEntry)
	context, err := appengine.Namespace(c, "core")
	if err != nil {
		return urlEntry, err
	}

	itemEntry.Id = id
	key := datastore.NewKey(context, "Url", itemEntry.Id, 0, nil)
	err = datastore.RunInTransaction(context, func(tc appengine.Context) error {
		err := datastore.Get(tc, key, itemEntry)
		if err != nil {
			return err
		}
		itemEntry.LastAccessed = time.Now()
		itemEntry.Clicks++
		_, err = datastore.Put(tc, key, itemEntry)
		return err
	}, nil)
	return itemEntry, err
}