package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
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
    mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

    mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)


    srv := &http.Server {
        Addr: ":" + port,
        Handler: mux,
    }

    log.Printf("Serving files from %s onport: %s\n", filepathRoot, port)
    log.Fatal(srv.ListenAndServe())
}

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
    type parameters struct {
        Body string `json:"body"`
    }
    type returnVals struct {
        CleanedBody string `json:"cleaned_body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
        return
    }

    const maxChirpLength = 140
    if len(params.Body) > maxChirpLength {
        respondWithError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }

    badWords := map[string]struct{}{
        "kerfuffle": {},
        "sharbert": {},
        "fornax": {},
    }
    cleaned := getCleanedBody(params.Body, badWords)

    respondWithJSON(w, http.StatusOK, returnVals{
        CleanedBody: cleaned,
    })
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
    words := strings.Split(body, " ")
    for i, word := range words {
        loweredWord := strings.ToLower(word)
        if _, ok := badWords[loweredWord]; ok {
            words[i] = "****"
        }
    }
    cleaned := strings.Join(words, " ")
    return cleaned
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
    if code > 499 {
        log.Printf("Responding with 5XX error: %s", msg)
    }
    type errorResponse struct {
        Error string `json:"error"`
    }
    respondWithJSON(w, code, errorResponse{
        Error: msg,
    })
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(500)
        return
    }
    w.WriteHeader(code)
    w.Write(dat)
}
