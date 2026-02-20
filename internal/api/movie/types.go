package movie

type Movie struct {
	Cat    string `json:"cat"`
	SID    string `json:"sid"`
	Name   string `json:"name"`
	Rating string `json:"rating"`
	Img    string `json:"img"`
	Year   string `json:"year"`
}

type MovieInfo struct {
	SID          string      `json:"sid"`
	Name         string      `json:"name"`
	OriginalName string      `json:"originalName"`
	Rating       string      `json:"rating"`
	Img          string      `json:"img"`
	Year         string      `json:"year"`
	Intro        string      `json:"intro"`
	Director     string      `json:"director"`
	Writer       string      `json:"writer"`
	Actor        string      `json:"actor"`
	Genre        string      `json:"genre"`
	Site         string      `json:"site"`
	Country      string      `json:"country"`
	Language     string      `json:"language"`
	Screen       string      `json:"screen"`
	Duration     string      `json:"duration"`
	Episodes     string      `json:"episodes"`
	Subname      string      `json:"subname"`
	IMDB         string      `json:"imdb"`
	Celebrities  []Celebrity `json:"celebrities"`
}

type Celebrity struct {
	ID       string `json:"id"`
	Img      string `json:"img"`
	Name     string `json:"name"`
	RoleType string `json:"-"`
	Role     string `json:"role"`
}

type CelebrityInfo struct {
	ID            string `json:"id"`
	Img           string `json:"img"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	Intro         string `json:"intro"`
	Gender        string `json:"gender"`
	Constellation string `json:"constellation"`
	Birthdate     string `json:"birthdate"`
	Birthplace    string `json:"birthplace"`
	Nickname      string `json:"nickname"`
	IMDB          string `json:"imdb"`
	Family        string `json:"family"`
}

type Photo struct {
	ID     string `json:"id"`
	Small  string `json:"small"`
	Medium string `json:"medium"`
	Large  string `json:"large"`
	Size   string `json:"size"`
	Width  string `json:"width"`
	Height string `json:"height"`
}
