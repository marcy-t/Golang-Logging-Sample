package infrastructure

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Golang-Logging-Sample/pkg/logger"

	"github.com/go-chi/chi/v5"
	// "github.com/gorilla/mux"
)

// Router
func NewRouter(controller *ControllHandler) (r *chi.Mux) {
	r = chi.NewRouter()

	// common := controller.Common

	r.Route("/api/insight/v1", func(r chi.Router) {

		r.Get("/ping", defaultRouter)

	})

	// root.PathPrefix("/api/insight/v1").Handler(http.FileServer(http.Dir("./public/")))
	// PathPrefix
	// api := root.PathPrefix("/api/insight/v1").Subrouter()
	// SamplePath
	// common := controller.Common
	// api.HandleFunc("/ping", eh(common.SampleHandler)).Methods(http.MethodGet, "OPTIONS")
	// api.HandleFunc("/taglist", eh(common.TagList)).Methods(http.MethodGet, "OPTIONS")
	return
}

func defaultRouter(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("title:%s", "hoghohohohohohohohohohohohohohohoh")))
	return
}

// // // Router
// func NewRouter(controller *ControllHandler) (root *mux.Router) {
// 	root = mux.NewRouter()
// 	// root = chi.Router()
// 	root.NotFoundHandler = http.HandlerFunc(h.NotFoundHandler)
// 	eh := errorRoutingDetected
// 	// PathPrefix
// 	api := root.PathPrefix("/api/insight/v1").Subrouter()
// 	// SamplePath
// 	common := controller.Common
// 	api.HandleFunc("/ping", eh(common.SampleHandler)).Methods(http.MethodGet, "OPTIONS")
// 	api.HandleFunc("/taglist", eh(common.TagList)).Methods(http.MethodGet, "OPTIONS")
// 	return
// }

// ミドルウェア
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request received")

		next.ServeHTTP(w, r)

		log.Println("Request handled")
	})
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
