package book

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/haigeek/douban-api-go/internal/httpclient"
)

const cacheSize = 100

type Service struct {
	client *httpclient.Client
	parser *parser
	cache  *expirable.LRU[string, DoubanBook]
}

func NewService(client *httpclient.Client) *Service {
	return &Service{
		client: client,
		parser: newParser(),
		cache:  expirable.NewLRU[string, DoubanBook](cacheSize, nil, 10*time.Minute),
	}
}

func (s *Service) Search(ctx context.Context, q string, count int) (DoubanBookResult, error) {
	list, err := s.getList(ctx, q, count)
	if err != nil {
		return DoubanBookResult{}, err
	}
	return DoubanBookResult{
		Code:  0,
		Msg:   "",
		Books: list,
	}, nil
}

func (s *Service) getList(ctx context.Context, q string, count int) ([]DoubanBook, error) {
	if q == "" {
		return []DoubanBook{}, nil
	}
	resp, err := s.client.Get(ctx, "https://www.douban.com/search", map[string]string{
		"cat": "1001",
		"q":   q,
	}, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return s.parser.parseSearchList(doc, count), nil
}

func (s *Service) GetBookInfoByISBN(ctx context.Context, isbn string) (DoubanBook, error) {
	if v, ok := s.cache.Get(isbn); ok {
		return v, nil
	}
	url := fmt.Sprintf("https://douban.com/isbn/%s/", isbn)
	return s.getBookInternal(ctx, url)
}

func (s *Service) GetBookInfo(ctx context.Context, id string) (DoubanBook, error) {
	if v, ok := s.cache.Get(id); ok {
		return v, nil
	}
	url := fmt.Sprintf("https://book.douban.com/subject/%s/", id)
	return s.getBookInternal(ctx, url)
}

func (s *Service) getBookInternal(ctx context.Context, rawURL string) (DoubanBook, error) {
	resp, err := s.client.Get(ctx, rawURL, nil, true)
	if err != nil {
		return DoubanBook{}, err
	}
	defer resp.Body.Close()

	finalID := ""
	parts := strings.Split(strings.Trim(resp.Request.URL.Path, "/"), "/")
	if len(parts) >= 2 {
		finalID = parts[len(parts)-1]
		if finalID == "" && len(parts) >= 2 {
			finalID = parts[len(parts)-2]
		}
	}
	if finalID == "" {
		finalID = strconv.FormatInt(time.Now().Unix(), 10)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return DoubanBook{}, err
	}

	info := s.parser.parseBookPage(doc, finalID)
	s.cache.Add(finalID, info)
	if info.ISBN13 != "" {
		s.cache.Add(info.ISBN13, info)
	}

	return info, nil
}

func parseFloat32(v string) (float32, error) {
	f, err := strconv.ParseFloat(v, 32)
	return float32(f), err
}

func parseInt(v string) (int, error) {
	return strconv.Atoi(v)
}
