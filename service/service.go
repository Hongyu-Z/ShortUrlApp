package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ShortUrlApp/cache"
	"ShortUrlApp/dao"
	"ShortUrlApp/database"
	"ShortUrlApp/models"
	"github.com/gorilla/mux"
)

type Service interface {
	UrlGetHandler(writer http.ResponseWriter, request *http.Request)
	UrlPostHandler(writer http.ResponseWriter, request *http.Request)
	UrlDeleteHandler(writer http.ResponseWriter, request *http.Request)
	UrlStatsGetHandler(writer http.ResponseWriter, request *http.Request)
}

type UrlService struct {
	urlRecordDao dao.UrlRecordDao
	UrlStatsDao  dao.UrlStatsDao
}

func NewUrlService() *UrlService {
	db := database.Init()
	cache := cache.Init()
	urlRecordDao := &dao.UrlRecordDaoImpl{
		DbStrategy:    &dao.UrlRecordDaoDBImpl{Db: db},
		CacheStrategy: &dao.UrlRecordDaoCacheImpl{Cache: cache},
	}
	urlStatsDao := &dao.UrlStatsDaoDBImpl{Db: db}

	return &UrlService{
		urlRecordDao: urlRecordDao,
		UrlStatsDao:  urlStatsDao,
	}
}

func (s *UrlService) UrlGetHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	shortUrl, ok := vars["id"]
	if !ok {
		log.Print("can't find id in request")
		handleInvalidRequest(writer, request)
		return
	}

	longUrl, err := s.getUrlRecord(shortUrl)
	if err != nil {
		log.Print(err)
		handleInvalidRequest(writer, request)
		return
	}

	log.Printf(fmt.Sprintf("successfully found record for %s, redirecting user to %s", shortUrl, longUrl))
	http.Redirect(writer, request, longUrl, http.StatusMovedPermanently)
}

func (s *UrlService) UrlPostHandler(writer http.ResponseWriter, request *http.Request) {
	var record *models.UrlRecord
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&record); err != nil {
		log.Print(err)
		handleInvalidRequest(writer, request)
		return
	}
	defer request.Body.Close()

	err := s.addUrlRecord(record)
	if err != nil {
		log.Print(err)
		handleInvalidRequest(writer, request)
		return
	}

	log.Printf("successfully added record: %v", record)
	writer.Write([]byte(fmt.Sprintf("successfully added record: %s -> %s, with expiration at %s", record.ShortUrl, record.LongUrl, record.ExpireAt)))
}

func (s *UrlService) UrlDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	shortUrl, ok := vars["id"]
	if !ok {
		log.Print("can't find id in request")
		handleInvalidRequest(writer, request)
		return
	}
	log.Printf("deleting shortUrl:%s", shortUrl)
	if !ok {
		log.Print("can't delete record")
		handleInvalidRequest(writer, request)
		return
	}

	err := s.deleteUrlRecord(shortUrl)
	if err != nil {
		log.Print(err)
		handleInvalidRequest(writer, request)
		return
	}

	log.Printf("successfully deleted record: " + shortUrl)
	writer.Write([]byte("successfully deleted record: " + shortUrl))
}

func (s *UrlService) UrlStatsGetHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	shortUrl, ok := vars["id"]
	if !ok {
		log.Print("can't find id in request")
		handleInvalidRequest(writer, request)
		return
	}

	oneDayCount, oneWeekCount, allTimeCount, err := s.getUrlStats(shortUrl)
	if err != nil {
		log.Print(err)
		handleInvalidRequest(writer, request)
		return
	}

	log.Printf("successfully processed get url stats request for %s. oneDayCount:%d oneWeekCount:%d allTimeCount:%d", shortUrl, oneDayCount, oneWeekCount, allTimeCount)
	writer.Write([]byte(fmt.Sprintf("oneDayCount:%d oneWeekCount:%d allTimeCount:%d", oneDayCount, oneWeekCount, allTimeCount)))
}
