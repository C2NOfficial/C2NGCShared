package models

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type MediaType string

func (m MediaType) String() string {
	return string(m)
}

type MediaItem struct {
	URL  string    `firestore:"url" json:"url"`
	Type MediaType `firestore:"type" json:"type"`
}
