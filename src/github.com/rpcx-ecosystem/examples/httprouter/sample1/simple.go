package main

import (
"fmt"
"log"
"net/http"

"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!中国\n")
}

func main() {
	router := httprouter.New()
	//Handle签名为
	// type Handle func(http.ResponseWriter, *http.Request, Params)
	router.GET("/", Index)

	log.Fatal(http.ListenAndServe(":8080", router))
}

