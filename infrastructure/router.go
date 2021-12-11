package infrastructure

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	h "github.com/marcy-t/Golang-Logging-Sample/pkg/handlers"
	"github.com/marcy-t/Golang-Logging-Sample/pkg/logger"
)

// Router
func NewRouter(controller *ControllHandler) (root *mux.Router) {
	root = mux.NewRouter()
	root.NotFoundHandler = http.HandlerFunc(h.NotFoundHandler)
	eh := errorRoutingDetected
	// PathPrefix
	api := root.PathPrefix("/api/v1/").Subrouter()
	// SamplePath
	common := controller.Common
	api.HandleFunc("/ping", eh(common.SampleHandler)).Methods(http.MethodGet)
	api.HandleFunc("/taglist", eh(common.TagList)).Methods(http.MethodGet)
	return
}

/*
	各パス先でおきたエラーやパニックを検知
*/
func errorRoutingDetected(handler func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				msg := fmt.Sprintf("%s", rec)
				logger.Error(
					logger.GetApplicationError(errors.New(msg)).
						AddMessage("Faild to panic."),
				)
			}
		}()
		if err := handler(w, r); err != nil {
			logger.Error(
				logger.GetApplicationError(err),
			)
		}
	}
}
