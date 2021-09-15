package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Cow struct {
	Name        string `json:"name,omitempty"`
	ID          string `json:"id,omitempty"`
	Date        string `json:"date,omitempty"`
	Image       string `json:"image,omitempty"`
	Finder      string `json:"finder,omitempty"`
	Description string `json:"description,omitempty"`
}

type cowsResource struct {
	sync.Mutex
	store map[string]Cow
}

func newHandler() *cowsResource {
	return &cowsResource{
		store: map[string]Cow{},
	}
}

func (rs cowsResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, //nolint:gomnd // Maximum value not ignored by any of major browsers
	}))

	r.Get("/", rs.List)    // GET /posts - Read a list of posts.
	r.Post("/", rs.Create) // POST /posts - Create a new post.

	r.Route("/{id}", func(r chi.Router) {
		r.Use(PostCtx)
		r.Get("/", rs.Get)       // GET /posts/{id} - Read a single post by :id.
		r.Put("/", rs.Update)    // PUT /posts/{id} - Update a single post by :id.
		r.Delete("/", rs.Delete) // DELETE /posts/{id} - Delete a single post by :id.
	})

	return r
}

// Request Handler - GET /posts - Read a list of posts.
func (rs cowsResource) List(w http.ResponseWriter, r *http.Request) {
	cows := make([]Cow, len(rs.store))

	rs.Lock()
	i := 0
	for _, cow := range rs.store {
		cows[i] = cow
		i++
	}
	sort.Slice(cows, func(i, j int) bool {
		return strings.ToLower(cows[i].Name) < strings.ToLower(cows[j].Name)
	})
	rs.Unlock()

	jsonBytes, err := json.Marshal(cows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Request Handler - POST /posts - Create a new post.
func (rs cowsResource) Create(w http.ResponseWriter, r *http.Request) {

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

	rs.Lock()
	rs.store[cow.ID] = cow
	defer rs.Unlock()
}

func PostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "id", chi.URLParam(r, "id"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Request Handler - GET /posts/{id} - Read a single post by :id.
func (rs cowsResource) Get(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("id").(string)
	found, ok := rs.store[id]

	if !ok {
		fmt.Sprintf("Unknown cow ID: " + id)
	}

	w.Header().Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(found)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Request Handler - PUT /posts/{id} - Update a single post by :id.
func (rs cowsResource) Update(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("id").(string)
	bodyBytes, err := ioutil.ReadAll(r.Body)

	var cow Cow
	err = json.Unmarshal(bodyBytes, &cow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	// TODO: check to see if it exists first
	rs.store[id] = cow
}

// Request Handler - DELETE /posts/{id} - Delete a single post by :id.
func (rs cowsResource) Delete(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("id").(string)

	// TODO: check to see if it exists first
	delete(rs.store, id)
}
