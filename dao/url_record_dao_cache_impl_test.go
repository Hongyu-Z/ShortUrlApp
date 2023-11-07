package dao

import (
	"testing"
	"time"

	"ShortUrlApp/models"
	"github.com/jellydator/ttlcache/v3"
	"github.com/stretchr/testify/assert"
)

func TestCacheSave(t *testing.T) {
	cache := ttlcache.New[string, string]()
	dao := &UrlRecordDaoCacheImpl{Cache: cache}

	record := &models.UrlRecord{
		ShortUrl: "short_url",
		LongUrl:  "long_url",
		ExpireAt: time.Now().Add(time.Hour * 24),
	}

	//adding a new record
	err := dao.Save(record)
	assert.Nil(t, err)

	//finding the added record
	result, err := dao.Find(record.ShortUrl)
	assert.Nil(t, err)
	assert.Equal(t, record.ShortUrl, result.ShortUrl)
	assert.Equal(t, record.LongUrl, result.LongUrl)

	//deleting the added record
	err = dao.Delete(record.ShortUrl)
	assert.Nil(t, err)

	//finding the added record again after deletion
	result, err = dao.Find(record.ShortUrl)
	assert.Nil(t, result)
	assert.Equal(t, "record doesn't exist", err.Error())
}
