package handler

import "net/http"

func CreateErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		badRequest(w, "Error example")
	}
}
