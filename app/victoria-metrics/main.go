package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	// httpListenAddr is the address to listen for HTTP connections.
	httpListenAddr = flag.String("httpListenAddr", ":8428", "TCP address to listen for HTTP connections")

	// retentionPeriod defines how long to keep data in months.
	retentionPeriod = flag.Int("retentionPeriod", 1, "Retention period in months")

	// storageDataPath is the path to the directory where VictoriaMetrics stores its data.
	storageDataPath = flag.String("storageDataPath", "victoria-metrics-data", "Path to storage data directory")

	// maxInsertRequestSize is the maximum size of a single insert request in bytes.
	maxInsertRequestSize = flag.Int("maxInsertRequestSize", 32*1024*1024, "The maximum size in bytes of a single insert request")

	// loggerLevel defines the logging level.
	loggerLevel = flag.String("loggerLevel", "INFO", "Minimum level of errors to log. Possible values: INFO, WARN, ERROR, FATAL, PANIC")
)

func main() {
	// Parse command-line flags.
	flag.Parse()

	// Validate retention period.
	if *retentionPeriod < 1 {
		fmt.Fprintf(os.Stderr, "retentionPeriod must be at least 1 month; got %d\n", *retentionPeriod)
		os.Exit(1)
	}

	fmt.Printf("Starting VictoriaMetrics\n")
	fmt.Printf("  httpListenAddr: %s\n", *httpListenAddr)
	fmt.Printf("  retentionPeriod: %d months\n", *retentionPeriod)
	fmt.Printf("  storageDataPath: %s\n", *storageDataPath)
	fmt.Printf("  loggerLevel: %s\n", *loggerLevel)

	// Ensure the storage data directory exists.
	if err := os.MkdirAll(*storageDataPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "cannot create storageDataPath %q: %v\n", *storageDataPath, err)
		os.Exit(1)
	}

	// Set up HTTP server with routing.
	s := &fasthttp.Server{
		Handler:            requestHandler,
		Name:               "VictoriaMetrics",
		ReadTimeout:        60 * time.Second,
		WriteTimeout:       60 * time.Second,
		MaxRequestBodySize: *maxInsertRequestSize,
	}

	fmt.Printf("Listening for HTTP connections at %s\n", *httpListenAddr)
	if err := s.ListenAndServe(*httpListenAddr); err != nil {
		fmt.Fprintf(os.Stderr, "cannot start HTTP server at %q: %v\n", *httpListenAddr, err)
		os.Exit(1)
	}
}

// requestHandler routes incoming HTTP requests to appropriate handlers.
func requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())

	switch path {
	case "/health":
		// Health check endpoint.
		ctx.SetStatusCode(http.StatusOK)
		fmt.Fprintf(ctx, "OK")

	case "/metrics":
		// Prometheus-compatible metrics endpoint.
		handleMetrics(ctx)

	case "/api/v1/write":
		// Prometheus remote write endpoint.
		handleRemoteWrite(ctx)

	case "/api/v1/query":
		// Prometheus query endpoint.
		handleQuery(ctx)

	case "/api/v1/query_range":
		// Prometheus range query endpoint.
		handleQueryRange(ctx)

	default:
		ctx.SetStatusCode(http.StatusNotFound)
		fmt.Fprintf(ctx, "Not found: %s", path)
	}
}

// handleMetrics serves internal VictoriaMetrics metrics in Prometheus format.
func handleMetrics(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetStatusCode(http.StatusOK)
	// TODO: expose actual internal metrics.
	fmt.Fprintf(ctx, "# VictoriaMetrics internal metrics\n")
}

// handleRemoteWrite handles Prometheus remote write requests.
func handleRemoteWrite(ctx *fasthttp.RequestCtx) {
	if !ctx.IsPost() {
		ctx.SetStatusCode(http.StatusMethodNotAllowed)
		return
	}
	// TODO: implement Prometheus remote write protocol parsing and storage.
	ctx.SetStatusCode(http.StatusNoContent)
}

// handleQuery handles instant PromQL queries.
func handleQuery(ctx *fasthttp.RequestCtx) {
	// TODO: implement PromQL query execution.
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	fmt.Fprintf(ctx, `{"status":"success","data":{"resultType":"vector","result":[]}}`)
}

// handleQueryRange handles range PromQL queries.
func handleQueryRange(ctx *fasthttp.RequestCtx) {
	// TODO: implement PromQL range query execution.
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(http.StatusOK)
	fmt.Fprintf(ctx, `{"status":"success","data":{"resultType":"matrix","result":[]}}`)
}
