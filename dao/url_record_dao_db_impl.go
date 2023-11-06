package dao

import (
	"database/sql"
	"log"

	"ShortUrlApp/models"
)

type UrlRecordDaoDBImpl struct {
	Db *sql.DB
}

func (impl *UrlRecordDaoDBImpl) Save(record *models.UrlRecord) error {
	query := "INSERT INTO urls(short_url, long_url, expire_at) VALUES (?, ?, ?)"
	stmt, err := impl.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(record.ShortUrl, record.LongUrl, record.ExpireAt)
	if err != nil {
		return err
	}

	return nil
}

func (impl *UrlRecordDaoDBImpl) Find(shortUrl string) (*models.UrlRecord, error) {
	query := "SELECT short_url, long_url, expire_at FROM urls WHERE short_url = ?"
	stmt, err := impl.Db.Prepare(query)
	defer stmt.Close()
	result := &models.UrlRecord{}
	err = stmt.QueryRow(shortUrl).Scan(&result.ShortUrl, &result.LongUrl, &result.ExpireAt)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return result, nil
}

func (impl *UrlRecordDaoDBImpl) Delete(shortUrl string) error {
	query := "DELETE from urls WHERE short_url = ?"
	stmt, err := impl.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(shortUrl)
	if err != nil {
		return err
	}

	return nil
}
