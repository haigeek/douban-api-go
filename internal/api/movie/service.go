package movie

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/haigeek/douban-api-go/internal/httpclient"
)

const cacheSize = 100

type Service struct {
	client     *httpclient.Client
	parser     *parser
	movieCache *expirable.LRU[string, MovieInfo]
	photoCache *expirable.LRU[string, []Photo]
}

func NewService(client *httpclient.Client) *Service {
	return &Service{
		client:     client,
		parser:     newParser(),
		movieCache: expirable.NewLRU[string, MovieInfo](cacheSize, nil, 1*time.Minute),
		photoCache: expirable.NewLRU[string, []Photo](cacheSize, nil, 1*time.Minute),
	}
}

func (s *Service) Search(ctx context.Context, q string, limit int, imageSize string) ([]Movie, error) {
	if q == "" {
		return []Movie{}, nil
	}
	doc, err := s.fetchDocument(ctx, "https://www.douban.com/search", map[string]string{
		"cat": "1002",
		"q":   q,
	})
	if err != nil {
		return nil, err
	}
	return s.parser.parseMovies(doc, limit, imageSize), nil
}

func (s *Service) SearchFull(ctx context.Context, q string, limit int, imageSize string) ([]MovieInfo, error) {
	movies, err := s.Search(ctx, q, limit, imageSize)
	if err != nil {
		return nil, err
	}
	list := make([]MovieInfo, 0, len(movies))
	for _, m := range movies {
		info, infoErr := s.GetMovieInfo(ctx, m.SID, imageSize)
		if infoErr != nil {
			return nil, infoErr
		}
		list = append(list, info)
	}
	return list, nil
}

func (s *Service) GetMovieInfo(ctx context.Context, sid, imageSize string) (MovieInfo, error) {
	cacheKey := fmt.Sprintf("movie_%s_%s", sid, imageSize)
	if v, ok := s.movieCache.Get(cacheKey); ok {
		return v, nil
	}

	doc, err := s.fetchDocument(ctx, fmt.Sprintf("https://movie.douban.com/subject/%s/", sid), nil)
	if err != nil {
		return MovieInfo{}, err
	}

	info := s.parser.parseMovieInfo(doc, sid, imageSize)
	s.movieCache.Add(cacheKey, info)
	return info, nil
}

func (s *Service) GetCelebrities(ctx context.Context, sid string) ([]Celebrity, error) {
	doc, err := s.fetchDocument(ctx, fmt.Sprintf("https://movie.douban.com/subject/%s/celebrities", sid), nil)
	if err != nil {
		return nil, err
	}
	return s.parser.parseCelebrities(doc), nil
}

func (s *Service) GetCelebrity(ctx context.Context, id string) (CelebrityInfo, error) {
	doc, err := s.fetchDocument(ctx, fmt.Sprintf("https://movie.douban.com/celebrity/%s/", id), nil)
	if err != nil {
		return CelebrityInfo{}, err
	}
	return s.parser.parseCelebrityInfo(doc, id), nil
}

func (s *Service) GetWallpaper(ctx context.Context, sid string) ([]Photo, error) {
	if v, ok := s.photoCache.Get(sid); ok {
		return v, nil
	}

	doc, err := s.fetchDocument(ctx, fmt.Sprintf("https://movie.douban.com/subject/%s/photos", sid), map[string]string{
		"type":    "W",
		"start":   strconv.Itoa(0),
		"sortby":  "size",
		"size":    "a",
		"subtype": "a",
	})
	if err != nil {
		return nil, err
	}

	photos := s.parser.parseWallpapers(doc)
	s.photoCache.Add(sid, photos)
	return photos, nil
}

func (s *Service) ProxyImage(ctx context.Context, rawURL string) (*http.Response, []byte, error) {
	resp, err := s.client.Get(ctx, rawURL, nil, false)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	clone := *resp
	clone.Body = http.NoBody
	return &clone, body, nil
}

func (s *Service) fetchDocument(ctx context.Context, rawURL string, query map[string]string) (*goquery.Document, error) {
	resp, err := s.client.Get(ctx, rawURL, query, true)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return goquery.NewDocumentFromReader(resp.Body)
}
