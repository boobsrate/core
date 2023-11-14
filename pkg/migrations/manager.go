package migrations

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // migrate dest
	_ "github.com/golang-migrate/migrate/v4/source/file"       // migrate source
)

const (
	defaultTimeout  = time.Minute * 2
	defaultInterval = time.Second * 5

	defaultDirectory = "migrations/"
)

// Manager manages migrations.
type Manager struct {
	dir string
	dsn string
}

// Migrate migrates the database to the latest version.
// Migrate only migrates forward: if file version <= current version, Migrate does nothing.
func (m *Manager) Migrate(ctx context.Context) error {
	version, err := m.currentVersion()
	if err != nil {
		return fmt.Errorf("get curent version: %v", err)
	}

	mig, err := m.new()
	if err != nil {
		return err
	}
	defer func() {
		cerr := m.close(mig)
		if err == nil {
			err = cerr
		}
	}()

	if ok, _ := m.check(mig, version); ok {
		return nil
	}

	return m.migrateVersion(ctx, mig, version)
}

// MigrateVersion migrates the database to the specified version.
func (m *Manager) MigrateVersion(ctx context.Context, version uint) (err error) {
	mig, err := m.new()
	if err != nil {
		return err
	}
	defer func() {
		cerr := m.close(mig)
		if err == nil {
			err = cerr
		}
	}()

	return m.migrateVersion(ctx, mig, version)
}

// Wait waits for when database migration version matches the latest version.
func (m *Manager) Wait(ctx context.Context) error {
	version, err := m.currentVersion()
	if err != nil {
		return fmt.Errorf("get current version: %v", err)
	}

	return m.WaitVersion(ctx, version)
}

// WaitVersion waits for when database migration version matches the specified version.
func (m *Manager) WaitVersion(ctx context.Context, version uint) (err error) {
	mig, err := m.new()
	if err != nil {
		return err
	}
	defer func() {
		cerr := m.close(mig)
		if err == nil {
			err = cerr
		}
	}()

	return m.waitVersion(ctx, mig, version)
}

// new returns a new migrate instance.
func (m *Manager) new() (*migrate.Migrate, error) {
	mig, err := migrate.New("file://"+m.dir, m.dsn)
	if err != nil {
		return nil, fmt.Errorf("create migrate instance: %v", err)
	}
	return mig, nil
}

// close closes the migrate instance.
func (m *Manager) close(mig *migrate.Migrate) error {
	var msgs []string
	serr, derr := mig.Close()
	if serr != nil {
		msgs = append(msgs, fmt.Sprintf("close source: %v", serr))
	}
	if derr != nil {
		msgs = append(msgs, fmt.Sprintf("close destination: %v", derr))
	}
	if len(msgs) > 0 {
		return errors.New(strings.Join(msgs, ", "))
	}
	return nil
}

func (m *Manager) currentVersion() (uint, error) {
	files, err := filepath.Glob(m.dir + "*.sql")
	if err != nil {
		return 0, err
	}
	if len(files) == 0 {
		return 0, fmt.Errorf("found 0 migration files")
	}

	file := files[len(files)-1]
	file = strings.TrimPrefix(file, m.dir)
	parts := strings.Split(file, "_")
	if len(parts) < 2 {
		return 0, fmt.Errorf("bad migration file name")
	}

	v, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse version from file name: %v", err)
	}

	return uint(v), nil
}

// check returns true if the check was successful.
func (m *Manager) check(mig *migrate.Migrate, targetVersion uint) (bool, error) {
	currentVersion, dirty, err := mig.Version()
	if err == migrate.ErrNilVersion {
		// "no version" never meets targetVersion requirement
		return false, fmt.Errorf("nil current version in database")
	}
	if err != nil {
		return false, fmt.Errorf("check version error: %v", err)
	}
	if dirty {
		return false, errors.New("dirty database version, either migration is in progress or something went wrong and manual fix required")
	}
	if currentVersion < targetVersion {
		return false, fmt.Errorf("database version %v does not match desired minimal version %v yet", currentVersion, targetVersion)
	}

	return true, nil
}

func (m *Manager) migrateVersion(_ context.Context, mig *migrate.Migrate, version uint) error {
	if err := mig.Migrate(version); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("migrate: %v", err)
	}
	return nil
}

func (m *Manager) waitVersion(ctx context.Context, mig *migrate.Migrate, version uint) error {
	// Do the first check immediately.

	if ok, _ := m.check(mig, version); ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	t := time.NewTicker(defaultInterval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			if ok, _ := m.check(mig, version); ok {
				return nil
			}
		}
	}
}

// ManagerOption is an option for manager configuration.
type ManagerOption func(*Manager)

// WithDirectory returns an option that sets directory with migrations.
func WithDirectory(dir string) ManagerOption {
	return func(manager *Manager) {
		manager.dir = strings.TrimSuffix(filepath.Clean(dir), "/") + "/"
	}
}

// NewManager NewManger returns a new migrations manager.
func NewManager(dsn string, opts ...ManagerOption) (*Manager, error) {
	man := &Manager{
		dsn: dsn,
		dir: defaultDirectory,
	}

	for _, opt := range opts {
		opt(man)
	}

	return man, nil
}
