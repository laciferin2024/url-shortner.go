package app

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/laciferin2024/url-shortner.go/models"
)

type Services interface {
	ShortenUrl(ctx context.Context, url string) (shortenedUrl string, err error)
	RetrieveOriginalUrl(shortUrl string) (url string, err error)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// randStringBytes - Create random short link
func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *service) ShortenUrl(ctx context.Context, url string) (shortenedUrl string, err error) {

	if url == "" {
		s.Log.Errorln("url is empty")
		return "", fmt.Errorf("url is empty")
	}

	// Check if URL already exists
	var existing models.Url
	err = s.db.NewSelect().Model(&existing).Where("urls = ?", url).Scan(ctx)
	if err == nil {
		return existing.ShortenedUrl, nil
	}

	for i := 0; i < 5; i++ {
		shortenedUrl = randStringBytes(10)

		newUrl := &models.Url{
			Url:          url,
			ShortenedUrl: shortenedUrl,
		}

		_, err = s.db.NewInsert().Model(newUrl).Exec(ctx)
		if err == nil {
			return shortenedUrl, nil
		}

		s.Log.Errorf("failed to insert url (attempt %d): %v", i+1, err)

		// If context is canceled, stop
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
	}
	return "", fmt.Errorf("failed to shorten url after retries: %v", err)
}

func (s *service) RetrieveOriginalUrl(shortUrl string) (url string, err error) {

	// 1. Check Cache
	if !s.getCache(shortUrl, &url) {
		// Cache Hit
		if !strings.Contains(url, "http") {
			url = fmt.Sprintf("https://%s", url)
		}
		// Async update stats (using shortUrl since we don't have ID)
		go s.updateStats(shortUrl)
		return
	}

	// 2. Cache Miss - Query DB
	var urlModel models.Url
	err = s.db.NewSelect().Model(&urlModel).Where("short_urls = ?", shortUrl).Scan(context.Background())

	if err != nil {
		if err == sql.ErrNoRows {
			err = fmt.Errorf("url %s not found", shortUrl)
		}
		return
	}

	url = urlModel.Url

	if !strings.Contains(url, "http") {
		url = fmt.Sprintf("https://%s", url)
	}

	// 3. Set Cache (Async or Sync? Sync is safer for next request)
	// Using 24 hours expiration
	_ = s.setCache(shortUrl, url, 24*time.Hour)

	// 4. Update Stats
	go s.updateStats(shortUrl)

	return
}

func (s *service) updateStats(shortUrl string) {
	_, err := s.db.NewUpdate().
		Model((*models.Url)(nil)).
		Set("click_count = click_count + 1").
		Set("last_accessed_at = ?", time.Now()).
		Where("short_urls = ?", shortUrl).
		Exec(context.Background())
	if err != nil {
		s.Log.Errorln("failed to update stats:", err)
	}
}
