package main

import (
	"fmt"
	"os"
	"time"
)

const folder = "/usr/src/app/files"
const file = folder + "/stamp"

func main() {
    if _, err := os.Stat(file); os.IsNotExist(err) {
        os.MkdirAll(folder, 0644)
    }

    for {
        lastTime := time.Now().UTC().String()
        fmt.Printf("Time is: %s\n", lastTime)
        err := os.WriteFile(file, []byte(lastTime), 0644)
        if err != nil {
            fmt.Printf("Error writing to file: %s\n", err.Error())
        }
        
        time.Sleep(5 * time.Second)
    }
}
