package postgres

import (
	"Yaratam/internal/domain"
	"errors"
	"github.com/dlmiddlecote/sqlstats"
	"github.com/golang-migrate/migrate/v4"
	postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type adapter struct {
	logger logrus.FieldLogger
	config *Config
	db     *sqlx.DB
}

func NewAdapter(logger logrus.FieldLogger, config *Config) (domain.Database, error) {
	a := &adapter{
		logger: logger,
		config: config,
	}

	db, err := sqlx.Open("pgx", config.ConnectionString())
	if err != nil {
		logger.Errorf("cannot open sql connection: %w", err)
		return nil, err
	}
	a.db = db

	// Create a new collector, the name will be used as a label on the metrics
	collector := sqlstats.NewStatsCollector("credit_history_db", db)

	// Register it with Prometheus
	prometheus.MustRegister(collector)

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifeTime)

	// Migrations block
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(config.MigrationsSourceURL, config.Name, driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return a, nil
}

func (a *adapter) SaveAppLogs(userID int, header string, body string, status int) error {
	if _, err := a.db.Exec(
		`INSERT INTO app_log(user_id, start_dt, header, body, status)
					   VALUES ($1, now(), $2, $3, $4)`,
		userID,
		header,
		body,
		status); err != nil {
		a.logger.WithError(err).Error("Error while saving info in app_log!")
		return domain.ErrInternalDatabase
	}

	return nil
}

/*
func (a *adapter) SelectShemaMigration()(int ,error)  {
	var
}*/
