package endpoint

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RegisterHandlers() *mux.Router {
	router := mux.NewRouter()
	router.Methods("GET").Path("/members").Handler(appHandler(listHandler))
	router.Methods("POST").Path("/members").Handler(appHandler(createHandler))
	router.Methods("PUT").Path("/members/{id:[0-9]+}").Handler(appHandler(updateHandler))
	router.Methods("DELETE").Path("/members/{id:[0-9]+}").Handler(appHandler(deleteHandler))

	router.Methods("GET").Path("/orders").Handler(appHandler(listOrdersHandler))
	router.Methods("POST").Path("/orders").Handler(appHandler(createOrderHandler))
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

func appErrorFormat(error error, format string, v interface{}) *appError {
	return &appError{
		Code:    500,
		Message: fmt.Sprintf(format, v),
		Error:   error,
	}
}

