package gservice

import (
	"appengine/aetest"
	"fmt"
	"testing"
)

func TestStripCurrencyData(t *testing.T) {
	var html = "<div id=currency_converter_result>1000 USD = <span class=bld>1.6000 BTC</span>"
	var expected = int64(16000)
	value := stripCurrencyData(html)
	if value != expected {
		t.Error(fmt.Sprintf("Expected %d, but actual %d", expected, value))
	}
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
	if currency.rate == 0 {
		t.Error(fmt.Sprintf("Currency rate should not be 0."))
	}
}
