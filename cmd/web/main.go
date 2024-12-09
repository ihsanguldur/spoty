package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/test", testHandler)

	fmt.Println("server is listening on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error while starting server: \n%v", err)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test"))
}
