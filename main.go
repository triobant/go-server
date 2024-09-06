package main

import (
    "log"
    "net/http"
)

type Dir string

func main() {
    const filepathRoot = "."
    const port = "8080"

    mux := http.NewServeMux()
    mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    mux.HandleFunc("/healthz", handlerReadiness)

    srv := &http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files from %s onport: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())
}


func handlerReadiness(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(http.StatusText(http.StatusOK)))
}
