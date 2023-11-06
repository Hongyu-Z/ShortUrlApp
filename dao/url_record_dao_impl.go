package dao

import (
	"errors"
	"time"

	"ShortUrlApp/models"
)

type UrlRecordDaoImpl struct {
	DbStrategy    *UrlRecordDaoDBImpl
	CacheStrategy *UrlRecordDaoCacheImpl
}

func (impl *UrlRecordDaoImpl) Save(record *models.UrlRecord) error {
	err := impl.DbStrategy.Save(record)
	if err != nil {
		return err
	}

	impl.CacheStrategy.Save(record)
	return nil
}

func (impl *UrlRecordDaoImpl) Find(shortUrl string) (*models.UrlRecord, error) {
	record, err := impl.CacheStrategy.Find(shortUrl)
	if err == nil {
		return record, nil
	}

	record, err = impl.DbStrategy.Find(shortUrl)
	if err != nil {
		return nil, err
	}
	if !record.ExpireAt.IsZero() && record.ExpireAt.Before(time.Now()) {
		impl.Delete(record.ShortUrl)
		return nil, errors.New("record doesn't exist")
	}

	impl.CacheStrategy.Save(record)
	return record, err
}

func (impl *UrlRecordDaoImpl) Delete(shortUrl string) error {
	err := impl.CacheStrategy.Delete(shortUrl)
	if err != nil {
		return err
	}

	err = impl.DbStrategy.Delete(shortUrl)
	return err
}
