package sqlite

import (
	"database/sql"

	"rsshub/internal/ports"
)

type WorkerAndIntervalRepositorySqlite struct {
	db *sql.DB
}

func GetWorkerAndIntervalRepositorySqlite(db *sql.DB) ports.WorkerAndIntervalRepository {
	return &WorkerAndIntervalRepositorySqlite{
		db: db,
	}
}

func (r *WorkerAndIntervalRepositorySqlite) SetStarted() error {
	query := `
		UPDATE workerandinterval
		SET isstarted = 'yes'
		WHERE id = 1
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *WorkerAndIntervalRepositorySqlite) CloseStarted() error {
	query := `
		UPDATE workerandinterval
		SET isstarted = 'no'
		WHERE id = 1
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *WorkerAndIntervalRepositorySqlite) SetWorkerNs(numWorkers int) error {
	query := `
		UPDATE workerandinterval
		SET workerns = ?
		WHERE id = 1
	`
	_, err := r.db.Exec(query, numWorkers)
	return err
}

func (r *WorkerAndIntervalRepositorySqlite) SetInterval(newInterval string) error {
	query := `
		UPDATE workerandinterval
		SET interval = ?
		WHERE id = 1
	`
	_, err := r.db.Exec(query, newInterval)
	return err
}

func (r *WorkerAndIntervalRepositorySqlite) GetInterval() (string, error) {
	query := `
		SELECT interval
		FROM workerandinterval
		WHERE id = 1
	`

	var interval string
	err := r.db.QueryRow(query).Scan(&interval)
	return interval, err
}

func (r *WorkerAndIntervalRepositorySqlite) GetStarted() (string, error) {
	query := `
		SELECT isstarted
		FROM workerandinterval
		WHERE id = 1
	`

	var started string
	err := r.db.QueryRow(query).Scan(&started)
	return started, err
}

func (r *WorkerAndIntervalRepositorySqlite) GetWorkerNs() (int, error) {
	query := `
		SELECT workerns
		FROM workerandinterval
		WHERE id = 1
	`

	var workers int
	err := r.db.QueryRow(query).Scan(&workers)
	return workers, err
}
