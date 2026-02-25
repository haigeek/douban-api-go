package main

import (
	"fmt"
	"log"

	"github.com/haigeek/douban-api-go/internal/api/movie"
	"github.com/haigeek/douban-api-go/internal/book"
	"github.com/haigeek/douban-api-go/internal/config"
	"github.com/haigeek/douban-api-go/internal/httpclient"
	"github.com/haigeek/douban-api-go/internal/media"
	"github.com/haigeek/douban-api-go/internal/server"
)

func main() {
	cfg := config.Load()

	client, err := httpclient.New(cfg.Cookie)
	if err != nil {
		log.Fatalf("create http client failed: %v", err)
	}

	movieService := movie.NewService(client)
	bookService := book.NewService(client)
	mediaService := media.NewService(client)
	h := server.NewHandlers(movieService, cfg)
	b := book.NewHandlers(bookService)
	m := media.NewHandlers(mediaService)
	r := server.NewRouter(h, b, m, cfg.Debug)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server start failed: %v", err)
	}
}
