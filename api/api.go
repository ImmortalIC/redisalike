package api

import (
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
	case http.MethodPost:
	case http.MethodDelete:
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
		response.Write([]byte("Method not allowed"))
	}

}
