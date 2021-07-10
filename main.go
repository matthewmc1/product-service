package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	mux "github.com/gorilla/mux"
)

type Product struct {
	ID         int      `json :"id"`
	BRAND      string   `json : "brand"`
	CATEGORIES []string `json : "categories"`
	PRICE      float64  `json: "price"`
	QUANTITY   int      `json: "quantity"`
}

func main() {

	l := log.New(os.Stdout, "product-service-api", 3)

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	l.Printf("Service has stared on process id %v and host %v", os.Getpid(), name)

	router := mux.NewRouter()
	router.HandleFunc("/", EmptyResponseHandler).Methods(http.MethodGet)
	router.HandleFunc("/products", GetAllProductsHandler).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":6743",
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	http.Handle("/", router)

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt)
	signal.Notify(sigChannel, os.Kill)

	sig := <-sigChannel
	log.Panicln("Receieved termination call", sig)

	tc, _ := context.WithTimeout(context.Background(), 3*time.Second)
	srv.Shutdown(tc)

}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	traceId := rand.Int()

	l := log.New(os.Stdout, "product-handler", 3)
	l.Printf("Requst on product handler at %v and trace ID %d", time.Now().UTC(), traceId)

	pl := ProductList

	if len(pl) < 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	json.NewEncoder(w).Encode(pl)
}

func EmptyResponseHandler(w http.ResponseWriter, r *http.Request) {
	l := log.New(os.Stdout, "empty-product-handler", 3)
	l.Printf("Request to empty product handler at %v", time.Now().UTC())

	w.WriteHeader(http.StatusNoContent)
}

var ProductList = []*Product{
	{
		ID:    1,
		BRAND: "Levi",
		CATEGORIES: []string{
			"jeans", "adult", "mens",
		},
		PRICE:    79.99,
		QUANTITY: 10,
	},
}