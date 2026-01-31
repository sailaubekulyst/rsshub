package ports

import "rsshub/internal/domain"

type ArticleRepository interface {
	AddNewArticle(NewArticle domain.Article) error
	GetAllArticles() ([]domain.Article, error)
}
