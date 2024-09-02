package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type responseMessage struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

func SendJson(rw http.ResponseWriter, statusCode int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(data)
}

func SendJsonResponse(rw http.ResponseWriter, statusCode int, msg, status string) {
	SendJson(rw, statusCode, responseMessage{Msg: msg, Status: status})
}

func SendError(rw http.ResponseWriter, message string, statusCode int) {
	rw.WriteHeader(statusCode)
	http.Error(rw, message, statusCode)
}

func GetIntVal(val string) (int, bool) {
	res, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return res, true
}
