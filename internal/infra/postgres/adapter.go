package postgres

import (
	"Yaratam/internal/domain"
	"Yaratam/internal/infra/postgres/models"
	"database/sql"
	_ "database/sql"
	"errors"
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

func (a *adapter) GetUser(chatID int) (*domain.User, error) {
	var user models.User
	if err := a.db.Get(&user, `SELECT id, username, chat_id, coalesce(path_id,-1) FROM users 
                             WHERE chat_id=$1`, chatID); err != nil {
		return &domain.User{}, a.wrapError(err, "Error while getting user")
	}
	return user.ToDomain(), nil
}

func (a *adapter) AddUser(user *domain.User) error {
	if _, err := a.db.Exec(`INSERT INTO users (username,chat_id) VALUES ($1,$2)`, user.UserName, user.ChatID); err != nil {
		a.logger.WithError(err).Error("Error while insert user")
		return domain.ErrInternalDatabase
	}
	return nil
}

func (a *adapter) AddPath(chatID int, path *domain.Path) (int, error) {
	var ID int
	if err := a.db.Get(&ID, `INSERT INTO path (display_name) VALUES ($1) RETURNING id`, path.DisplayName); err != nil {
		a.logger.WithError(err).Error("Error while insert path")
		return 0, domain.ErrInternalDatabase
	}

	if _, err := a.db.Exec(`INSERT INTO users_paths (path_id, user_id) SELECT $1,id from users where chat_id=$2 LIMIT 1`, ID, chatID); err != nil {
		a.logger.WithError(err).Error("Error while insert users_paths")
		return 0, domain.ErrInternalDatabase
	}
	return ID, nil
}

func (a *adapter) ChangeUserPath(chatID, pathID int) error {
	if _, err := a.db.Exec(`UPDATE users SET path_id=$1  where chat_id=$2 `, pathID, chatID); err != nil {
		a.logger.WithError(err).Error("Error while insert users_paths")
		return domain.ErrInternalDatabase
	}
	return nil
}

func (a *adapter) AddFile(chatID int, path string) error {
	user, err := a.GetUser(chatID)
	if err != nil {
		return err
	}
	if user.PathID == -1 {
		return domain.ErrInvalidInputData
	}
	if _, err := a.db.Exec(`INSERT INTO file ( user_id, path_id,paths) SELECT id, path_id, $2 from users where chat_id=$1 LIMIT 1`, chatID, path); err != nil {
		a.logger.WithError(err).Error("Error while insert file")
		return domain.ErrInternalDatabase
	}
	return nil
}

func (a *adapter) GetFiles(chatID int64) ([]*domain.File, error) {
	var file models.Files

	if err := a.db.Select(&file, `SELECT paths FROM file WHERE path_id in (
			SELECT path_id FROM users_paths up 
    			join users u on u.id = up.user_id 
                	        AND u.path_id=up.path_id 
               	WHERE u.chat_id=$1)`, chatID); err != nil {
		a.logger.WithError(err).Error(`Error while getting file`)
		return nil, domain.ErrInternalDatabase
	}
	return file.Domain(), nil
}

func (a *adapter) GetPaths(chatID int) ([]*domain.Path, error) {
	var paths models.Paths

	if err := a.db.Select(&paths, `SELECT id, display_name FROM path WHERE deleted=false AND  id in (
			SELECT path_id FROM users_paths up 
    			join users u on u.id = up.user_id 
                	        AND u.path_id=up.path_id 
               	WHERE u.chat_id=$1)`, chatID); err != nil {
		a.logger.WithError(err).Error(`Error while getting file`)
		return nil, domain.ErrInternalDatabase
	}
	return paths.Domain(), nil
}

func (a *adapter) DeletePaths(chatID int) error {
	if _, err := a.db.Exec(`UPDATE path SET deleted= true  where id = (SELECT path_id from users where chat_id=$2 LIMIT 1) `, chatID); err != nil {
		a.logger.WithError(err).Error("Error while deleted paths")
		return domain.ErrInternalDatabase
	}
	return nil
}
