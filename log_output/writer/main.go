package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
    str := uuid.New().String()
    for {
        t := time.Now().UTC().String()
        fmt.Printf("%s: %s\n", t, str)
        time.Sleep(5 * time.Second)
    }
}
