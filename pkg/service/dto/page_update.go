package dto

type PageUpdate struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	ImageUri string `json:"imageUri"`
	Text     string `json:"text"`
}
