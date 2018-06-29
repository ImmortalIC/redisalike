package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
)

type responseValue struct {
	ValueJSON string `json:"value"`
}

func Get(key string, index string) (string, error) {
	url := path.Join(apiAddr, "storage", key)
	if index != "" {
		url = path.Join(url, index)
	}
	response := new(responseValue)
	err := request(url, http.MethodGet, response, nil)
	return response.ValueJSON, err
}

func (r *responseValue) ReadBody(body io.ReadCloser) error {
	const prefix = `{"value":`
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	dataString := string(data)

	if !strings.HasPrefix(dataString, prefix) {
		return fmt.Errorf("Wrong response format for Get method")
	}
	r.ValueJSON = strings.TrimRight(strings.TrimPrefix(dataString, prefix), "}")
	r.ValueJSON = strings.Trim(r.ValueJSON, `"`)
	return nil
}
