package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type TodoPostReq struct {
    Todo string
}

var todos []string

func todoHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        w.Header().Set("Content-Type", "application/json")
        err := json.NewEncoder(w).Encode(todos)
        if err != nil {
            fmt.Printf("Error responding to get request: %s\n", err.Error())
        }
    case http.MethodPost:
        err := r.ParseForm()
        if err != nil {
            fmt.Printf("Error parsing post request form: %s\n", err.Error())
            return
        }

        todo := r.PostForm.Get("todo")
        if len(todo) == 0 {
            fmt.Printf("No value for todo\n")
            http.Error(w, "No value for todo", http.StatusBadRequest)
            return
        } else if len(todo) > 140 {
            http.Error(w, "Too long todo", http.StatusBadRequest)
            return
        }

        todos = append(todos, todo)
        http.Redirect(w, r, r.Header.Get("Referer"), 302)
    default:
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
    }
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

    http.HandleFunc("/todos", todoHandler)
    log.Fatal(http.ListenAndServe(addr, nil))
}
