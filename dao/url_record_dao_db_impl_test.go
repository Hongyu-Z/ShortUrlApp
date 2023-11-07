package dao

import (
	"testing"
	"time"

	"ShortUrlApp/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUrlRecordDaoDBImpl_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dao := &UrlRecordDaoDBImpl{Db: db}
	record := &models.UrlRecord{
		ShortUrl: "short_url",
		LongUrl:  "long_url",
		ExpireAt: time.Now(),
	}
	mock.ExpectPrepare("INSERT INTO urls\\(short_url, long_url, expire_at\\) VALUES \\(\\?, \\?, \\?\\)")
	mock.ExpectExec("INSERT INTO urls\\(short_url, long_url, expire_at\\) VALUES \\(\\?, \\?, \\?\\)").
		WithArgs(record.ShortUrl, record.LongUrl, record.ExpireAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.Save(record)

	assert.Nil(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUrlRecordDaoDBImpl_Find(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dao := &UrlRecordDaoDBImpl{Db: db}

	shortUrl := "short_url"
	longUrl := "long_url"
	expireAt := time.Now()

	rows := sqlmock.NewRows([]string{"short_url", "long_url", "expire_at"}).
		AddRow(shortUrl, longUrl, expireAt)
	mock.ExpectPrepare("SELECT short_url, long_url, expire_at FROM urls WHERE short_url = ?")
	mock.ExpectQuery("SELECT short_url, long_url, expire_at FROM urls WHERE short_url = ?").
		WithArgs(shortUrl).
		WillReturnRows(rows)

	result, err := dao.Find(shortUrl)

	assert.Nil(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	assert.Equal(t, shortUrl, result.ShortUrl)
	assert.Equal(t, longUrl, result.LongUrl)
	assert.Equal(t, expireAt, result.ExpireAt)
}

func TestUrlRecordDaoDBImpl_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	dao := &UrlRecordDaoDBImpl{Db: db}
	shortUrl := "short_url"
	mock.ExpectPrepare("DELETE from urls WHERE short_url = ?")
	mock.ExpectExec("DELETE from urls WHERE short_url = ?").
		WithArgs(shortUrl).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.Delete(shortUrl)

	assert.Nil(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
