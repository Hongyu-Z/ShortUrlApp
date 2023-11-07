package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ShortUrlApp/dao"
	"ShortUrlApp/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUrlGetHandler(t *testing.T) {
	mockUrlStatsDao := &dao.MockUrlStatsDao{}
	mockUrlRecordDao := &dao.MockUrlRecordDao{}
	mockUrlStatsDao.On("Save", mock.Anything).Twice().Return(nil)
	service := &UrlService{UrlStatsDao: mockUrlStatsDao, UrlRecordDao: mockUrlRecordDao}
	router := mux.NewRouter()
	router.HandleFunc("/url/{id}", service.UrlGetHandler).Methods(http.MethodGet)

	t.Run("valid id", func(t *testing.T) {
		// You should replace "valid_id" with an actual valid id
		mockUrlRecordDao.On("Find", "valid_url").Return(&models.UrlRecord{LongUrl: "long_url"}, nil)
		req, _ := http.NewRequest("GET", "/url/valid_url", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMovedPermanently, rr.Code)
		assert.Equal(t, "<a href=\"/url/long_url\">Moved Permanently</a>.\n\n", rr.Body.String())
	})

	t.Run("invalid url", func(t *testing.T) {
		mockUrlRecordDao.On("Find", "invalid_url").Return(nil, errors.New("record doesn't exist"))
		req, _ := http.NewRequest("GET", "/url/invalid_url", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request\n", rr.Body.String())
	})
}

func TestUrlPostHandler(t *testing.T) {
	mockUrlStatsDao := &dao.MockUrlStatsDao{}
	mockUrlRecordDao := &dao.MockUrlRecordDao{}
	service := &UrlService{UrlStatsDao: mockUrlStatsDao, UrlRecordDao: mockUrlRecordDao}
	router := mux.NewRouter()
	router.HandleFunc("/url", service.UrlPostHandler).Methods(http.MethodPost)

	t.Run("valid request", func(t *testing.T) {
		// You should replace "valid_id" with an actual valid id
		mockUrlRecordDao.On("Save", mock.Anything).Once().Return(nil)
		record := &models.UrlRecord{
			ShortUrl: "short_url",
			LongUrl:  "long_url",
		}
		payload, _ := json.Marshal(record)
		req, _ := http.NewRequest("POST", "/url", strings.NewReader(string(payload)))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "successfully added record: short_url -> long_url, with expiration at 0001-01-01 00:00:00 +0000 UTC", rr.Body.String())
	})

	t.Run("invalid request", func(t *testing.T) {
		mockUrlRecordDao.On("Save", mock.Anything).Once().Return(errors.New("record already exists"))
		record := &models.UrlRecord{
			ShortUrl: "short_url",
			LongUrl:  "long_url",
		}
		payload, _ := json.Marshal(record)
		req, _ := http.NewRequest("POST", "/url", strings.NewReader(string(payload)))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request\n", rr.Body.String())
	})
}

func TestUrlDeleteHandler(t *testing.T) {
	mockUrlStatsDao := &dao.MockUrlStatsDao{}
	mockUrlRecordDao := &dao.MockUrlRecordDao{}
	service := &UrlService{UrlStatsDao: mockUrlStatsDao, UrlRecordDao: mockUrlRecordDao}
	router := mux.NewRouter()
	router.HandleFunc("/url/{id}", service.UrlDeleteHandler).Methods(http.MethodDelete)

	t.Run("valid request", func(t *testing.T) {
		mockUrlRecordDao.On("Delete", mock.Anything).Once().Return(nil)
		req, _ := http.NewRequest(http.MethodDelete, "/url/valid_url", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "successfully deleted record: valid_url", rr.Body.String())
	})

	t.Run("invalid request", func(t *testing.T) {
		mockUrlRecordDao.On("Delete", mock.Anything).Once().Return(errors.New("record doesn't exist"))
		req, _ := http.NewRequest(http.MethodDelete, "/url/valid_url", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request\n", rr.Body.String())
	})
}

func TestUrlStatsGetHandler(t *testing.T) {
	mockUrlStatsDao := &dao.MockUrlStatsDao{}
	mockUrlRecordDao := &dao.MockUrlRecordDao{}
	service := &UrlService{UrlStatsDao: mockUrlStatsDao, UrlRecordDao: mockUrlRecordDao}
	router := mux.NewRouter()
	router.HandleFunc("/url/{id}/stats", service.UrlStatsGetHandler).Methods(http.MethodGet)

	t.Run("valid request", func(t *testing.T) {
		mockUrlStatsDao.On("GetCount", mock.Anything).Once().Return(1, 2, 3, nil)
		req, _ := http.NewRequest(http.MethodGet, "/url/valid_url/stats", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "oneDayCount:1 oneWeekCount:2 allTimeCount:3", rr.Body.String())
	})

	t.Run("invalid request", func(t *testing.T) {
		mockUrlStatsDao.On("GetCount", mock.Anything).Once().Return(0, 0, 0, errors.New("random error"))
		req, _ := http.NewRequest(http.MethodGet, "/url/invalid_url/stats", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, "Invalid request\n", rr.Body.String())
	})
}
