// +build !appengine

package main

import (
	"net/http"

	"github.com/Hakkin/cddb/app/handler"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/cddb", handler.CDDB)
	http.HandleFunc("/cddb/", handler.CDDB)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
