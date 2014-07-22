package api

import (
    "appengine"
    "appengine/aetest"
    "appengine/datastore"
    "testing"
    T "mv2/testing"
)

func TestFetchFromDatastore(t *testing.T) {
    assert := T.Assert(t)
    
    c, err := aetest.NewContext(nil)
    assert.That(err).IsNil()
    c2, err := appengine.Namespace(appengine.Context(c), "core")
    assert.That(err).IsNil()
	defer c.Close()

	name := "test"
	key := datastore.NewKey(c2, "Config", name, 0, nil)
	config := &ConfigEntry{Value: "abcd"}
	_, err = datastore.Put(c2, key, config)
	assert.That(err).IsNil()
	
	value2, err := fetchConfigFromDataStore(c2, name)
	assert.That(err).IsNil()
	assert.That(value2).IsEqualTo("abcd")
}
	