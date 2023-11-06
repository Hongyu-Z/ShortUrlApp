package dao

import (
	"ShortUrlApp/models"
)

type UrlRecordDao interface {
	Save(record *models.UrlRecord) error
	Find(shortUrl string) (*models.UrlRecord, error)
	Delete(shortUrl string) error
}
