package movie

import (
	"html"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type parser struct {
	reID            *regexp.Regexp
	reBackground    *regexp.Regexp
	reSID           *regexp.Regexp
	reCat           *regexp.Regexp
	reYear          *regexp.Regexp
	reDirector      *regexp.Regexp
	reWriter        *regexp.Regexp
	reActor         *regexp.Regexp
	reGenre         *regexp.Regexp
	reCountry       *regexp.Regexp
	reLanguage      *regexp.Regexp
	reDuration      *regexp.Regexp
	reEpisodes      *regexp.Regexp
	reScreen        *regexp.Regexp
	reSubname       *regexp.Regexp
	reIMDB          *regexp.Regexp
	reSite          *regexp.Regexp
	reStripTags     *regexp.Regexp
	reNameMatch     *regexp.Regexp
	reRole          *regexp.Regexp
	reLifeDate      *regexp.Regexp
	reGender        *regexp.Regexp
	reConstellation *regexp.Regexp
	reBirthdate     *regexp.Regexp
	reBirthplace    *regexp.Regexp
	reCareer        *regexp.Regexp
	reNickname      *regexp.Regexp
	reFamily        *regexp.Regexp
	reCelebrityIMDB *regexp.Regexp
}

func newParser() *parser {
	return &parser{
		reID:            regexp.MustCompile(`/([0-9]+?)/`),
		reBackground:    regexp.MustCompile(`url\((.+?)\)`),
		reSID:           regexp.MustCompile(`sid: ([0-9]+?),`),
		reCat:           regexp.MustCompile(`\[(.+?)\]`),
		reYear:          regexp.MustCompile(`\(([0-9]+?)\)`),
		reDirector:      regexp.MustCompile(`导演\s*:\s*(.+?)(?:\n|$)`),
		reWriter:        regexp.MustCompile(`编剧\s*:\s*(.+?)(?:\n|$)`),
		reActor:         regexp.MustCompile(`主演\s*:\s*(.+?)(?:\n|$)`),
		reGenre:         regexp.MustCompile(`类型\s*:\s*(.+?)(?:\n|$)`),
		reCountry:       regexp.MustCompile(`制片国家/地区\s*:\s*(.+?)(?:\n|$)`),
		reLanguage:      regexp.MustCompile(`语言\s*:\s*(.+?)(?:\n|$)`),
		reDuration:      regexp.MustCompile(`片长\s*:\s*(.+?)(?:\n|$)`),
		reEpisodes:      regexp.MustCompile(`集数\s*:\s*(.+?)(?:\n|$)`),
		reScreen:        regexp.MustCompile(`上映日期\s*:\s*(.+?)(?:\n|$)`),
		reSubname:       regexp.MustCompile(`又名\s*:\s*(.+?)(?:\n|$)`),
		reIMDB:          regexp.MustCompile(`IMDb\s*:\s*(.+?)(?:\n|$)`),
		reSite:          regexp.MustCompile(`官方网站\s*:\s*(.+?)(?:\n|$)`),
		reStripTags:     regexp.MustCompile(`(?s)<[^>]*>`),
		reNameMatch:     regexp.MustCompile(`(.+第\w季|[\p{Han}\w：！，·]+)\s*(.*)`),
		reRole:          regexp.MustCompile(`\([饰配] (.+?)\)`),
		reLifeDate:      regexp.MustCompile(`生卒日期:\s*\n?\s*(.+?)\s*至`),
		reGender:        regexp.MustCompile(`性别:\s*\n?\s*(.+?)\n`),
		reConstellation: regexp.MustCompile(`星座:\s*\n?\s*(.+?)\n`),
		reBirthdate:     regexp.MustCompile(`出生日期:\s*\n?\s*(.+?)\n`),
		reBirthplace:    regexp.MustCompile(`出生地:\s*\n?\s*(.+?)\n`),
		reCareer:        regexp.MustCompile(`职业:\s*\n?\s*(.+?)\n`),
		reNickname:      regexp.MustCompile(`更多外文名:\s*\n?\s*(.+?)\n`),
		reFamily:        regexp.MustCompile(`家庭成员:\s*\n?\s*(.+?)\n`),
		reCelebrityIMDB: regexp.MustCompile(`imdb编号:\s*\n?\s*(.+?)\n`),
	}
}

func (p *parser) parseMovies(doc *goquery.Document, limit int, imageSize string) []Movie {
	movies := make([]Movie, 0)
	doc.Find("div.result-list").First().Find(".result").Each(func(_ int, s *goquery.Selection) {
		rating := strings.TrimSpace(s.Find("div.rating-info>.rating_nums").Text())
		if rating == "" {
			rating = "0"
		}

		onclick, _ := s.Find("div.title a").Attr("onclick")
		img := p.getImgBySize(attrOrEmpty(s.Find("a.nbg>img"), "src"), imageSize)
		sid := p.captureGroup(p.reSID, onclick)
		name := strings.TrimSpace(s.Find("div.title a").Text())
		titleMark := strings.TrimSpace(s.Find("div.title>h3>span").Text())
		cat := p.captureGroup(p.reCat, titleMark)
		subject := strings.TrimSpace(s.Find("div.rating-info>.subject-cast").Text())
		year := p.parseYear(subject)

		if cat == "电影" || cat == "电视剧" {
			movies = append(movies, Movie{
				Cat:    cat,
				SID:    sid,
				Name:   name,
				Rating: rating,
				Img:    img,
				Year:   year,
			})
		}
	})

	if limit > 0 && len(movies) > limit {
		return movies[:limit]
	}
	return movies
}

func (p *parser) parseMovieInfo(doc *goquery.Document, sid, imageSize string) MovieInfo {
	content := doc.Find("#content")

	nameStr := strings.TrimSpace(content.Find("h1>span:first-child").Text())
	name := nameStr
	originalName := ""
	if m := p.reNameMatch.FindStringSubmatch(nameStr); len(m) >= 3 {
		name = strings.TrimSpace(m[1])
		originalName = strings.TrimSpace(m[2])
	}

	year := p.captureGroup(p.reYear, strings.TrimSpace(content.Find("h1>span.year").Text()))
	rating := strings.TrimSpace(content.Find("div.rating_self strong.rating_num").Text())
	if rating == "" {
		rating = "0"
	}
	img := p.getImgBySize(attrOrEmpty(content.Find("a.nbgnbg>img"), "src"), imageSize)
	intro := strings.TrimSpace(strings.ReplaceAll(content.Find("div.indent>span").Text(), "©豆瓣", ""))
	infoText := p.parseInfoText(content.Find("#info"))

	director := p.captureGroup(p.reDirector, infoText)
	writer := p.captureGroup(p.reWriter, infoText)
	actor := p.captureGroup(p.reActor, infoText)
	genre := p.captureGroup(p.reGenre, infoText)
	site := p.captureGroup(p.reSite, infoText)
	country := p.captureGroup(p.reCountry, infoText)
	language := p.captureGroup(p.reLanguage, infoText)
	screen := p.captureGroup(p.reScreen, infoText)
	duration := p.captureGroup(p.reDuration, infoText)
	episodes := p.captureGroup(p.reEpisodes, infoText)
	subname := p.captureGroup(p.reSubname, infoText)
	imdb := p.captureGroup(p.reIMDB, infoText)

	celebrities := make([]Celebrity, 0)
	first := content.Find("#celebrities li.celebrity").First()
	if first.Length() > 0 {
		id := p.captureGroup(p.reID, attrOrEmpty(first.Find("div.info a.name"), "href"))
		imgStr := attrOrEmpty(first.Find("div.avatar"), "style")
		imgURL := p.getImgBySize(p.captureGroup(p.reBackground, imgStr), imageSize)
		name := strings.TrimSpace(first.Find("div.info a.name").Text())
		role := strings.TrimSpace(first.Find("div.info span.role").Text())
		celebrities = append(celebrities, Celebrity{
			ID:       id,
			Img:      imgURL,
			Name:     name,
			RoleType: "",
			Role:     role,
		})
	}

	return MovieInfo{
		SID:          sid,
		Name:         name,
		OriginalName: originalName,
		Rating:       rating,
		Img:          img,
		Year:         year,
		Intro:        intro,
		Director:     director,
		Writer:       writer,
		Actor:        actor,
		Genre:        genre,
		Site:         site,
		Country:      country,
		Language:     language,
		Screen:       screen,
		Duration:     duration,
		Episodes:     episodes,
		Subname:      subname,
		IMDB:         imdb,
		Celebrities:  celebrities,
	}
}

func (p *parser) parseInfoText(info *goquery.Selection) string {
	if info.Length() == 0 {
		return ""
	}

	rawHTML, err := info.Html()
	if err != nil {
		return strings.TrimSpace(info.Text())
	}

	text := strings.NewReplacer(
		"<br/>", "\n",
		"<br />", "\n",
		"<br>", "\n",
	).Replace(rawHTML)
	text = p.reStripTags.ReplaceAllString(text, "")
	text = html.UnescapeString(text)

	lines := strings.Split(text, "\n")
	clean := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.Join(strings.Fields(strings.TrimSpace(line)), " ")
		if line != "" {
			clean = append(clean, line)
		}
	}
	return strings.Join(clean, "\n")
}

