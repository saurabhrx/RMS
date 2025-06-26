package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type clientError struct {
	StatusCode    int    `json:"status_code"`
	MessageToUser string `json:"message_to_user"`
}

func ResponseError(w http.ResponseWriter, statusCode int, messageToUser string) {
	logrus.Errorf("status : %d, message : %s", statusCode, messageToUser)
	clientErr := &clientError{
		StatusCode:    statusCode,
		MessageToUser: messageToUser,
	}
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientErr); err != nil {
		logrus.Errorf("failed to send the error %+v", err)
	}
}
