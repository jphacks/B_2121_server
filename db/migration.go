package db

import (
	"io/ioutil"
	"path"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"github.com/golang-migrate/migrate/source"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/jmoiron/sqlx"
	"golang.org/x/xerrors"
)

const migrationDir = "migrations"

func Migrate(db *sqlx.DB, baseDir string) error {
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return xerrors.Errorf("failed to create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path.Join(baseDir, migrationDir), "mysql", driver)
	if err != nil {
		return xerrors.Errorf("failed to create a new instance of migrate: %w", err)
	}
	latestVer, err := getLatestVersion(path.Join(baseDir, migrationDir))
	if err != nil {
		return xerrors.Errorf("failed to get the latest database schema version from migration files: %w", err)
	}

	currentVer, _, err := m.Version()

	if err != nil && !xerrors.Is(err, migrate.ErrNilVersion) {
		return xerrors.Errorf("failed to get the current database version: %w", err)
	}

	if currentVer == latestVer {
		return nil
	}

	if currentVer > latestVer {
		return xerrors.Errorf("unknown schema version %d is in the database. (latest supported version is %d)", currentVer, latestVer)
	}

	err = m.Up()
	if err != nil {
		return xerrors.Errorf("database migration failed: %w", err)
	}

	return nil
}

func getLatestVersion(dir string) (uint, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return 0, xerrors.Errorf("failed to get migration files: %w", err)
	}

	maxVersion := uint(0)
	for _, f := range files {
		m, err := source.Parse(f.Name())
		if err != nil {
			return 0, xerrors.Errorf("failed to parse name of migration files: %w", err)
		}
		if maxVersion < m.Version {
			maxVersion = m.Version
		}
	}

	return maxVersion, nil
}
