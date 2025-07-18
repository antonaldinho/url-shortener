package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"url-shortener/models"
)

type Body struct {
	URL string `json:"url"`
}

type URLs struct {
	l *log.Logger
}

func NewURLs(l *log.Logger) *URLs {
	return &URLs{l}
}

func (u *URLs) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		u.getURLFromId(rw, r)
		return
	case http.MethodPost:
		u.createShortURL(rw, r)
	}
}

func (u *URLs) getURLFromId(rw http.ResponseWriter, r *http.Request) {
	// Reading id from URL query params
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(rw, "Id parameter is required.", http.StatusBadRequest)
		return
	}

	url, err := models.GetURL(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
	}
	rw.Header().Add("content-type", "application/json")
	json.NewEncoder(rw).Encode(url)
}

func (u *URLs) createShortURL(rw http.ResponseWriter, r *http.Request) {
	//Parsing the body
	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		u.l.Println("error parsing the body")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	url, err := models.CreateShortURL(body.URL)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Header().Add("content-type", "application/json")
	json.NewEncoder(rw).Encode(url)
}
