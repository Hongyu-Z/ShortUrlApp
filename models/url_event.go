package models

import "time"

type UrlStatsEvent struct {
	ShortUrl  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
}
