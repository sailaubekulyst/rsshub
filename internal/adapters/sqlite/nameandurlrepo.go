package sqlite

import (
	"database/sql"

	"rsshub/internal/domain"
	"rsshub/internal/ports"
)

type NameAndUrlRepositorySqlite struct {
	db *sql.DB
}

func GetNameAndUrlRepositorySqlite(db *sql.DB) ports.NameAndUrlRepository {
	return &NameAndUrlRepositorySqlite{
		db: db,
	}
}

func (r *NameAndUrlRepositorySqlite) AddNewNameAndUrl(newUrl domain.NameAndUrl) error {
	query := `
		INSERT INTO nameandurls (name, url)
		VALUES (?, ?)
	`

	_, err := r.db.Exec(query, newUrl.Name, newUrl.Url)
	return err
}

func (r *NameAndUrlRepositorySqlite) DeleteUrlByName(name string) error {
	query := `
		DELETE FROM nameandurls
		WHERE name = ?
	`

	_, err := r.db.Exec(query, name)
	return err
}

func (r *NameAndUrlRepositorySqlite) GetAllUrls() ([]domain.NameAndUrl, error) {
	query := `
		SELECT name, url
		FROM nameandurls
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.NameAndUrl

	for rows.Next() {
		var item domain.NameAndUrl
		err := rows.Scan(&item.Name, &item.Url)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *NameAndUrlRepositorySqlite) GetIDByUrl(url string) (int, error) {
	query := `
		SELECT id
		FROM nameandurls
		WHERE url = ?
	`

	var id int
	err := r.db.QueryRow(query, url).Scan(&id)
	if err != nil {
		return 0, err // может быть sql.ErrNoRows
	}

	return id, nil
}
