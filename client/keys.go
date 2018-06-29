package client

import (
	"io"
	"net/http"
	"path"
)

type keysResponse struct {
	Keys []string `json:"keys"`
}

func (k *keysResponse) ReadBody(body io.ReadCloser) error {
	return readAndUnmarshall(body, k)
}

func Keys() ([]string, error) {
	url := path.Join(apiAddr, "keys")
	response := new(keysResponse)
	err := request(url, http.MethodGet, response, nil)
	return response.Keys, err
}
