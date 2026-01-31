package ports

import "rsshub/internal/domain"

type NameAndUrlRepository interface {
	AddNewNameAndUrl(NewUrl domain.NameAndUrl) error
	DeleteUrlByName(Name string) error
	GetAllUrls() ([]domain.NameAndUrl, error)
	GetIDByUrl(url string) (int, error)
}
