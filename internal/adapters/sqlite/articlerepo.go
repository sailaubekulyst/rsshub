package sqlite

import (
	"database/sql"

	"rsshub/internal/domain"
	"rsshub/internal/ports"
)

type ArticleRepositorySqlite struct {
	db *sql.DB
}

func GetArticleRepositorySqlite(db *sql.DB) ports.ArticleRepository {
	return &ArticleRepositorySqlite{
		db: db,
	}
}

func (r *ArticleRepositorySqlite) AddNewArticle(newArticle domain.Article) error {
	query := `
		INSERT INTO articles (title, link, description, feed_id)
		VALUES (?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		newArticle.Title,
		newArticle.Link,
		newArticle.Description,
		newArticle.Feedid,
	)

	return err
}

func (r *ArticleRepositorySqlite) GetAllArticles() ([]domain.Article, error) {
	query := `
		SELECT title, link, description, feed_id
		FROM articles
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []domain.Article

	for rows.Next() {
		var article domain.Article
		err := rows.Scan(
			&article.Title,
			&article.Link,
			&article.Description,
			&article.Feedid,
		)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
