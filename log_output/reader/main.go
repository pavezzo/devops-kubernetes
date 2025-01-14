package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)


const stampfile = "/usr/src/app/files/stamp"
const pongfile = "/usr/src/app/files/pong_counter"

var randomString string

func statusHandler(w http.ResponseWriter, r *http.Request) {
    _ = r

    time, err := os.ReadFile(stampfile)
    if err != nil {
        fmt.Printf("Error reading file: %s\n", err.Error())
    }

    counter, err := os.ReadFile(pongfile)
    if err != nil {
        fmt.Printf("Error reading file: %s\n", err.Error())
    }

    fmt.Fprintf(w, "%s: %s\nPing / Pongs: %s\n", time, randomString, counter)
}

func main() {
    randomString = uuid.New().String()

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
    http.HandleFunc("/status", statusHandler)

    log.Fatal(http.ListenAndServe(addr, nil))
}
