package main

import (
	"Forum_Greyddit/backend"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", backend.HomeHandler)
	http.HandleFunc("/register", backend.RegisterHandler)
	http.HandleFunc("/post/create", backend.PostCreateHandler)
	log.Println("Bonjour et bienvenus sur Greyddit: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
