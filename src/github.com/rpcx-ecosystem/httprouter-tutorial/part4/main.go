package main

import (
	"log"
	"net/http"
	part "github.com/rpcx-ecosystem/httprouter-tutorial/part4/part4"
)

func main() {
	router := part.NewRouter(part.AllRoutes())
	log.Fatal(http.ListenAndServe(":8080", router))
}
