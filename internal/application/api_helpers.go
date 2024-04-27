package application

import (
	"encoding/json"
	"net/http"
	"fmt"	
	"strconv"
)

type responseMessage struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
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

func getIntVal(val string) (int, bool) {
	res, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return res, true
}