func (p *parser) parseCelebrities(doc *goquery.Document) []Celebrity {
	result := make([]Celebrity, 0)
	doc.Find("#content ul.celebrities-list li.celebrity").Each(func(_ int, s *goquery.Selection) {
		id := p.captureGroup(p.reID, attrOrEmpty(s.Find("div.info a.name"), "href"))
		img := p.captureGroup(p.reBackground, attrOrEmpty(s.Find("div.avatar"), "style"))
		nameText := strings.TrimSpace(s.Find("div.info a.name").Text())
		name := nameText
		if parts := strings.Fields(nameText); len(parts) > 0 {
			name = parts[0]
		}
		rawRole := strings.TrimSpace(s.Find("div.info span.role").Text())
		roleType := ""
		if parts := strings.Fields(rawRole); len(parts) > 0 {
			roleType = parts[0]
		}
		role := p.captureGroup(p.reRole, rawRole)
		if role == "" {
			role = roleType
		}

		if roleType == "导演" || roleType == "配音" || roleType == "演员" {
			result = append(result, Celebrity{
				ID:       id,
				Img:      img,
				Name:     name,
				RoleType: roleType,
				Role:     role,
			})
		}
	})

	if len(result) > 15 {
		return result[:15]
	}
	return result
}

func (p *parser) parseCelebrityInfo(doc *goquery.Document, id string) CelebrityInfo {
	content := doc.Find("#content")
	img := attrOrEmpty(content.Find("#headline .nbg img"), "src")
	name := strings.TrimSpace(content.Find("h1").Text())
	intro := strings.TrimSpace(content.Find("#intro span.all").Text())
	if intro == "" {
		intro = strings.TrimSpace(content.Find("#intro div.bd").Text())
	}

	infoText := content.Find("div.info").Text()
	gender := strings.TrimSpace(p.captureGroup(p.reGender, infoText))
	constellation := strings.TrimSpace(p.captureGroup(p.reConstellation, infoText))
	birthdate := strings.TrimSpace(p.captureGroup(p.reBirthdate, infoText))
	if birthdate == "" {
		birthdate = strings.TrimSpace(p.captureGroup(p.reLifeDate, infoText))
	}
	birthplace := strings.TrimSpace(p.captureGroup(p.reBirthplace, infoText))
	role := strings.TrimSpace(p.captureGroup(p.reCareer, infoText))
	nickname := strings.TrimSpace(p.captureGroup(p.reNickname, infoText))
	family := strings.TrimSpace(p.captureGroup(p.reFamily, infoText))
	imdb := strings.TrimSpace(p.captureGroup(p.reCelebrityIMDB, infoText))

	return CelebrityInfo{
		ID:            id,
		Img:           img,
		Name:          name,
		Role:          role,
		Intro:         intro,
		Gender:        gender,
		Constellation: constellation,
		Birthdate:     birthdate,
		Birthplace:    birthplace,
		Nickname:      nickname,
		IMDB:          imdb,
		Family:        family,
	}
}

