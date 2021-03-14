package dto

import "github.com/artrey/simple-pages/pkg/models"

type PageInfo struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	ImageUri  string `json:"imageUri"`
	CreatedAt int64  `json:"createdAt"`
}

func PageInfoFromModelPage(p *models.Page) *PageInfo {
	return &PageInfo{
		Id:        p.Id,
		Title:     p.Title,
		ImageUri:  p.ImageUri,
		CreatedAt: p.CreatedAt,
	}
}
