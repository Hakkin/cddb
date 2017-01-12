// +build !appengine

package main

import (
	"github.com/hakkin/cddb"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/cddb", cddb.CddbHttp)
	http.HandleFunc("/cddb/", cddb.CddbHttp)
	http.ListenAndServe(":8080", nil)
}
