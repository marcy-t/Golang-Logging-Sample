package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Golang-Logging-Sample/infrastructure"
	"github.com/Golang-Logging-Sample/pkg/logger"
)

func main() {
	if err := logger.NewLogger(); err != nil {
		logger.Fatal(
			logger.GetApplicationError(err).
				AddMessage("An error has occured"),
		)
	}

	// db, err := rdb.NewHandler()
	// if err != nil {
	// 	logger.Fatal(
	// 		logger.GetApplicationError(err).
	// 			AddMessage("Faild to get database connection."),
	// 	)
	// }

	// handler := infrastructure.NewServer(db)

	handler := infrastructure.NewServer(nil)

	// Port Conf
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("0.0.0.0:%s", port),
		WriteTimeout: 180 * time.Second,
		ReadTimeout:  180 * time.Second,
		IdleTimeout:  300 * time.Second,
	}

	tag := logger.NewTag("host", srv.Addr)
	logger.Info("00", "Starting Serverer Listening", tag)

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(
			logger.GetApplicationError(err).
				AddMessage("Server Down..."),
		)
	}

}
