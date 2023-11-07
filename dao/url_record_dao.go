package dao

import (
	"ShortUrlApp/models"
)

//go:generate mockery --name UrlRecordDao --inpackage --case underscore
type UrlRecordDao interface {
	Save(record *models.UrlRecord) error
	Find(shortUrl string) (*models.UrlRecord, error)
	Delete(shortUrl string) error
}
