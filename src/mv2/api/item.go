package api

import (
	"appengine"
	"appengine/datastore"
	"time"
	"strconv"
	"strings"
)

const (
	intIdPrefix = "z0"
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

func AddNewItem(c appengine.Context, itemEntry *ItemEntry) (*ItemEntry, error) {
	key := datastore.NewIncompleteKey(c, "Url", nil)
	key, err := datastore.Put(c, key, itemEntry)
	if err != nil {
		return nil, err
	}
	intId := key.IntID()
	strId := intIdPrefix + strconv.FormatInt(intId, 36)
	itemEntry.Id = strId
	return itemEntry, nil
}

func parseId(strId string) (string, int64, error) {
	if !strings.HasPrefix(strId, intIdPrefix) {
		return strId, int64(0), nil
	}
	strId = strId[len(intIdPrefix):]
	intId, err := strconv.ParseInt(strId, 36, 64)
	return "", intId, err
}

func (this *ItemEntry) GetKey(c appengine.Context) *datastore.Key {
	strId, intId, err := parseId(this.Id)
	if err != nil {
		c.Errorf("Error parsing item id (%s): %v", this.Id, err)
	}
	return datastore.NewKey(c, "Url", strId, intId, nil)
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
