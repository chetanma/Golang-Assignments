package main

import (
	"net/http"	
	"github.com/gorilla/mux"
)

func routing(){

	sources :=make([]string,0,100)
	sources =append(sources, "http://localhost:8001")
	sources =append(sources, "http://localhost:8002")

	r := mux.NewRouter()
	productsRouter := r.PathPrefix("/hotel").Subrouter()
	productsRouter.HandleFunc("", FetchHotelsHandler(sources)).Methods("GET")
	productsRouter.HandleFunc("/book", BookHotelHandler()).Methods("PUT")
	
	http.ListenAndServe(":8003", r)	
}
