package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Servidor local recebeu:", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `<html><body><h1>vaai</h1></body></html>`)
	})

	log.Println("Servidor local escutando em http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
