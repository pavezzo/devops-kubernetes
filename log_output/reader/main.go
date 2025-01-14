package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)


const file = "/usr/src/app/files/stamp"

var randomString string


type JsonResponse struct {
    Time string
    String string
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    _ = r
    time, err := os.ReadFile(file)
    if err != nil {
        fmt.Printf("Error reading file: %s\n", err.Error())
    }
    data := JsonResponse{ Time: string(time), String: randomString }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(data)
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
