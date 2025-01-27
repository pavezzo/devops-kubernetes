package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

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

func pongHandler(w http.ResponseWriter, r *http.Request) {
    var num int
    dbconn, err := getConnection()
    if err == nil {
        err = dbconn.QueryRow(context.Background(), "UPDATE pingpong SET num = num + 1 WHERE id=$1 RETURNING num as new_num", "counter").Scan(&num)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Postgresql update error: %s\n", err.Error())
        }
    } else {
        fmt.Fprintf(os.Stderr, "Couldn't get connection to database: %s\n", err.Error())
    }
    fmt.Fprintf(w, "pong %d", num)
}

func pongStatusHandler(w http.ResponseWriter, r *http.Request) {
    var num int
    dbconn, err := getConnection()
    if err == nil {
        err := dbconn.QueryRow(context.Background(), "SELECT num FROM pingpong where id=$1", "counter").Scan(&num)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Postgresql read error: %s\n", err.Error())
        }
    } else {
        fmt.Fprintf(os.Stderr, "Couldn't get connection to database: %s\n", err.Error())
    }
    fmt.Fprintf(w, "%d", num)
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

    conn, err := getConnection()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Couldn't get connection to database: %s\n", err.Error())
        os.Exit(1)
    }
    defer conn.Close(context.Background())

    cmd := `CREATE TABLE IF NOT EXISTS pingpong (
                id varchar(20) NOT NULL,
                num integer NOT NULL,
                PRIMARY KEY (id)
            )`
    _, err = conn.Exec(context.Background(), cmd)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Couldn't create postgresql table: %s\n", err.Error())
        os.Exit(1)
    }

    _, err = conn.Exec(context.Background(), "INSERT INTO pingpong (id, num) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING", "counter", 0)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Couldn't insert into pingpong table: %s\n", err.Error())
        os.Exit(1)
    }

    addr := fmt.Sprintf(":%d", port)
    fmt.Printf("Server started in port %d\n", port)

    http.HandleFunc("/pingpong", pongHandler)
    http.HandleFunc("/pingpongstatus", pongStatusHandler)
    log.Fatal(http.ListenAndServe(addr, nil))
}
