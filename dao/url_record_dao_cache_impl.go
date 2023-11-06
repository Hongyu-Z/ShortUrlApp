package dao

import (
	"errors"
	"time"

	"ShortUrlApp/models"
	"github.com/jellydator/ttlcache/v3"
)

type UrlRecordDaoCacheImpl struct {
	Cache *ttlcache.Cache[string, string]
}

func (impl *UrlRecordDaoCacheImpl) Save(record *models.UrlRecord) error {
	duration := ttlcache.NoTTL
	if !record.ExpireAt.IsZero() {
		duration = record.ExpireAt.Sub(time.Now())
	}

	impl.Cache.Set(record.ShortUrl, record.LongUrl, duration)
	return nil
}

func (impl *UrlRecordDaoCacheImpl) Find(shortUrl string) (*models.UrlRecord, error) {
	item := impl.Cache.Get(shortUrl)
	if item == nil {
		return nil, errors.New("record doesn't exist")
	}

	return &models.UrlRecord{ShortUrl: shortUrl, LongUrl: item.Value(), ExpireAt: item.ExpiresAt()}, nil
}

func (impl *UrlRecordDaoCacheImpl) Delete(shortUrl string) error {
	impl.Cache.Delete(shortUrl)

	return nil
}
