package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {
	lox := "https://minio-node-test.sovcombank.ru/credit-history/users/341/newHalva_70_2.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=XBNQZAP02C05O6EHQ9J6%2F20220425%2F%2Fs3%2Faws4_request&X-Amz-Date=20220425T134652Z&X-Amz-Expires=432000&X-Amz-SignedHeaders=host&X-Amz-Signature=4f0287eab903c9b9c9ab1ff31c9e4ec544f44418a2ed7fdb4ea7af7f4f9c48d9"
	u, err := url.Parse(lox)
	if err != nil {
		panic(err)
	}

	//fmt.Println(u.Path)
	path := strings.SplitAfter(u.Path, "/")
	fmt.Println(path[len(path)-1])
	/*config, err := configs.Parse()
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
	*/
}
