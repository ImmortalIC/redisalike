package api

import (
	"encoding/json"
	"fmt"
	"github.com/ImmortalIC/redisalike/storage"
	"io/ioutil"
	"net/http"
	"strings"
)

func StartServer() error {
	http.HandleFunc("/storage/", storageHandler)
	http.HandleFunc("/keys", keysHandler)
	return http.ListenAndServe(":8080", nil)
}

func storageHandler(response http.ResponseWriter, request *http.Request) {
	uri := strings.TrimPrefix(request.URL.Path, "/storage/")
	uriParts := strings.Split(uri, "/")
	key := uriParts[0]
	var index string
	if len(uriParts) > 1 {
		index = uriParts[1]
	}
	switch request.Method {
	case http.MethodGet:
		responseBody, err := Get(key, index)
		if err != nil {
			if strings.HasPrefix(err.Error(), "No such key") {
				response.WriteHeader(http.StatusNotFound)
			} else {
				response.WriteHeader(http.StatusInternalServerError)
			}
			response.Write(marshallError(err))
			return
		}
		body, err := responseBody.JSON()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(marshallError(err))
			return
		}
		response.WriteHeader(http.StatusOK)
		response.Write(body)
	case http.MethodPost:
		defer request.Body.Close()
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(marshallError(err))
			return
		}
		err = Post(key, body)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write(marshallError(err))
			return
		}
		response.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		err := storage.Remove(key)
		if err != nil {
			if strings.HasPrefix(err.Error(), "No such key") {
				response.WriteHeader(http.StatusNotFound)
			} else {
				response.WriteHeader(http.StatusInternalServerError)
			}
			response.Write(marshallError(err))
			return
		}
		response.WriteHeader(http.StatusNoContent)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write(marshallError(fmt.Errorf("Method not allowed")))
	}

}

func keysHandler(response http.ResponseWriter, request *http.Request) {
	keys := storage.Keys()
	responseBody, err := json.Marshal(KeysResponse{Keys: keys})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write(marshallError(err))
	}
	response.WriteHeader(http.StatusOK)
	response.Write(responseBody)
}

func Get(key, index string) (Responser, error) {
	result, err := storage.Get(key)
	if err != nil {
		return nil, err
	}
	if index == "" {
		switch result.Value().(type) {
		case string:
			return StringResponse{
				Value: result.Value().(string),
			}, nil
		case []string:
			return ListResponse{
				Value: result.Value().([]string),
			}, nil
		case map[string]string:
			return DictResponse{
				Value: result.Value().(map[string]string),
			}, nil
		default:
			panic("Somethig gone wrong")
		}
	}
	indexedValue, err := result.ByKey(index)
	if err != nil {
		return nil, err
	}
	return StringResponse{
		Value: indexedValue,
	}, nil
}

func Post(key string, body []byte) error {
	var unmarshalledBody WriteRequestBody
	err := json.Unmarshal(body, &unmarshalledBody)
	if err != nil {
		return err
	}
	switch unmarshalledBody.Value[0] {
	case '[':
		var list []string
		err := json.Unmarshal([]byte(unmarshalledBody.Value), &list)
		if err != nil {
			return err
		}
		return storage.Set(key, list, unmarshalledBody.TimeToTile)
	case '{':
		var dict map[string]string
		err := json.Unmarshal([]byte(unmarshalledBody.Value), &dict)
		if err != nil {
			return err
		}
		return storage.Set(key, dict, unmarshalledBody.TimeToTile)
	default:
		return storage.Set(key, unmarshalledBody.Value, unmarshalledBody.TimeToTile)
	}

}
