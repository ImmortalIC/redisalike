package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

type emtpyResponse struct {
}

type writeRequest struct {
	Value []byte        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

func (e *emtpyResponse) ReadBody(body io.ReadCloser) error {
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		return fmt.Errorf("Wrong response format")
	}
	return nil
}

func Set(key, body string) error {
	url := path.Join(apiAddr, "storage", key)
	reqBody := writeRequest{
		Value: []byte(body),
		TTL:   666 * time.Minute, //forgot about ttl when planning client. too lazy to make it normal from CL parameters
	}
	reqBodyMarshalled, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	resp := new(emtpyResponse)
	return request(url, http.MethodPost, resp, bytes.NewBuffer(reqBodyMarshalled))
}

func Remove(key string) error {
	url := path.Join(apiAddr, "storage", key)
	resp := new(emtpyResponse)
	return request(url, http.MethodDelete, resp, nil)
}
