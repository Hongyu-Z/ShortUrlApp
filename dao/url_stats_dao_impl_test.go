package dao

import (
	"testing"
	"time"

	"ShortUrlApp/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUrlStatsDaoImpl_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	dao := &UrlStatsDaoDBImpl{Db: db}
	event := &models.UrlStatsEvent{
		ShortUrl:  "short_url",
		CreatedAt: time.Now(),
	}
	mock.ExpectPrepare("INSERT INTO stats\\(short_url, created_at\\) values\\(\\?, \\?\\)")
	mock.ExpectExec("INSERT INTO stats\\(short_url, created_at\\) values\\(\\?, \\?\\)").
		WithArgs(event.ShortUrl, event.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = dao.Save(event)

	assert.Nil(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUrlStatsDaoImpl_GetCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	dao := &UrlStatsDaoDBImpl{Db: db}
	shortUrl := "short_url"
	mock.ExpectPrepare("SELECT count\\(\\*\\) FROM stats WHERE short_url = \\? AND created_at >= \\?")
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM stats WHERE short_url = \\? AND created_at >= \\?").
		WithArgs(shortUrl, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
	mock.ExpectQuery("SELECT count\\(\\*\\) FROM stats WHERE short_url = \\? AND created_at >= \\?").
		WithArgs(shortUrl, 0).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))

	_, _, _, err = dao.GetCount(shortUrl)

	assert.Nil(t, err)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
