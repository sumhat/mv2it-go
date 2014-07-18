package api

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
)

const (
    config_cache_prefix = "cache_config_"
)

type ConfigEntry struct {
	Key         string    `datastore:"-"`
	Value       string    `datastore:"Value,noindex"`
}

func fetchConfigFromDataStore(c appengine.Context, key string) (value string, err error) {
	context, err := appengine.Namespace(c, "core")
	if err != nil {
	    c.Errorf("Error switching to core namespace: %", err)
		return
	}

    configEntry := new(ConfigEntry)
    dsKey := datastore.NewKey(context, "Config", key, 0, nil)
    err = datastore.Get(context, dsKey, configEntry)
    if err != nil {
        c.Errorf("Error fetching config %s: %s", key, err)
    } else {
        value = configEntry.Value
    }
    return
}

func GetConfig(c appengine.Context, key string) (value string, err error) {
    cachedKey := config_cache_prefix + key
    cachedItem, err := memcache.Get(c, cachedKey);
    if err == nil {
        value = string(cachedItem.Value)
        return
    }
    
    value, err = fetchConfigFromDataStore(c, key);
    if err == nil {
        err = memcache.Set(c, &memcache.Item{Key: cachedKey, Value: []byte(value)})
        if err != nil {
            c.Errorf("Error caching config: %s", err)
        }
    }
    return
}