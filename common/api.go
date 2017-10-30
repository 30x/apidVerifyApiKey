package common

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func WriteError(w http.ResponseWriter, reason string, errorCode int, statusCode int) {
	w.WriteHeader(statusCode)
	resp := ErrorResponse{
		ResponseCode:    strconv.Itoa(errorCode),
		ResponseMessage: reason,
		StatusCode:      statusCode,
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte("unable to marshal errorResponse: " + err.Error()))
	} else {
		w.Write(bytes)
	}
}
