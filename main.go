package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), registerHandlers()))
}

func registerHandlers() *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/members").Handler(appHandler(listHandler))
	return router
}

type appError struct {
	Code int
	Message string
	Error error
}

type appHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r * http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)
		http.Error(w, e.Message, e.Code)
	}
}

type Member struct {
	Name string
}

type Members []Member

func listHandler(w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(Members{
		Member {
			Name: "taro",
		},
	})
	w.Write(response)
	return nil
}

