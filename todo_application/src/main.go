package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const folder = "/usr/src/app/files/"
const imageFile = folder + "image.jpg"

func indexHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "public/index.html")
}

func fileNeedsUpdate(t time.Duration) bool {
    info, err := os.Stat(imageFile)
    if err != nil {
        fmt.Printf("Couldn't read file info: %s\n", err.Error())
        return true
    }

    now := time.Now()
    if now.Sub(info.ModTime()) >= t {
        return true
    }

    return false
}

func getNewImage() error {
    if _, err := os.Stat(imageFile); os.IsNotExist(err) {
        os.MkdirAll(folder, 0644)
    }

    resp, err := http.Get("https://picsum.photos/1000")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    err = os.WriteFile(imageFile, data, 0644)
    if err != nil {
        return err
    }

    return nil
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

    go func() {
        for {
            if fileNeedsUpdate(time.Minute * 60) {
                if err := getNewImage(); err != nil {
                    fmt.Printf("Error getting new image %s\n", err)
                }
            }
            time.Sleep(time.Second * 1)
        }
    }()

    http.HandleFunc("/", indexHandler)
    http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))
    log.Fatal(http.ListenAndServe(addr, nil))
}
