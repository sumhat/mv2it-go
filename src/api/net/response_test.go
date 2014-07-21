package net

import (
    "testing"
    "encoding/json"
    "fmt"
)

func TestJsonResponse(t *testing.T) {
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
    if rjData.Error != "ok" {
        t.Error(fmt.Sprintf("Expected ok, but actually %s", rjData.Error))
    }
    
    vi := rjData.Data.(map[string]interface{})["Vi"].(float64)
    vs := rjData.Data.(map[string]interface{})["Vs"].(string)
    if int64(vi) != int64(tData.Vi) {
        t.Error(fmt.Sprintf("Expected %d, but actually %d", tData.Vi, int64(vi)))
    }
    if vs != tData.Vs {
        t.Error(fmt.Sprintf("Expected %s, but actually %s", tData.Vs, vs))
    }
}
    