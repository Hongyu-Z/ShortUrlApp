package dao

import "ShortUrlApp/models"

//go:generate mockery --name UrlStatsDao --inpackage --case underscore
type UrlStatsDao interface {
	Save(event *models.UrlStatsEvent) error
	GetCount(shortUrl string) (int, int, int, error)
}
