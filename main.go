package main

import (
    "log"
    "net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
    const filepathRoot = "."
    const port = "8080"

    apiCfg := apiConfig {
        fileserverHits: 0,
    }

    mux := http.NewServeMux()
    fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
    mux.Handle("/app/", fsHandler)

    mux.HandleFunc("GET /api/healthz", handlerReadiness)
    mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

    srv := &http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files from %s onport: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())
}
