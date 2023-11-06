package service

import (
	"log"
	"time"

	"ShortUrlApp/models"
)

func (s *UrlService) getUrlRecord(shortUrl string) (string, error) {
	record, err := s.urlRecordDao.Find(shortUrl)
	if err != nil {
		return "", err
	}

	go s.reportStats(shortUrl)

	return record.LongUrl, nil
}

func (s *UrlService) reportStats(shortUrl string) {
	event := &models.UrlStatsEvent{ShortUrl: shortUrl, CreatedAt: time.Now()}
	err := s.UrlStatsDao.Create(event)
	if err != nil {
		log.Print(err)
	}
}
