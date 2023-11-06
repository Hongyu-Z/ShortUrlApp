package dao

import "ShortUrlApp/models"

type UrlStatsDao interface {
	Create(event *models.UrlStatsEvent) error
	GetCount(shortUrl string) (int, int, int, error)
}
