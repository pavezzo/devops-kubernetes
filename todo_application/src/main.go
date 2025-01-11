package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "public/index.html")
}


func main() {
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

    http.HandleFunc("/", indexHandler)
    log.Fatal(http.ListenAndServe(addr, nil))
}
