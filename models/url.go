package models

import (
	"time"

	"github.com/hiroBzinga/bun"
)

type Url struct {
	bun.BaseModel `bun:"table:urls"`
	ID            int    `bun:"id"`
	Url           string `bun:"urls"`
	ShortenedUrl  string `bun:"short_urls"`

	ClickCount     int64     `bun:"click_count,default:0" json:"clickCount"`
	LastAccessedAt time.Time `bun:"last_accessed_at,nullzero" json:"lastAccessedAt,omitempty"`

	CreatedAt time.Time `bun:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `bun:"updated_at" json:"updatedAt,omitempty"`
	DeletedAt time.Time `bun:"deleted_at,nullzero,soft_delete" json:"deletedAt,omitempty"`
}
