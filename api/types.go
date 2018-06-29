package api

import (
	"encoding/json"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}
type Responser interface {
	JSON() ([]byte, error)
}

type StringResponse struct {
	Value string `json:"value"`
}
type ListResponse struct {
	Value []string `json:"value"`
}
type DictResponse struct {
	Value map[string]string `json:"value"`
}

func (r StringResponse) JSON() ([]byte, error) {
	return json.Marshal(r)
}

func (r ListResponse) JSON() ([]byte, error) {
	return json.Marshal(r)
}

func (r DictResponse) JSON() ([]byte, error) {
	return json.Marshal(r)
}

func marshallError(err error) []byte {
	errorObject := ErrorResponse{
		Error: err.Error(),
	}
	res, _ := json.Marshal(errorObject)
	return res
}

type WriteRequestBody struct {
	Value      []byte         `json:"value"`
	TimeToTile *time.Duration `json:"ttl"`
}

type KeysResponse struct {
	Keys []string `json:"keys"`
}
