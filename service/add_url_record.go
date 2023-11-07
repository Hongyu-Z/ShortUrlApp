package service

import (
	"ShortUrlApp/models"
)

func (s *UrlService) addUrlRecord(record *models.UrlRecord) error {
	return s.UrlRecordDao.Save(record)
}
