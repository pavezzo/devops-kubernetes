package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var lastTime string
var randomString string

type JsonResponse struct {
    Time string
    String string
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    _ = r
    data := JsonResponse{ Time: lastTime, String: randomString }
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
    go http.ListenAndServe(addr, nil)

    for {
        lastTime = time.Now().UTC().String()
        fmt.Printf("%s: %s\n", lastTime, randomString)
        time.Sleep(5 * time.Second)
    }
}
