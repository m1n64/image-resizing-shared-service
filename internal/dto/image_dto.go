package dto

type ImageWithThumbnails struct {
	ID            string           `json:"id"`
	OriginalUrl   string           `json:"original_url"`
	CompressedUrl string           `json:"compressed_url"`
	Size          int64            `json:"size"`
	MIME          string           `json:"mime"`
	Status        string           `json:"status"`
	ErrorMessage  *string          `json:"error_message,omitempty"`
	Thumbnails    []ThumbnailShort `json:"thumbnails"`
}

type ThumbnailShort struct {
	Size string `json:"size"`
	Url  string `json:"url"`
	Type string `json:"type"`
}
