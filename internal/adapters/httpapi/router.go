package httpapi

import (
	"crypto-aggregator-service/config"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.elastic.co/apm/module/apmchiv5/v2"
	"go.uber.org/zap"
)

// HTTPServer The HTTP server
type HTTPServer struct {
	Logger *zap.SugaredLogger
	sc     config.ServerConfigurations
	Router *chi.Mux
}

// NewHTTPServer Initializes a new httpapi server
func NewHTTPServer(logger *zap.SugaredLogger, serverConf config.ServerConfigurations) *HTTPServer {
	router := chi.NewRouter()

	// APM middleware
	router.Use(apmchiv5.Middleware())

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.AllowContentType("application/json"))

	// Set a timeout value on the request models (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	return &HTTPServer{
		Logger: logger,
		sc:     serverConf,
		Router: router,
	}
}

// Start Fires the httpapi server
func (r *HTTPServer) Start() {
	listeningAddr := ":" + strconv.Itoa(r.sc.Port)
	r.Logger.Infof("Server listening on port %s", listeningAddr)

	// Create a new HTTP server with timeouts
	server := &http.Server{
		Addr:         listeningAddr,
		Handler:      r.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	//err := httpapi.ListenAndServe(listeningAddr, r.Router)
	err := server.ListenAndServe()
	if err != nil {
		r.Logger.Fatalf("Failed to start httpapi server. %v", err)
	}
}
