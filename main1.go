package main

import (
    "fmt"
    "net/http"
    "github.com/goombaio/namegenerator"
    "time"
)

func main() {
    http.HandleFunc("/", help)
    http.ListenAndServe("localhost:8089", nil)
    
}

func help(writer http.ResponseWriter, request *http.Request) {
    seed := time.Now().UTC().Unix()
    nameG := namegenerator.NewNameGenerator(seed)
    fmt.Fprintf(writer, "Hello %s", nameG.Generate())
}
