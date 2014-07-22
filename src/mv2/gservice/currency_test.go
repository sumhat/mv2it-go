package gservice

import (
	"appengine/aetest"
	T "mv2/testing"
	"testing"
)

func TestStripCurrencyData(t *testing.T) {
	var html = "<div id=currency_converter_result>1000 USD = <span class=bld>1.6000 BTC</span>"
	var expected = int64(16000)
	value := stripCurrencyData(html)
	T.Assert(t).That(value).IsEqualTo(expected)
}

func TestFetchFromGFinance(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	currency, err := fetchCurrencyFromGFinance(c, "EUR", "USD")
	if err != nil {
		t.Error(err)
	}

	T.Assert(t).That(currency.Rate).IsGreaterThan(int64(0))
}
