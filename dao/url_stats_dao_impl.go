package dao

import (
	"database/sql"
	"time"

	"ShortUrlApp/models"
)

type UrlStatsDaoDBImpl struct {
	Db *sql.DB
}

func (impl UrlStatsDaoDBImpl) Save(event *models.UrlStatsEvent) error {
	stmt, err := impl.Db.Prepare("INSERT INTO stats(short_url, created_at) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(event.ShortUrl, event.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (impl UrlStatsDaoDBImpl) GetCount(shortUrl string) (int, int, int, error) {
	stmt, err := impl.Db.Prepare("SELECT count(*) FROM stats WHERE short_url = ? AND created_at >= ?")
	if err != nil {
		return 0, 0, 0, err
	}
	defer stmt.Close()
	var oneDayCount, oneWeekCount, allTimeCount int

	err1 := stmt.QueryRow(shortUrl, time.Now().AddDate(0, 0, -1)).Scan(&oneDayCount)
	err2 := stmt.QueryRow(shortUrl, time.Now().AddDate(0, 0, -7)).Scan(&oneWeekCount)
	err3 := stmt.QueryRow(shortUrl, 0).Scan(&allTimeCount)
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, err
	}

	return oneDayCount, oneWeekCount, allTimeCount, nil
}
