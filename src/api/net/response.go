package net

import (
    "encoding/json"
)

type JsonResponse struct {
    Data interface{} `json:"data"`
    Error string `json:"error"`
}

func NewJsonResponse() *JsonResponse {
    return new(JsonResponse)
}

func (this *JsonResponse) WithData(v interface{}) *JsonResponse {
    this.Data = v
    return this
}

func (this *JsonResponse) WithError(err string) *JsonResponse {
    this.Error = err
    return this
}

func (this *JsonResponse) Json() (value []byte, err error) {
    value, err = json.Marshal(this)
    return
}