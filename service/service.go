package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
}

func NewUrlService() *UrlService {
	db := database.Init()
	urlRecordDao := &dao.UrlRecordDaoDBImpl{Db: db}
	return &UrlService{
		urlRecordDao: urlRecordDao,
	}
}

func (s *UrlService) UrlGetHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	shortUrl, ok := vars["id"]
	if !ok {
		handleInvalidRequest(writer, request)
		return
	}
	log.Printf("getting shortUrl:%s", shortUrl)
	if !ok {
		handleInvalidRequest(writer, request)
		return
	}

	longUrl, err := s.getUrlRecord(shortUrl)
	if err != nil {
		handleInvalidRequest(writer, request)
		return
	}

	http.Redirect(writer, request, longUrl, http.StatusMovedPermanently)
}

func (s *UrlService) UrlPostHandler(writer http.ResponseWriter, request *http.Request) {
	var record *models.UrlRecord
	decoder := json.NewDecoder(request.Body)
	if err := decoder.Decode(&record); err != nil {
		log.Printf("invalid")
		handleInvalidRequest(writer, request)
	}
	defer request.Body.Close()

	log.Printf("%s -> %s", record.ShortUrl, record.LongUrl)

	err := s.addUrlRecord(record)
	if err != nil {
		handleInvalidRequest(writer, request)
		return
	}

	writer.Write([]byte("added record: " + record.ShortUrl + " -> " + record.LongUrl))
}

func (s *UrlService) UrlDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	shortUrl, ok := vars["id"]
	if !ok {
		handleInvalidRequest(writer, request)
		return
	}
	log.Printf("deleting shortUrl:%s", shortUrl)
	if !ok {
		handleInvalidRequest(writer, request)
		return
	}

	err := s.deleteUrlRecord(shortUrl)
	if err != nil {
		handleInvalidRequest(writer, request)
		return
	}

	writer.Write([]byte("deleted record: " + shortUrl))
}

func (s *UrlService) UrlStatsGetHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	shortUrl, ok := vars["id"]
	if !ok {
		handleInvalidRequest(writer, request)
		return
	}

	oneDayCount, oneWeekCount, allTimeCount, err := s.getUrlStats(shortUrl)
	if err != nil {
		handleInvalidRequest(writer, request)
		return
	}

	writer.Write([]byte(fmt.Sprintf("oneDayCount:%d oneWeekCount:%d allTimeCount:%d", oneDayCount, oneWeekCount, allTimeCount)))
}
