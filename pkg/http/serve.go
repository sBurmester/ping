package http

import (
	"net/http"
	"time"
)

// handlerPing handles the /ping endpoint, responding with "pong".
func handlerPing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func addHandler(router *http.ServeMux) {
	router.HandleFunc("/ping", handlerPing)
	router.Handle("/livez", Live)
	router.Handle("/readyz", Ready)

}

func StartServer() error {
	router := http.NewServeMux()
	addHandler(router)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return server.ListenAndServe()
}
