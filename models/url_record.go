package models

import "time"

type UrlRecord struct {
	ShortUrl string    `json:"short_url"`
	LongUrl  string    `json:"long_url"`
	ExpireAt time.Time `json:"expire_at"`
}
