package book

type DoubanBookResult struct {
	Code  uint32       `json:"code"`
	Msg   string       `json:"msg"`
	Books []DoubanBook `json:"books"`
}

type DoubanBook struct {
	ID          string   `json:"id"`
	Author      []string `json:"author"`
	AuthorIntro string   `json:"author_intro"`
	Translators []string `json:"translators"`
	Images      Image    `json:"images"`
	Binding     string   `json:"binding"`
	Category    string   `json:"category"`
	Rating      Rating   `json:"rating"`
	ISBN13      string   `json:"isbn13"`
	Pages       string   `json:"pages"`
	Price       string   `json:"price"`
	Pubdate     string   `json:"pubdate"`
	Publisher   string   `json:"publisher"`
	Producer    string   `json:"producer"`
	Serials     string   `json:"serials"`
	Subtitle    string   `json:"subtitle"`
	Summary     string   `json:"summary"`
	Title       string   `json:"title"`
	Tags        []Tag    `json:"tags"`
	Origin      string   `json:"origin"`
}

type Image struct {
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
}

type Tag struct {
	Name string `json:"name"`
}

type Rating struct {
	Average float32 `json:"average"`
}
