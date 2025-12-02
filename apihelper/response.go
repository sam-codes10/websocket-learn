package apihelper

import "net/http"

type ApiRes struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Code    int      `json:"code"`
	Data    interface{} `json:"data"`
}

func SendInternalServerError(message string) (int, ApiRes) {
	return http.StatusInternalServerError, ApiRes{
		Status: false,
		Message: message,
		Code: http.StatusInternalServerError,
	}
}


