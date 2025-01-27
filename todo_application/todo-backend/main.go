package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type TodoPostReq struct {
    Todo string
}

var conn *pgx.Conn

func getConnection() (*pgx.Conn, error) {
    if conn == nil || conn.IsClosed() {
        dbport := os.Getenv("DB_PORT")
        dbpass := os.Getenv("DB_PASSWORD")
        dbhost := os.Getenv("DB_HOST")
        dbuser := os.Getenv("DB_USER")
        dbname := os.Getenv("DB_NAME")
        connstr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", dbuser, dbpass, dbhost, dbport, dbname)

        var err error
        conn, err = pgx.Connect(context.Background(), connstr)
        if err != nil {
            return nil, err
        }
    }

    return conn, nil
}


func getTodos() ([]string, error) {
    var todos []string

    dbconn, err := getConnection()
    if err != nil {
        return todos, err
    }

    rows, err := dbconn.Query(context.Background(), "SELECT todo FROM todos")
    if err != nil {
        return todos, err
    }
    defer rows.Close()

    for rows.Next() {
        var todo string
        err := rows.Scan(&todo)
        if err != nil {
            return todos, err
        }

        todos = append(todos, todo)
    }

    return todos, nil
}

func addTodo(todo string) error {
    if len(todo) > 140 {
        msg := fmt.Sprintf("\"%s\" is too long", todo)
        return errors.New(msg)
    } else if len(todo) == 0 {
        return errors.New("todo is empty")
    }

    dbconn, err := getConnection()
    if err != nil {
        return err
    }

    _, err = dbconn.Exec(context.Background(), "INSERT INTO todos (todo) VALUES ($1)", todo)
    if err != nil {
        return err
    }
    fmt.Printf("Added todo with following contents: \"%s\"\n", todo)

    return nil
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        w.Header().Set("Content-Type", "application/json")
        todos, err := getTodos()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error reading todos from database: %s\n", err.Error())
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            return
        }

        err = json.NewEncoder(w).Encode(todos)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error responding to get request: %s\n", err.Error())
        }
    case http.MethodPost:
        err := r.ParseForm()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error parsing post request form: %s\n", err.Error())
            return
        }

        todo := r.PostForm.Get("todo")
        err = addTodo(todo)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error adding todo: %s\n", err.Error())
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

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

    fmt.Printf("Initializing postgres connection\n")
    conn, err := getConnection()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to connect to database: %s\n", err.Error())
        os.Exit(1)
    }
    defer conn.Close(context.Background())
    fmt.Printf("Succesfully opened connection to database\n")

    cmd := `CREATE TABLE IF NOT EXISTS todos (
                id SERIAL PRIMARY KEY,
                todo varchar(140) NOT NULL
            )`
    _, err = conn.Exec(context.Background(), cmd)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Couldn't create postgresql table: %s\n", err.Error())
        os.Exit(1)
    }

    addr := fmt.Sprintf(":%d", port)
    fmt.Printf("Server started in port %d\n", port)

    http.HandleFunc("/todos", todoHandler)
    log.Fatal(http.ListenAndServe(addr, nil))
}
