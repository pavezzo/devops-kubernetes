package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const folder = "/usr/src/app/files"
const file = folder + "/pong_counter"

var counter int

func pongHandler(w http.ResponseWriter, r *http.Request) {
    writeToFile()
    fmt.Fprintf(w, "pong %d", counter)
    counter += 1
}

func writeToFile() {
    str := strconv.Itoa(counter)
    err := os.WriteFile(file, []byte(str), 0644)
    if err != nil {
        fmt.Printf("Error writing to file: %s\n", err.Error())
    }
}


func main() {
    if _, err := os.Stat(file); os.IsNotExist(err) {
        os.MkdirAll(folder, 0644)
    }

    counterstr, err := os.ReadFile(file)
    if err != nil {
        counter = 0
        writeToFile()
    } else {
        counter, _ = strconv.Atoi(string(counterstr))
    }

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
