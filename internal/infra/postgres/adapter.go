package postgres

import (
	"Yaratam/internal/domain"
	"Yaratam/internal/infra/postgres/models"
	"database/sql"
	_ "database/sql"
	"errors"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
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

func (a *adapter) wrapError(err error, ErrorInfo string) error {
	switch err {
	case sql.ErrNoRows:
		return domain.ErrNotFound
	default:
		a.logger.WithError(err).Error(ErrorInfo)
		return domain.ErrInternalDatabase

	}
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

func (a *adapter) GetTime() (time.Time, error) {
	var tim time.Time
	err := a.db.Get(&tim, `SELECT now-nows FROM test_timne`)
	if err != nil {
		return time.Time{}, err
	}
	return tim, nil
}

func (a *adapter) GetUser(chatID int) (*domain.User, error) {
	var user models.User
	if err := a.db.Get(&user, `SELECT id, username, chat_id FROM users 
                             WHERE chat_id=$1`, chatID); err != nil {
		return &domain.User{}, a.wrapError(err, "Error while getting user")
	}
	return user.ToDomain(), nil
}
