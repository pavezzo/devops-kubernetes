package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)


const stampfile = "/usr/src/app/files/stamp"
const informationfile = "/usr/src/app/information/information.txt"
const counterAddr = "http://ping-pong-svc:2346/pingpongstatus"

var randomString string
var envMessage string
var informationContents string

func statusHandler(w http.ResponseWriter, r *http.Request) {
    _ = r

    time, err := os.ReadFile(stampfile)
    if err != nil {
        fmt.Printf("Error reading file: %s\n", err.Error())
    }

    counter := -1
    resp, err := http.Get(counterAddr)
    if err != nil {
        fmt.Printf("Error getting pong status: %s\n", err.Error())
    } else {
        data, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("Error reading pong status response body: %s\n", err.Error())
        } else {
            num, err := strconv.Atoi(string(data))
            if err == nil {
                counter = num
            }
        }
    }

    fmt.Fprintf(w, "file content: %s\nenv variable MESSAGE=%s\n%s: %s\nPing / Pongs: %d\n", informationContents, envMessage, time, randomString, counter)
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

    envMessage, ok = os.LookupEnv("MESSAGE")
    if !ok {
        fmt.Printf("Couldn't find MESSAGE from environment variables\n")
    }

    informationBytes, err := os.ReadFile(informationfile)
    if err != nil {
        fmt.Printf("Couldn't read information.txt: %s\n", err.Error())
    } else {
        informationContents = string(informationBytes)
    }


    addr := fmt.Sprintf(":%d", port)
    fmt.Printf("Server started in port %d\n", port)
    http.HandleFunc("/status", statusHandler)

    log.Fatal(http.ListenAndServe(addr, nil))
}
