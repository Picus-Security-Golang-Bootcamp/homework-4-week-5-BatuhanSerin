package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	postgres "github.com/BatuhanSerin/postgresql/common/db"
	"github.com/BatuhanSerin/postgresql/domain/author"
	"github.com/BatuhanSerin/postgresql/domain/book"
	httpErrors "github.com/BatuhanSerin/postgresql/server/http_errors"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func Server() {
	r := mux.NewRouter()

	handlers.AllowedOrigins([]string{"https://www.example.com"})
	handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	handlers.AllowedMethods([]string{"POST", "GET", "PUT", "PATCH"})

	//0.0.0.0:8090/products/bb
	r.Use(loggingMiddleware)
	r.Use(authenticationMiddleware)

	s := r.PathPrefix("/list").Subrouter()
	//"/list/{name}/"
	s.HandleFunc("", BookList)
	s.HandleFunc("/{id}", BookListById)
	s.HandleFunc("/id/{id}", BookListByAuthorOrBookId)
	s.HandleFunc("/{name}", BookListByName)

	p := r.PathPrefix("/user").Subrouter()
	p.HandleFunc("/", userCreate).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:         "localhost:8090",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	ShutdownServer(srv, time.Second*10)
}

// type ApiResponse struct {
// 	Data interface{} `json:"data"`
// }

func BookRepo() *book.BookRepository {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := postgres.NewPsqlDB()
	if err != nil {
		log.Fatal("Postgres cannot init ", err)
	}

	log.Println("Postgres connected")

	bookRepo := book.NewBookRepository(db)
	bookRepo.Migrations()
	bookRepo.InsertData()

	return bookRepo
}
func AuthorRepo() *author.AuthorRepository {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := postgres.NewPsqlDB()
	if err != nil {
		log.Fatal("Postgres cannot init ", err)
	}

	log.Println("Postgres connected")

	authorRepo := author.NewAuthorRepository(db)
	authorRepo.Migrations()
	authorRepo.InsertData()
	return authorRepo
}

func BookList(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//r.URL.Query().Get("param")

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	d := BookRepo().FinAll()

	resp, _ := json.Marshal(d)
	w.Write(resp)
}

func BookListById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	d := BookRepo().FindBookById(id)

	resp, _ := json.Marshal(d)
	w.Write(resp)
}

func BookListByAuthorOrBookId(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	d := BookRepo().FindByAuthorOrBookId(id)

	resp, _ := json.Marshal(d)
	w.Write(resp)
}

func BookListByName(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// id, _ := strconv.Atoi(vars["name"])

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")

	d := BookRepo().FindByName(vars["name"])

	resp, _ := json.Marshal(d)
	w.Write(resp)
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// type errorsResponse struct {
// 	message string `json:"message"`
// }

func userCreate(w http.ResponseWriter, r *http.Request) {
	var u User

	if r.Header.Get("Content-Type") != "application/json" {
		err := httpErrors.ParseErrors(httpErrors.NotAllowedImageHeader)
		w.Write([]byte(err.Error()))
		return
	}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.Write([]byte(httpErrors.
			ParseErrors(httpErrors.BadRequest).
			Error()))
		return
	}

	personData, err := json.Marshal(u)
	if err != nil {
		w.Write([]byte(httpErrors.
			ParseErrors(httpErrors.BadRequest).
			Error()))
		return
	}
	w.Write(personData)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if strings.HasPrefix(r.URL.Path, "/products/") {
			if token != "" {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Token not found", http.StatusUnauthorized)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

//https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97
//https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/
func ShutdownServer(srv *http.Server, timeout time.Duration) {
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
