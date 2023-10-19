package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Golang-Logging-Sample/pkg/logger"
)

// Return Error
type error struct {
	ErrorMessage *messages `json:"errors"`
}

type messages struct {
	StatusCode int    `json:"stats_code"`
	Message    string `json:"message"`
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf(`%s url="%s" msg="%s"`, r.Method, r.URL, "Not Found API. Please Check to http://api_list_xxxx.com")
	errJson, err := json.Marshal(&error{
		ErrorMessage: &messages{
			StatusCode: http.StatusNotFound,
			Message:    msg,
		},
	})
	if err != nil {
		logger.Fatal(
			logger.GetApplicationError(err).
				AddMessage("Faild to Marshal from NotFoundHandler"),
		)
	}
	// ErrorResponse
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(errJson)))
	w.WriteHeader(http.StatusNotFound)
	if _, err := w.Write(errJson); err != nil {
		logger.Fatal(
			logger.GetApplicationError(err).
				AddMessage("Faild to Marshal from NotFoundHandler"),
		)
	}
}
