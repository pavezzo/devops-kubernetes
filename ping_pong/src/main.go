package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var counter int

func pongHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "pong %d", counter)
    counter += 1
}


func main() {
    counter = 0

    port := 8000
    portstr, ok := os.LookupEnv("PORT")
    if ok {
        envport, err := strconv.Atoi(portstr)
        if err == nil {
            port = envport
        }
    }
    addr := fmt.Sprintf(":%d", port)
    fmt.Printf("Server started in port %d\n", port)

    http.HandleFunc("/pingpong", pongHandler)
    log.Fatal(http.ListenAndServe(addr, nil))
}
