package utils

import (
	"RMS/models"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"math"
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

func HaversineDistance(body models.Distance) float64 {
	const radius float64 = 6371
	lat1 := body.UserLat * math.Pi / 180
	lon1 := body.UserLong * math.Pi / 180
	lat2 := body.RestaurantLat * math.Pi / 180
	lon2 := body.RestaurantLong * math.Pi / 180

	dlat := lat2 - lat1
	dlong := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlong/2)*math.Sin(dlong/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return radius * c

}
