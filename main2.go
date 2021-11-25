package main

import (
    "fmt"
    "net/http"
    "log"
)

func main() {
    http.HandleFunc("/", help)
    log.Fatal(http.ListenAndServe("localhost:8089", nil))

}

func help(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprintf(writer, "Hello %s", request.URL.Path[1:])
}