func (p *parser) parseWallpapers(doc *goquery.Document) []Photo {
	photos := make([]Photo, 0)
	doc.Find(".poster-col3>li").Each(func(_ int, s *goquery.Selection) {
		id := strings.TrimSpace(attrOrEmpty(s, "data-id"))
		if id == "" {
			return
		}
		small := "https://img2.doubanio.com/view/photo/s/public/p" + id + ".jpg"
		medium := "https://img2.doubanio.com/view/photo/m/public/p" + id + ".jpg"
		large := "https://img2.doubanio.com/view/photo/l/public/p" + id + ".jpg"
		size := strings.TrimSpace(s.Find("div.prop").Text())
		width := ""
		height := ""
		if size != "" {
			arr := strings.SplitN(size, "x", 2)
			if len(arr) == 2 {
				width = strings.TrimSpace(arr[0])
				height = strings.TrimSpace(arr[1])
			}
		}
		photos = append(photos, Photo{
			ID:     id,
			Small:  small,
			Medium: medium,
			Large:  large,
			Size:   size,
			Width:  width,
			Height: height,
		})
	})
	return photos
}

func (p *parser) parseYear(text string) string {
	parts := strings.Split(text, "/")
	if len(parts) == 0 {
		return ""
	}
	return strings.TrimSpace(parts[len(parts)-1])
}

func (p *parser) captureGroup(re *regexp.Regexp, text string) string {
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(strings.Trim(m[1], "\"'"))
}

func (p *parser) getImgBySize(url, imageSize string) string {
	if imageSize == "m" || imageSize == "l" {
		return strings.Replace(url, "s_ratio_poster", imageSize, 1)
	}
	return url
}

func attrOrEmpty(sel *goquery.Selection, key string) string {
	v, ok := sel.Attr(key)
	if !ok {
		return ""
	}
	return strings.TrimSpace(v)
}
