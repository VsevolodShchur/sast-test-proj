package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres
)

type item struct {
	ID    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

func returnItem(w http.ResponseWriter, i item) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&i)
}

type errorResponse struct {
	Error string `json:"error"`
}

func returnError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errResp := errorResponse{Error: err.Error()}
	json.NewEncoder(w).Encode(&errResp)
}

func newGetItemHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		itemId := chi.URLParam(r, "id")
		_, err := strconv.ParseInt(itemId, 10, 64)
		if err != nil {
			returnError(w, http.StatusBadRequest, fmt.Errorf("id param is not int"))
			return
		}

		rows, err := db.Query("SELECT * FROM items WHERE id = " + itemId)
		if err != nil {
			returnError(w, http.StatusInternalServerError, err)
			return
		}

		if !rows.Next() {
			returnError(w, http.StatusNotFound, fmt.Errorf("item with id %s not found", itemId))
			return
		}

		var id int
		var title string
		err = rows.Scan(&id, &title)
		if err != nil {
			returnError(w, http.StatusInternalServerError, err)
			return
		}
		returnItem(w, item{ID: id, Title: title})
	}
}

type config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
	DbSSLMode  string
}

func getConfigFromEnv() config {
	return config{
		DbHost:     os.Getenv("POSTGRES_HOST"),
		DbPort:     os.Getenv("POSTGRES_PORT"),
		DbUser:     os.Getenv("POSTGRES_USER"),
		DbPassword: os.Getenv("POSTGRES_PASSWORD"),
		DbName:     os.Getenv("POSTGRES_DB"),
		DbSSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}

func main() {
	cfg := getConfigFromEnv()
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSSLMode,
		))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := chi.NewRouter()
	router.Get("/items/{id}", newGetItemHandler(db))
	srv := http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	log.Fatalln(srv.ListenAndServe())
}