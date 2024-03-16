package application

import (
	"encoding/json"
	"net/http"
	"strings"
)

type responseMessage struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

func (a app) validateToken(rw http.ResponseWriter, tokenReq string) bool {
	tokenReq = strings.TrimSpace(tokenReq)
	if a.token != tokenReq {
		sendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
		return false
	}
	return true
}

func sendJson(rw http.ResponseWriter, statusCode int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(data)
}

func sendJsonResponse(rw http.ResponseWriter, statusCode int, msg, status string) {
	sendJson(rw, statusCode, responseMessage{Msg: msg, Status: status})
}

func sendError(rw http.ResponseWriter, message string, statusCode int) {
	rw.WriteHeader(statusCode)
	http.Error(rw, message, statusCode)
}
