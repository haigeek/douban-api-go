package media

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/haigeek/douban-api-go/internal/httpclient"
)

const recentHotAPI = "https://m.douban.com/rexxar/api/v2/subject/recent_hot/%s"

type Service struct {
	client *httpclient.Client
}

func NewService(client *httpclient.Client) *Service {
	return &Service{client: client}
}

func (s *Service) RecentHot(ctx context.Context, subject, category, mediaType string, start, limit int) (HotMediaResponse, error) {
	resp, err := s.client.Get(ctx, fmt.Sprintf(recentHotAPI, subject), map[string]string{
		"start":    fmt.Sprintf("%d", start),
		"limit":    fmt.Sprintf("%d", limit),
		"category": category,
		"type":     mediaType,
	}, true)
	if err != nil {
		return HotMediaResponse{}, err
	}
	defer resp.Body.Close()

	var out HotMediaResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return HotMediaResponse{}, err
	}
	return out, nil
}

func (s *Service) HotTV(ctx context.Context, start, limit int) (HotMediaResponse, error) {
	return s.RecentHot(ctx, "tv", "tv", "tv", start, limit)
}

func (s *Service) HotMovie(ctx context.Context, start, limit int) (HotMediaResponse, error) {
	return s.RecentHot(ctx, "movie", "热门", "全部", start, limit)
}

func (s *Service) LatestMovie(ctx context.Context, start, limit int) (HotMediaResponse, error) {
	return s.RecentHot(ctx, "movie", "最新", "全部", start, limit)
}

func (s *Service) HighRatingMovie(ctx context.Context, start, limit int) (HotMediaResponse, error) {
	return s.RecentHot(ctx, "movie", "豆瓣高分", "全部", start, limit)
}
