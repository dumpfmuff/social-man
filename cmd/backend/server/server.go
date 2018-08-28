package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/osimono/social-man/mongo"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"strings"
	"time"
)

func Init(protocol, host string, port int) {
	r := mux.NewRouter()

	bundle := clientopt.ClientBundle{}
	timeout := time.Duration(2) * time.Second
	client, _ := mgo.NewClientWithOptions("mongodb://@localhost:27017",
		bundle.ServerSelectionTimeout(timeout),
		bundle.ConnectTimeout(timeout),
		bundle.SocketTimeout(timeout))

	connectErr := client.Connect(context.Background())
	if connectErr != nil {
		panic("no mongo")
	}
	t := mongo.NewStorage(client)

	api := r.PathPrefix("/api/").Subrouter()
	api.Handle("/tenants", &AllTenantHandler{&t}).Methods(http.MethodGet)
	api.Handle("/tenants", &NewTenantHandler{&t}).Methods(http.MethodPost)
	http.Handle("/", r)

	// Serve static assets directly.
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("web/ui/build")))

	// Catch-all: Serve our JavaScript application's entry-point (index.html).
	r.PathPrefix("/").HandlerFunc(indexHandler("web/ui/build/index.html"))

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

func indexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
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
