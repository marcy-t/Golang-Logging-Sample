package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/marcy-t/Golang-Logging-Sample/infrastructure"
	"github.com/marcy-t/Golang-Logging-Sample/infrastructure/rdb"
	"github.com/marcy-t/Golang-Logging-Sample/pkg/logger"
)

func main() {
	if err := logger.NewLogger(); err != nil {
		logger.Fatal(
			logger.GetApplicationError(err).
				AddMessage("An error has occured"),
		)
	}

	db, err := rdb.NewHandler()
	if err != nil {
		logger.Fatal(
			logger.GetApplicationError(err).
				AddMessage("Faild to get database connection."),
		)
	}

	handler := infrastructure.NewServer(db)

	// Port Conf
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	srv := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("127.0.0.1:%s", port),
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
