package main

import (
	"Yaratam/internal/configs"
	"Yaratam/internal/domain"
	"Yaratam/internal/infra/http"
	"Yaratam/internal/infra/postgres"
	"Yaratam/pkg/logging"
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	config, err := configs.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Printf("Invalid args: %v\n", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(config.Logger)
	if err != nil {
		panic(err)
	}

	// Init PostgreSQL
	db, err := postgres.NewAdapter(logger, config.Postgres)
	if err != nil {
		logger.WithError(err).Fatal("Error while creating a new database adapter!")
	}
	// Init service
	service := domain.NewService(logger, db)

	// Init HTTP adapter
	httpAdapter, err := http.NewAdapter(logger, config.HTTP, service)
	if err != nil {
		logger.WithError(err).Fatal("Error creating new HTTP adapter!")
	}
	fmt.Println(service.GetTimes())
	shutdown := make(chan error, 1)

	go func(shutdown chan<- error) {
		shutdown <- httpAdapter.ListenAndServe()
	}(shutdown)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		logger.WithField("signal", s).Info("Got the signal!")
	case err := <-shutdown:
		logger.WithError(err).Error("Error running the application!")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	logger.Info("Stopping application...")

	if err := httpAdapter.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Error shutting down the HTTP server!")
	}

	time.Sleep(time.Second)

	logger.Info("The application stopped.")

}
