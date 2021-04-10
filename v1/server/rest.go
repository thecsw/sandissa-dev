package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	RESTClientCert string = "./rest/cert1.pem"
	RESTClientKey  string = "./rest/privkey1.pem"
)

// startRouter starts the mux router and blocks until a crash or
// a SIGINT signal.
func startRouter() {
	r := mux.NewRouter()
	// Standard GET methods to retrieve blogs and pages
	r.HandleFunc("/", rootPage).Methods(http.MethodGet)
	r.HandleFunc("/temp", tempHandler).Methods(http.MethodGet)
	r.HandleFunc("/temps", tempsHandler).Methods(http.MethodGet)
	r.HandleFunc("/led", ledHandler).Methods(http.MethodPost)

	r.Use(loggingMiddleware)

	// Declare and define our HTTP handler
	handler := cors.Default().Handler(r)
	srv := &http.Server{
		Handler: handler,
		Addr:    ":5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	// Fire up the router
	go func() {
		if err := srv.ListenAndServeTLS(RESTClientCert, RESTClientKey); err != nil {
			lerr("Failed to fire up the router", err, params{})
		}
	}()
	// Listen to SIGINT and other shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// rootPage is a generic placeholder HTTP handler
func rootPage(w http.ResponseWriter, r *http.Request) {
	httpHTML(w, "hello, world")
}

func tempHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := getTempDB()
	if err != nil {
		httpJSON(w, nil, http.StatusInternalServerError, err)
		return
	}
	httpJSON(w, httpMessageReturn{Message: temp.Value}, http.StatusOK, nil)
}

func tempsHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := getTempsDB()
	if err != nil {
		httpJSON(w, nil, http.StatusInternalServerError, err)
		return
	}
	toReturn := make([]float64, len(temp))
	for i, v := range temp {
		toReturn[i] = v.Value
	}
	httpJSON(w, httpMessageReturn{Message: toReturn}, http.StatusOK, nil)
}

func ledHandler(w http.ResponseWriter, r *http.Request) {
	request := &ledStatusRequset{}
	json.NewDecoder(r.Body).Decode(request)
	toSend := "off"
	if request.Status {
		toSend = "on"
	}
	publish(1, topicLED, toSend)
	httpJSON(w, httpMessageReturn{Message: "OK"}, http.StatusOK, nil)
}

// httpJSON is a generic http object passer.
func httpJSON(w http.ResponseWriter, data interface{}, status int, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	if err != nil && status >= 400 && status < 600 {
		json.NewEncoder(w).Encode(httpErrorReturn{Error: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(data)
}

// httpHTML sends a good HTML response
func httpHTML(w http.ResponseWriter, data interface{}) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, data)
}

// ledStatusRequset is the POST body for LED
type ledStatusRequset struct {
	Status bool `json:"status"`
}

// httpMessageReturn defines a generic HTTP return message.
type httpMessageReturn struct {
	Message interface{} `json:"message"`
}

// httpErrorReturn defines a generic HTTP error message.
type httpErrorReturn struct {
	Error string `json:"error"`
}
