package dto

type CreateLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required,url"`
	Slug        string `json:"slug" binding:"omitempty,min=3,max=50"`
}

type LinkResponse struct {
	ID        int    `json:"id"`
	Original  string `json:"original_url"`
	Slug      string `json:"slug"`
	ShortURL  string `json:"short_url"`
	CreatedAt string `json:"created_at"`
}
