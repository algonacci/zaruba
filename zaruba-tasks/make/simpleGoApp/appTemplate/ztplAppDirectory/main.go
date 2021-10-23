package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	// get httpPort
	httpPortStr := os.Getenv("APP_HTTP_PORT")
	httpPort, err := strconv.Atoi(httpPortStr)
	if err != nil {
		httpPort = 3000
	}

	// handle URL
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s /", r.Method)
		fmt.Fprintf(w, "Hello world 🐹")
	})

	// serve
	fmt.Printf("Serve HTTP on port %d\n", httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
		log.Fatal(err)
	}

}
