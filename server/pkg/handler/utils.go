package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeErrorResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode(&errorResponse{Error: message})
	if err != nil {
		log.Println(err)
	}
}

func badRequest(w http.ResponseWriter, message string) {
	writeErrorResponse(w, message, http.StatusBadRequest)
}

func unauthorizedRequest(w http.ResponseWriter, message string) {
	writeErrorResponse(w, message, http.StatusUnauthorized)
}

func writeResponse(w http.ResponseWriter, body interface{}) {
	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.Println(err)
	}
}

var ParamNotFound = errors.New("parameter not found")
var CannotParseParam = errors.New("cannot parse param")

func getRequiredStringFromVars(w http.ResponseWriter, vars map[string]string, paramName string) (string, error) {
	value, ok := vars[paramName]
	if !ok {
		badRequest(w, fmt.Sprintf("please provide `%s`", paramName))
		return "", ParamNotFound
	}

	return value, nil
}

func getRequiredIntFromVars(w http.ResponseWriter, vars map[string]string, paramName string) (int, error) {
	stringValue, err := getRequiredStringFromVars(w, vars, paramName)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(stringValue, 10, 64)
	if err != nil {
		log.Println(fmt.Errorf("cannot parse `%s`: %w", paramName, err))
		badRequest(w, fmt.Sprintf("cannot parse `%s`", paramName))
		return 0, CannotParseParam
	}

	return int(value), nil
}
