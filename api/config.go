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
	Key         string    `datastore:"Key"`
	Value string `datastore:"Value,noindex"`
}

func fetchConfigFromDataStore(c appengine.Context, key string) (value string, err error) {
	context, err := appengine.Namespace(c, "core")
	if err != nil {
		return
	}

	query := datastore.NewQuery("Config").Filter("Key =", key)
	var configs []ConfigEntry
    _, err = query.GetAll(context, &configs)
    if len(configs) > 0 {
        value = configs[0].Value
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