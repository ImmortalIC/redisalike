package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"io"
)

const apiAddr = "http://localhost:8080"

var client = &http.Client{
	Timeout: 30 * time.Second,
}

type requestResult interface {
	ReadBody(body io.ReadCloser) error
}

type storageError struct {
	Message string `json:"error"`
}

func (s storageError) Error() string {
	return s.Message
}

func (s *storageError) ReadBody(body io.ReadCloser) error {
	return readAndUnmarshall(body, s)
}

func request(requestURL string, method string, result requestResult, body io.Reader) error {
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	var resp *http.Response

	resp, err = client.Do(req)

	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return result.ReadBody(resp.Body)
	}
	var responseErr storageError
	err = responseErr.ReadBody(resp.Body)
	if err != nil {
		return err
	}
	return responseErr
}

func readAndUnmarshall(data io.ReadCloser, object interface{}) error {
	defer data.Close()
	dataBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(dataBytes, object)
}
