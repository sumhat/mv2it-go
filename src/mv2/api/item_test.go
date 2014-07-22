package api

import (
    "testing"
    T "mv2/testing"
    "time"
)

func TestItemExpired(t *testing.T) {
    assert := T.Assert(t)
    item := &ItemEntry{TimeoutSeconds: 0, DateCreated: time.Now().UTC().Add(time.Duration(-100) * time.Second)}
    assert.That(item.HasExpired()).IsFalse()
    
    item = &ItemEntry{TimeoutSeconds: 100, DateCreated: time.Now().UTC().Add(time.Duration(-1000) * time.Second)}
    assert.That(item.HasExpired()).IsTrue()
}