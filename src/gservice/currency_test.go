package gservice

import (
    "testing"
    "fmt"
)

func TestStripCurrencyData(t *testing.T) {
    var html = "<div id=currency_converter_result>1000 USD = <span class=bld>1.6000 BTC</span>"
    var expected = int64(16000)
    value := stripCurrencyData(html)
    if value != expected {
        t.Error(fmt.Sprintf("Expected %d, but actual %d", expected, value))
    }
}