package main

import (
	"encoding/json"
	"errors"
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
	router.Methods("POST").Path("/members").Handler(appHandler(createHandler))
	router.Methods("PUT").Path("/members/{id:[0-9]+}").Handler(appHandler(updateHandler))
	router.Methods("DELETE").Path("/members/{id:[0-9]+}").Handler(appHandler(deleteHandler))
	return router
}

type appError struct {
	Code    int
	Message string
	Error   error
}

type appHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	response, jsonError := json.Marshal(Members{
		Member{
			Name: "Name1",
		},
		Member{
			Name: "Name2",
		},
	})
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	return nil
}

func createHandler(w http.ResponseWriter, r *http.Request) *appError {
	response, jsonError := json.Marshal(Members{
		Member{
			Name: "new name",
		},
	})
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/members/2")
	w.WriteHeader(201)
	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	return nil
}

func updateHandler(w http.ResponseWriter, r *http.Request) *appError {
	response, jsonError := json.Marshal(Members{
		Member{
			Name: "updated Name " + mux.Vars(r)["id"],
		},
	})
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	return nil
}

func deleteHandler(w http.ResponseWriter, r *http.Request) *appError {
	requestId := mux.Vars(r)["id"]
	if requestId == "2" {
		return appErrorFormat(errors.New("invalid reques id"), "invalid id:  %s", requestId)
	}
	w.Header().Set("Content-Type", "application/json")
	return nil
}

func appErrorFormat(error error, format string, v interface{}) *appError {
	return &appError{
		Code:    500,
		Message: fmt.Sprintf(format, v),
		Error:   error,
	}
}
