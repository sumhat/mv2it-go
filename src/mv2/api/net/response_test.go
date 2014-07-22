package net

import (
	"encoding/json"
	"fmt"
	"testing"
	T "mv2/testing"
)

func TestJsonResponse(t *testing.T) {
	assert := T.Assert(t)
	type data struct {
		Vi int
		Vs string
	}
	tData := &data{Vi: 10, Vs: "test"}
	response := NewJsonResponse().WithData(tData).WithError("ok")
	jResponse, err := response.Json()
	if err != nil {
		t.Error(fmt.Sprintf("Error calling Json(): %s", err))
	}
	rjData := new(JsonResponse)
	err = json.Unmarshal(jResponse, rjData)
	if err != nil {
		t.Error(fmt.Sprintf("Error unmarshalling data: %s", err))
	}
	assert.That(rjData.Error).IsEqualTo("ok")

	vi := rjData.Data.(map[string]interface{})["Vi"].(float64)
	vs := rjData.Data.(map[string]interface{})["Vs"].(string)
	assert.That(int64(vi)).IsEqualTo(int64(tData.Vi))
	assert.That(vs).IsEqualTo(tData.Vs)
}
