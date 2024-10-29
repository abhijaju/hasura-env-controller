package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

func getEnvVars(w http.ResponseWriter, r *http.Request) {
	envVars := make(map[string]string)
	for _, kv := range os.Environ() {
		pair := strings.SplitN(kv, "=", 2)
		envVars[pair[0]] = pair[1]
	}

	b, err := json.Marshal(envVars)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(b))
}

func main() {
	// setting up logging
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info("Starting the httpserver...")
	r := mux.NewRouter()
	r.HandleFunc("/", getEnvVars)
	http.ListenAndServe(":9090", r)
}
