package media

import "encoding/json"

type HotMediaResponse struct {
	Category      string           `json:"category"`
	Type          string           `json:"type"`
	Total         int              `json:"total"`
	Tags          []HotMediaTag    `json:"tags,omitempty"`
	RecommendTags []HotMediaOption `json:"recommend_tags,omitempty"`
	Items         []HotMediaItem   `json:"items"`
}

type HotMediaTag struct {
	Category string           `json:"category"`
	Selected bool             `json:"selected"`
	Title    string           `json:"title"`
	Types    []HotMediaOption `json:"types,omitempty"`
}

type HotMediaOption struct {
	Selected bool   `json:"selected"`
	Type     string `json:"type"`
	Title    string `json:"title"`
}

type HotMediaItem struct {
	ID               string          `json:"id"`
	Type             string          `json:"type"`
	Title            string          `json:"title"`
	Year             string          `json:"year,omitempty"`
	URI              string          `json:"uri,omitempty"`
	URL              string          `json:"url,omitempty"`
	CardSubtitle     string          `json:"card_subtitle,omitempty"`
	Rating           *HotMediaRating `json:"rating,omitempty"`
	Pic              *HotMediaPic    `json:"pic,omitempty"`
	Playable         *bool           `json:"playable,omitempty"`
	NullRatingReason string          `json:"null_rating_reason,omitempty"`
	EpisodesInfo     string          `json:"episodes_info,omitempty"`
	HonorInfos       json.RawMessage `json:"honor_infos,omitempty"`
}

type HotMediaRating struct {
	Count     int     `json:"count"`
	Max       int     `json:"max"`
	StarCount float64 `json:"star_count"`
	Value     float64 `json:"value"`
}

type HotMediaPic struct {
	Large  string `json:"large"`
	Normal string `json:"normal"`
}
