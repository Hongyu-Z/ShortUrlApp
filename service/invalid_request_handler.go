package service

import "net/http"

func handleInvalidRequest(writer http.ResponseWriter, request *http.Request) {
	http.Error(writer, "Invalid request", http.StatusBadRequest)
}
