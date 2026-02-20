package book

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type parser struct {
	reID               *regexp.Regexp
	reInfoPair         *regexp.Regexp
	reRemoveSplitSpace *regexp.Regexp
}

func newParser() *parser {
	return &parser{
		reID:               regexp.MustCompile(`sid: ([0-9]+?),`),
		reInfoPair:         regexp.MustCompile(`([^\s]+?):\s*([^\n]+)`),
		reRemoveSplitSpace: regexp.MustCompile(`\s+?/\s+`),
	}
}

func (p *parser) parseSearchList(doc *goquery.Document, count int) []DoubanBook {
	books := make([]DoubanBook, 0)
	doc.Find("div.result-list").First().Find(".result").Each(func(_ int, s *goquery.Selection) {
		onclick, _ := s.Find("div.title a").Attr("onclick")
		title := strings.TrimSpace(s.Find("div.title a").Text())
		summary := strings.TrimSpace(s.Find("p").Text())
		large, _ := s.Find(".pic img").Attr("src")
		rate := strings.TrimSpace(s.Find(".rating_nums").Text())
		subStr := strings.TrimSpace(s.Find(".subject-cast").Text())

		author, publisher, pubdate := p.parseSubjectCast(subStr)
		id := captureGroup(p.reID, onclick)

		avg := float32(0)
		if rate != "" {
			if parsed, err := parseFloat32(rate); err == nil {
				avg = parsed
			}
		}

		books = append(books, DoubanBook{
			ID:          id,
			Author:      author,
			AuthorIntro: "",
			Translators: []string{},
			Images: Image{
				Large:  large,
				Medium: "",
				Small:  "",
			},
			Binding:   "",
			Category:  "",
			Rating:    Rating{Average: avg},
			ISBN13:    "",
			Pages:     "",
			Price:     "",
			Pubdate:   pubdate,
			Publisher: publisher,
			Producer:  "",
			Serials:   "",
			Subtitle:  "",
			Summary:   summary,
			Title:     title,
			Tags:      []Tag{},
			Origin:    "",
		})
	})

	if count > 0 && len(books) > count {
		return books[:count]
	}
	return books
}

func (p *parser) parseBookPage(doc *goquery.Document, id string) DoubanBook {
	wrapper := doc.Find("#wrapper")
	title := strings.TrimSpace(wrapper.Find("h1>span:first-child").Text())
	largeImg, _ := wrapper.Find("a.nbg").Attr("href")
	smallImg, _ := wrapper.Find("a.nbg>img").Attr("src")
	content := wrapper.Find("#content")

	tags := make([]Tag, 0)
	wrapper.Find("a.tag").Each(func(_ int, t *goquery.Selection) {
		tags = append(tags, Tag{Name: strings.TrimSpace(t.Text())})
	})

	ratingStr := strings.TrimSpace(content.Find("div.rating_self strong.rating_num").Text())
	rating := Rating{Average: 0}
	if ratingStr != "" {
		if avg, err := parseFloat32(ratingStr); err == nil {
			rating.Average = avg
		}
	}

	summary, _ := content.Find("#link-report .hidden .intro").First().Html()
	summary = strings.TrimSpace(summary)
	if summary == "" {
		summary, _ = content.Find("#link-report .intro").First().Html()
		summary = strings.TrimSpace(summary)
	}

	authorIntro, _ := content.Find(".related_info .indent:not([id]) > .all.hidden .intro").First().Html()
	authorIntro = strings.TrimSpace(authorIntro)
	if authorIntro == "" {
		authorIntro, _ = content.Find(".related_info .indent:not([id]) .intro").First().Html()
		authorIntro = strings.TrimSpace(authorIntro)
	}

	infoText := strings.TrimSpace(content.Find("#info").Text())
	infoMap := p.parseInfoText(infoText)

	author := p.getTexts(infoMap, "作者")
	translators := p.getTexts(infoMap, "译者")
	producer := p.getText(infoMap, "出品方")
	serials := p.getText(infoMap, "丛书")
	origin := p.getText(infoMap, "原作名")
	publisher := p.getText(infoMap, "出版社")
	pubdate := p.getText(infoMap, "出版年")
	pages := p.getText(infoMap, "页数")
	price := p.getText(infoMap, "定价")
	binding := p.getText(infoMap, "装帧")
	subtitle := p.getText(infoMap, "副标题")
	isbn13 := p.getText(infoMap, "ISBN")

	return DoubanBook{
		ID:          id,
		Author:      author,
		AuthorIntro: authorIntro,
		Translators: translators,
		Images: Image{
			Small:  smallImg,
			Medium: largeImg,
			Large:  largeImg,
		},
		Binding:   binding,
		Category:  "",
		Rating:    rating,
		ISBN13:    isbn13,
		Pages:     pages,
		Price:     price,
		Pubdate:   pubdate,
		Publisher: publisher,
		Producer:  producer,
		Serials:   serials,
		Subtitle:  subtitle,
		Summary:   summary,
		Title:     title,
		Tags:      tags,
		Origin:    origin,
	}
}

func (p *parser) parseSubjectCast(text string) ([]string, string, string) {
	subjects := strings.Split(text, "/")
	lenSub := len(subjects)
	pubdate := ""
	publisher := ""
	author := make([]string, 0)

	if lenSub >= 3 {
		pubdate = strings.TrimSpace(subjects[lenSub-1])
		publisher = strings.TrimSpace(subjects[lenSub-2])
		for i := 0; i < lenSub-2; i++ {
			author = append(author, strings.TrimSpace(subjects[i]))
		}
	} else if lenSub == 2 {
		author = append(author, strings.TrimSpace(subjects[0]))
		second := strings.TrimSpace(subjects[1])
		if _, err := parseInt(second); err == nil {
			pubdate = second
		} else {
			publisher = second
		}
	} else if lenSub == 1 && strings.TrimSpace(subjects[0]) != "" {
		author = append(author, strings.TrimSpace(subjects[0]))
	}

	return author, publisher, pubdate
}

func (p *parser) parseInfoText(s string) map[string]string {
	m := make(map[string]string)
	fixed := p.reRemoveSplitSpace.ReplaceAllString(s, "/")
	matches := p.reInfoPair.FindAllStringSubmatch(fixed, -1)
	for _, cap := range matches {
		if len(cap) >= 3 {
			m[strings.TrimSpace(cap[1])] = strings.TrimSpace(cap[2])
		}
	}
	return m
}

func (p *parser) getText(info map[string]string, key string) string {
	if v, ok := info[key]; ok {
		return v
	}
	return ""
}

func (p *parser) getTexts(info map[string]string, key string) []string {
	raw, ok := info[key]
	if !ok {
		return []string{}
	}
	items := strings.Split(raw, "/")
	res := make([]string, 0, len(items))
	for _, it := range items {
		v := strings.TrimSpace(it)
		if v != "" {
			res = append(res, v)
		}
	}
	return res
}

func captureGroup(re *regexp.Regexp, text string) string {
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}
