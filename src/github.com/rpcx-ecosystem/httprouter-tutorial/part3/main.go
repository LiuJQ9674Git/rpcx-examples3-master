package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	part "github.com/rpcx-ecosystem/httprouter-tutorial/part3/part3"
)

func main() {
	router := httprouter.New()
	router.GET("/", part.Index)
	router.GET("/books", part.BookIndex)
	router.GET("/books/:isdn", part.BookShow)

	//// Create a couple of sample Book entries
	//part.bookstore["123"] = &part.Book{
	//	ISDN:   "123",
	//	Title:  "Silence of the Lambs",
	//	Author: "Thomas Harris",
	//	Pages:  367,
	//}
	//
	//part.bookstore["124"] = &part.Book{
	//	ISDN:   "124",
	//	Title:  "To Kill a Mocking Bird",
	//	Author: "Harper Lee",
	//	Pages:  320,
	//}

	log.Fatal(http.ListenAndServe(":8080", router))
}
