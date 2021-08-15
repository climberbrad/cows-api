package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// (GET) localhost:8080/v0/cows
// (POST) localhost:8080/v0/cows/id
// (PUT) localhost:8080/v0/cows -d {}

type Cow struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Date   string `json:"date"`
	Image  string `json:"image"`
	Finder string `json:"finder"`
}

type cowHandlers struct {
	sync.Mutex
	store map[string]Cow
}

func (h *cowHandlers) cows(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Unsupported Method."))
		return
	}
}

func (h *cowHandlers) get(w http.ResponseWriter, r *http.Request) {
	cows := make([]Cow, len(h.store))

	h.Lock()
	i := 0
	for _, cow := range h.store {
		cows[i] = cow
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(cows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *cowHandlers) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	var cow Cow
	err = json.Unmarshal(bodyBytes, &cow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Needed content-type: application/json but got '%s'",
			ct)))
	}

	h.Lock()
	h.store[cow.ID] = cow
	defer h.Unlock()
}

func newCowHandlers() *cowHandlers {
	return &cowHandlers{
		store: map[string]Cow{},
	}
}

func main() {
	cowHandlers := newCowHandlers()
	http.HandleFunc("/v0/cows", cowHandlers.cows)
	http.ListenAndServe(":8080", nil)

}
