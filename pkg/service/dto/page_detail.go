package dto

import "github.com/artrey/simple-pages/pkg/models"

type PageDetail struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	ImageUri  string `json:"imageUri"`
	Text      string `json:"text"`
	CreatedAt int64  `json:"createdAt"`
}

func PageDetailFromModelPage(p *models.Page) *PageDetail {
	return &PageDetail{
		Id:        p.Id,
		Title:     p.Title,
		ImageUri:  p.ImageUri,
		Text:      p.Text,
		CreatedAt: p.CreatedAt,
	}
}
