package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func routing(port int) {

	store := HotelStoreJson{}

	r := mux.NewRouter()
	productsRouter := r.PathPrefix("/hotel").Subrouter()
	productsRouter.HandleFunc("", FetchHotelsHandler(&store)).Methods("GET")
	productsRouter.HandleFunc("", BasicAuthMiddleware(CreateHotelsHandler(&store))).Methods("POST")
	productsRouter.HandleFunc("", BasicAuthMiddleware(EditHotelsHandler(&store))).Methods("PUT")
	productsRouter.HandleFunc("/book", BookHotelHandler(&store)).Methods("PUT")
	productsRouter.HandleFunc("", BasicAuthMiddleware(DeleteHotelsHandler(&store))).Methods("DELETE")

	http.ListenAndServe(":"+strconv.Itoa(port), r)
}
