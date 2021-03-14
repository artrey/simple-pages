package memory

import (
	"errors"
	"github.com/artrey/simple-pages/pkg/models"
	"sync"
	"time"
)

type Storage struct {
	pages []*models.Page
	mu    sync.RWMutex
}

var (
	ErrNotFound = errors.New("not found")
)

func New() *Storage {
	return &Storage{
		pages: make([]*models.Page, 0),
	}
}

func (s *Storage) GetPages() ([]*models.Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.pages, nil
}

func (s *Storage) GetPageById(id int64) (*models.Page, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, page := range s.pages {
		if page.Id == id {
			return page, nil
		}
	}
	return nil, ErrNotFound
}

func (s *Storage) CreatePage(title, imageUri, text string) (*models.Page, error) {
	s.mu.RLock()
	id := int64(1)
	if len(s.pages) > 0 {
		id = s.pages[len(s.pages)-1].Id + 1
	}
	s.mu.RUnlock()

	page := &models.Page{
		Id:        id,
		Title:     title,
		ImageUri:  imageUri,
		CreatedAt: time.Now().Unix(),
		Text:      text,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.pages = append(s.pages, page)

	return page, nil
}

func (s *Storage) UpdatePageById(id int64, title, imageUri, text string) (*models.Page, error) {
	page, err := s.GetPageById(id)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	page.Title = title
	page.ImageUri = imageUri
	page.Text = text

	return page, nil
}

func (s *Storage) DeletePageById(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for idx, page := range s.pages {
		if page.Id == id {
			s.pages = append(s.pages[:idx], s.pages[idx+1:]...)
			return nil
		}
	}

	return ErrNotFound
}
