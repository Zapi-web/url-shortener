package server

type saveRequest struct {
	LongURL string `json:"url"`
	UserID  uint64 `json:"user_id"`
}

type saveResponse struct {
	ShortURL string `json:"short_url"`
}
