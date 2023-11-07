package dao

import (
	"errors"
	"testing"
	"time"

	"ShortUrlApp/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jellydator/ttlcache/v3"
	"github.com/stretchr/testify/assert"
)

func TestUrlRecordDaoImpl_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	cache := ttlcache.New[string, string]()
	dao := &UrlRecordDaoImpl{DbStrategy: &UrlRecordDaoDBImpl{Db: db}, CacheStrategy: &UrlRecordDaoCacheImpl{Cache: cache}}
	record := &models.UrlRecord{
		ShortUrl: "short_url",
		LongUrl:  "long_url",
		ExpireAt: time.Now(),
	}
	t.Run("happy path", func(t *testing.T) {
		cache.DeleteAll()
		mock.ExpectPrepare("INSERT INTO urls\\(short_url, long_url, expire_at\\) VALUES \\(\\?, \\?, \\?\\)")
		mock.ExpectExec("INSERT INTO urls\\(short_url, long_url, expire_at\\) VALUES \\(\\?, \\?, \\?\\)").
			WithArgs(record.ShortUrl, record.LongUrl, record.ExpireAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		dao.Save(record)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, 1, cache.Len())
	})

	t.Run("when DB write fails, cache write should be skipped", func(t *testing.T) {
		cache.DeleteAll()
		mock.ExpectPrepare("INSERT INTO urls\\(short_url, long_url, expire_at\\) VALUES \\(\\?, \\?, \\?\\)")
		mock.ExpectExec("INSERT INTO urls\\(short_url, long_url, expire_at\\) VALUES \\(\\?, \\?, \\?\\)").
			WithArgs(record.ShortUrl, record.LongUrl, record.ExpireAt).
			WillReturnError(errors.New("Random error"))

		dao.Save(record)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, 0, cache.Len())
	})
}

func TestUrlRecordDaoImpl_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	cache := ttlcache.New[string, string]()
	dao := &UrlRecordDaoImpl{DbStrategy: &UrlRecordDaoDBImpl{Db: db}, CacheStrategy: &UrlRecordDaoCacheImpl{Cache: cache}}
	record := &models.UrlRecord{
		ShortUrl: "short_url",
		LongUrl:  "long_url",
		ExpireAt: time.Now().Add(time.Hour * 24),
	}
	expiredRecord := &models.UrlRecord{
		ShortUrl: "short_url",
		LongUrl:  "long_url",
		ExpireAt: time.Now().Add(time.Hour * -24),
	}

	t.Run("cache hit", func(t *testing.T) {
		cache.DeleteAll()
		cache.Set(record.ShortUrl, record.LongUrl, record.ExpireAt.Sub(time.Now()))

		result, err := dao.Find(record.ShortUrl)
		assert.Nil(t, err)
		assert.Equal(t, record.ShortUrl, result.ShortUrl)
		assert.Equal(t, record.LongUrl, result.LongUrl)
	})

	t.Run("cache miss", func(t *testing.T) {
		cache.DeleteAll()
		rows := sqlmock.NewRows([]string{"short_url", "long_url", "expire_at"}).
			AddRow(record.ShortUrl, record.LongUrl, record.ExpireAt)
		mock.ExpectPrepare("SELECT short_url, long_url, expire_at FROM urls WHERE short_url = ?")
		mock.ExpectQuery("SELECT short_url, long_url, expire_at FROM urls WHERE short_url = ?").
			WithArgs(record.ShortUrl).
			WillReturnRows(rows)

		result, err := dao.Find(record.ShortUrl)

		assert.Nil(t, err)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, record.ShortUrl, result.ShortUrl)
		assert.Equal(t, record.LongUrl, result.LongUrl)
		assert.Equal(t, record.ExpireAt, result.ExpireAt)
	})

	t.Run("cache miss - expired record", func(t *testing.T) {
		cache.DeleteAll()
		rows := sqlmock.NewRows([]string{"short_url", "long_url", "expire_at"}).
			AddRow(expiredRecord.ShortUrl, expiredRecord.LongUrl, expiredRecord.ExpireAt)
		mock.ExpectPrepare("SELECT short_url, long_url, expire_at FROM urls WHERE short_url = ?")
		mock.ExpectQuery("SELECT short_url, long_url, expire_at FROM urls WHERE short_url = ?").
			WithArgs(record.ShortUrl).
			WillReturnRows(rows)

		result, err := dao.Find(record.ShortUrl)

		assert.NotNil(t, err)
		assert.Equal(t, "record doesn't exist", err.Error())
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Nil(t, result)
	})
}

func TestUrlRecordDaoImpl_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	cache := ttlcache.New[string, string]()
	dao := &UrlRecordDaoImpl{DbStrategy: &UrlRecordDaoDBImpl{Db: db}, CacheStrategy: &UrlRecordDaoCacheImpl{Cache: cache}}
	record := &models.UrlRecord{
		ShortUrl: "short_url",
		LongUrl:  "long_url",
		ExpireAt: time.Now().Add(time.Hour * 24),
	}
	t.Run("when cache has the record", func(t *testing.T) {
		cache.DeleteAll()
		cache.Set(record.ShortUrl, record.LongUrl, record.ExpireAt.Sub(time.Now()))
		mock.ExpectPrepare("DELETE from urls WHERE short_url = ?")
		mock.ExpectExec("DELETE from urls WHERE short_url = ?").
			WithArgs(record.ShortUrl).
			WillReturnResult(sqlmock.NewResult(1, 1))
		assert.Equal(t, 1, cache.Len())

		err = dao.Delete(record.ShortUrl)

		assert.Nil(t, err)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, 0, cache.Len())
	})

	t.Run("when cache doesn't have the record, DB should still execute the DELETE command", func(t *testing.T) {
		cache.DeleteAll()
		mock.ExpectPrepare("DELETE from urls WHERE short_url = ?")
		mock.ExpectExec("DELETE from urls WHERE short_url = ?").
			WithArgs(record.ShortUrl).
			WillReturnResult(sqlmock.NewResult(1, 1))
		assert.Equal(t, 0, cache.Len())

		err = dao.Delete(record.ShortUrl)

		assert.Nil(t, err)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		assert.Equal(t, 0, cache.Len())
	})
}
