package main

import (
	"Yaratam/internal/configs"
	"Yaratam/internal/domain"
	"Yaratam/internal/infra/bot"
	"Yaratam/internal/infra/httpreq"
	"Yaratam/internal/infra/postgres"
	"Yaratam/pkg/logging"
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
	logger.Info("Start works ))")
	fmt.Println(config.Telegram.Token)
	// Init PostgreSQL
	db, err := postgres.NewAdapter(logger, config.Postgres)
	if err != nil {
		logger.WithError(err).Fatal("Error while creating a new database adapter!")
	}

	//Init HTTP req
	httpreq, err := httpreq.NewAdapter(logger, *config.HTTPReq)
	if err != nil {
		logger.WithError(err).Error("Error while creating a new httpreq adapter!")
	}
	// Init service
	service := domain.NewService(logger, db, httpreq)

	//Init Telegram
	telegramBotAdapter, err := bot.NewAdapter(config.Telegram, service, logger)
	if err != nil {
		logger.WithError(err).Fatal("Error creating new Telegram adapter!")
	}

	stop := make(chan error, 1)

	// Receive errors form start bot func into error channel
	go func(stop chan<- error) {
		stop <- telegramBotAdapter.StartBot()
	}(stop)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		logger.WithField("signal", s).Info("Got the signal!")
	case err := <-stop:
		logger.WithError(err).Error("Error running the application!")
	}

	logger.Info("Stopping application...")
	telegramBotAdapter.StopBot()

	time.Sleep(time.Second)

	logger.Info("The application stopped.")

}
