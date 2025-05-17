package http

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const name = "go.opentelemetry.io/otel/homelab/ping"

var (
	tracer  = otel.Tracer(name)
	meter   = otel.Meter(name)
	pingCnt metric.Int64Counter
)

func init() {
	var err error
	pingCnt, err = meter.Int64Counter("ping.request.count",
		metric.WithDescription("The number of ping requests"),
		metric.WithUnit("{ping}"))
	if err != nil {
		panic(err)
	}
}

// handlerPing handles the /ping endpoint, responding with "pong".
func handlerPing(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	pingCnt.Add(ctx, 1)

	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("pong"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func addHandler(router *http.ServeMux) {
	router.HandleFunc("GET /ping", handlerPing)
	router.Handle("GET /livez", Live)
	router.Handle("GET /readyz", Ready)
	router.Handle("GET /metrics", promhttp.Handler())

}

func StartServer() error {
	Live.Up()
	router := http.NewServeMux()
	addHandler(router)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Add more timeouts for better security and reliability
	server.WriteTimeout = 10 * time.Second
	server.IdleTimeout = 60 * time.Second

	Ready.Up()

	return server.ListenAndServe()
}
