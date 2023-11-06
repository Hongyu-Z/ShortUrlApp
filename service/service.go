package service

import (
	"net/http"
)

type Service interface {
	UrlGetHandler(writer http.ResponseWriter, request *http.Request)
	UrlPostHandler(writer http.ResponseWriter, request *http.Request)
	UrlDeleteHandler(writer http.ResponseWriter, request *http.Request)
	UrlStatsGetHandler(writer http.ResponseWriter, request *http.Request)
}

type UrlService struct {
}

func NewUrlService() *UrlService {

}

func (s *UrlService) UrlGetHandler(writer http.ResponseWriter, request *http.Request) {

}

func (s *UrlService) UrlPostHandler(writer http.ResponseWriter, request *http.Request) {

}

func (s *UrlService) UrlDeleteHandler(writer http.ResponseWriter, request *http.Request) {

}

func (s *UrlService) UrlStatsGetHandler(writer http.ResponseWriter, request *http.Request) {

}
