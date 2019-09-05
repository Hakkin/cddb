package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Hakkin/cddb/app/handler"
)

var (
	addr string
	port string
)

func init() {
	flag.StringVar(&addr, "address", "", "The address to listen on")
	flag.StringVar(&port, "port", "8080", "The port to listen on")
}

func main() {
	flag.Parse()

	if addrEnv := os.Getenv("ADDR"); addrEnv != "" {
		addr = addrEnv
	}

	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/cddb", handler.CDDB)
	http.HandleFunc("/cddb/", handler.CDDB)

	log.Printf("Listening on %s:%s\n", addr, port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", addr, port), nil); err != nil {
		log.Fatal("Server exited with error: ", err)
	}
}
