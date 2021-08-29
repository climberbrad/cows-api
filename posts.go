package main

import (
"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Cow struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Date   string `json:"date"`
	Image  string `json:"image"`
	Finder string `json:"finder"`
}

type postsResource struct{
	sync.Mutex
	store map[string]Cow
}

func newHandler() *postsResource {
	return &postsResource{
		store: map[string]Cow{},
	}
}

func (rs postsResource) Routes() chi.Router {
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
func (rs postsResource) List(w http.ResponseWriter, r *http.Request) {

	cows := make([]Cow, len(rs.store))

	rs.Lock()
	i := 0
	for _, cow := range rs.store {
		cows[i] = cow
		i++
	}
	rs.Unlock()

	w.Header().Set("Content-Type", "application/json")

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
func (rs postsResource) Create(w http.ResponseWriter, r *http.Request) {

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
func (rs postsResource) Get(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("id").(string)
	found := rs.store[id]

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
func (rs postsResource) Update(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("id").(string)
	bodyBytes, err := ioutil.ReadAll(r.Body)

	fmt.Println("ID %s", id)

	var cow Cow
	err = json.Unmarshal(bodyBytes, &cow)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	rs.store[id] = cow
}

// Request Handler - DELETE /posts/{id} - Delete a single post by :id.
func (rs postsResource) Delete(w http.ResponseWriter, r *http.Request) {

	id := r.Context().Value("id").(string)
	delete(rs.store, id)
}
