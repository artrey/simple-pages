package storage

import "github.com/artrey/simple-pages/pkg/models"

type Interface interface {
	GetPages() ([]*models.Page, error)
	GetPageById(id int64) (*models.Page, error)
	CreatePage(title, imageUri, text string) (*models.Page, error)
	UpdatePageById(id int64, title, imageUri, text string) (*models.Page, error)
	DeletePageById(id int64) error
}
