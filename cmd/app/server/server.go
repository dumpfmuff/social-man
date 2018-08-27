package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/osimono/social-man/cmd/app/mongo"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"strings"
	"time"
)

func Init(protocol, host string, port int) {
	r := mux.NewRouter()
	r.HandleFunc("/clients", mongo.AllClients).Methods(http.MethodGet)
	r.HandleFunc("/clients", mongo.NewTenant).Methods(http.MethodPost)
	http.Handle("/", r)

	serverAddress := fmt.Sprintf("%v:%v", host, port)
	fmt.Println("Starting up on " + protocol + "://" + serverAddress)

	srv := &http.Server{
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		// Good practice: enforce timeouts for servers you create!
	}

	http2.ConfigureServer(srv, nil)

	if strings.EqualFold(protocol, "https") {
		log.Fatal(http.ListenAndServeTLS(serverAddress, "tls/server.crt", "tls/server.key", LogWrapper(r)))
	} else {
		log.Fatal(http.ListenAndServe(serverAddress, LogWrapper(r)))
	}

}

func LogWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loggingResponseWriter := newLoggingResponseWriter(w)
		h.ServeHTTP(loggingResponseWriter, r)
		log.Printf("%v - %v - %v", r.Method, r.URL, loggingResponseWriter.statusCode)
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
