package ShortUrlApp

import (
	"log"
	"net/http"
	"time"

	"ShortUrlApp/service"
	"github.com/gorilla/mux"
)

func main() {
	log.Printf("service starting")

	service := service.NewUrlService()
	r := mux.NewRouter()
	r.HandleFunc("/url/{id}", service.UrlGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/url/{id}", service.UrlDeleteHandler).Methods(http.MethodDelete)
	r.HandleFunc("/url/{id}/stats", service.UrlStatsGetHandler).Methods(http.MethodGet)
	r.HandleFunc("/url", service.UrlPostHandler).Methods(http.MethodPost)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
