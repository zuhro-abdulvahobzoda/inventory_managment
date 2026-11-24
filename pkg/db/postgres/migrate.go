package postgres

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5" // registers the "pgx5" database driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func Migrate(cfg Config, source fs.FS, dir string) (err error) {
	src, err := iofs.New(source, dir)
	if err != nil {
		return fmt.Errorf("open migrations source: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, cfg.dsnWithScheme("pgx5"))
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if err == nil {
			err = errors.Join(srcErr, dbErr)
		}
	}()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
