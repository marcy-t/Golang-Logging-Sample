package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Golang-Logging-Sample/domain"
	"github.com/Golang-Logging-Sample/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Object struct {
	AppResp *domain.ApplicationResponse
}

// Router
func NewRouter(controller *ControllHandler) (r *chi.Mux) {
	r = chi.NewRouter()

	// chi middleware log
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/insight/v1", func(r chi.Router) {
		r.Use(authenticatedUsersOnly)

		obj := Object{}
		eh := obj.ErrorRoutingDetected

		r.With(obj.ResponseJson).Get("/ping", eh(controller.Common.SampleHandler))
	})

	return
}

// 認証済みユーザー確認
func authenticatedUsersOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("########################")
		log.Println("########################")
		log.Println("Before AuthenticatedUsersOnly")
		log.Println("########################")
		log.Println("########################")
		next.ServeHTTP(w, r)
		log.Println("########################")
		log.Println("########################")
		log.Println("After AuthenticatedUsersOnly")
		log.Println("########################")
		log.Println("########################")
	})
}

func (re *Object) ResponseJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		object, err := json.Marshal(re.AppResp.Message)
		if err != nil {
			logger.Error(
				logger.GetApplicationError(err).AddMessage("An error has occured JsonMarshal "),
			)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(object)))
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.WriteHeader(re.AppResp.StatusCode)
		w.Write(object)
	})
}

/*
各パス先でおきたエラーやパニックを検知
*/
func (re *Object) ErrorRoutingDetected(handler func(http.ResponseWriter, *http.Request) (*domain.ApplicationResponse, error)) func(http.ResponseWriter, *http.Request) {
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
		resp, errType := handler(w, r)
		if errType != nil {
			logger.Error(
				logger.GetApplicationError(errType).ErrorResponseJSON(w, r),
			)
		}

		re.AppResp = resp // debug

	}
}
