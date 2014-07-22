package gservice

import (
    "testing"
    T "mv2/testing"
)

func TestGetTranslationUrl(t *testing.T) {
    assert := T.Assert(t)
    url, err := getTranslationUrl("key", "query", "source", "dest")
    assert.That(err).IsNil()
    assert.That(url).IsEqualTo("https://www.googleapis.com/language/translate/v2?key=key&q=query&source=source&target=dest")
}