package bookapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	_logger "github.com/putto112620002/go-http/logger"
)

var ()

type BookApi struct {
	server  http.Server
	logger  _logger.Logger
	addr    string
	storage BookStorage
}

func NewBookApi() *BookApi {
	return &BookApi{
		logger:  *_logger.NewLogger(true),
		addr:    ":8080",
		storage: &BookMemoryStorage{books: []Book{}},
	}
}

func (api *BookApi) handleGetBooks(w http.ResponseWriter, r *http.Request) {

	books := api.storage.GetBooks()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(books)

}

func (api *BookApi) handleGetBookByISBN(w http.ResponseWriter, r *http.Request) {


	isbn := strings.TrimPrefix(r.URL.Path, "books/")
	api.logger.Info("retreiving book with ISBN: %s", isbn)
	book := api.storage.GetBookByISBN(isbn)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)

}

func (api *BookApi) handleCreateBook(w http.ResponseWriter, r *http.Request) {

	
	var newBook Book
	json.NewDecoder(r.Body).Decode(&newBook)
	defer r.Body.Close()
	if newBook.ISBN == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if newBook.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	api.logger.Info("adding book with ISBN: %s", newBook.ISBN)
	start := time.Now()
	api.storage.AddBook(newBook)
	api.logger.Info("added book took: %v seconds", time.Since(start).Seconds())
	w.WriteHeader(http.StatusCreated)

}

func (api *BookApi) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.handleCreateBook(w, r)
			return
		}

		if r.Method == http.MethodGet {
			api.handleGetBooks(w, r)
			return
		}

	})
	
	mux.HandleFunc("/books/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.handleGetBookByISBN(w, r)
			return
		}
	})

	api.server = http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		api.logger.Info("server is running on port %d", 8080)
		if err := api.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			api.logger.Error("error starting server: %v", err)
		}

	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	<-sig

	api.logger.Info("server shutting down...")
	if err := api.server.Shutdown(ctx); err != nil {
		api.logger.Error("error shutting down server: %v", err)
	}
	api.logger.Info("server gracefully shut down")

	defer func() {
		api.logger.Info("cleaning up extra resources")
		cancel()
	}()

}
