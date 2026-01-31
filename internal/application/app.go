package application

import (
	"database/sql"

	"rsshub/internal/adapters/cli"
	"rsshub/internal/adapters/sqlite"
	"rsshub/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	cli *cli.Cli
}

func GetNewApp(args []string) (*Application, error) {
	var app Application
	db, err := sql.Open("sqlite3", "rsshub.db")
	if err != nil {
		return nil, err
	}
	err = runMigrations(db)
	if err != nil {
		return &app, err
	}
	NameAndUrlRepo := sqlite.GetNameAndUrlRepositorySqlite(db)
	ArticleRepo := sqlite.GetArticleRepositorySqlite(db)
	WorkerAndIntervalRepo := sqlite.GetWorkerAndIntervalRepositorySqlite(db)
	service, _ := service.GetService(NameAndUrlRepo, ArticleRepo, WorkerAndIntervalRepo)
	app.cli = cli.GetCli(service, args)
	return &app, nil
}

func (app *Application) Run() error {
	return app.cli.Run()
}

func runMigrations(db *sql.DB) error {
	queries := []string{
		// 1. nameandurls
		`
		CREATE TABLE IF NOT EXISTS nameandurls (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			url  TEXT NOT NULL UNIQUE
		);
		`,

		// 2. articles
		`
		CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			link TEXT,
			description TEXT,
			feed_id INTEGER
		);
		`,

		// 3. workerandinterval
		`
		CREATE TABLE IF NOT EXISTS workerandinterval (
			id INTEGER PRIMARY KEY,
			interval TEXT NOT NULL,
			workerns INTEGER NOT NULL,
			isstarted TEXT NOT NULL
		);

		`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM workerandinterval WHERE id = ?`, 1).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = db.Exec(`INSERT OR IGNORE INTO workerandinterval (id, interval, workerns, isstarted) VALUES (?, ?, ?, ?)`, "1", "3m", 5, "no")
		if err != nil {
			return err
		}

	}
	return nil
}
