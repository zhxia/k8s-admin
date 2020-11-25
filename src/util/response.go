package util

import (
    "encoding/json"
    "net/http"
)

const (
    ResultOk    = 0
    ResultError = -1
)

type JsonResponse struct {
    Code int32       `json:"code"`
    Data interface{} `json:"data"`
    Msg  string      `json:"msg"`
}

func (resp *JsonResponse) Output(w http.ResponseWriter) error {
    data, err := json.Marshal(resp)
    if err != nil {
        return err
    }
    w.Header().Add("Server", "AKC-Server")
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, err = w.Write(data)
    return err
}
func NewJsonResponse(data interface{}, msg string, code int32) *JsonResponse {
    if msg == "" {
        msg = "success"
    }
    jr := JsonResponse{
        Code: code,
        Data: data,
        Msg:  msg,
    }
    return &jr
}
