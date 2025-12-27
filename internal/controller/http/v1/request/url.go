package request

type ShortenRequest struct {
	URL string `json:"url" validate:"required,url"`
}
