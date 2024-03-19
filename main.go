package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	envName, ok := os.LookupEnv("ENVIRONMENT")
	if !ok || envName == "" {
		envName = "dev"
	}
	envName = strings.ToLower(envName)

	pathToConfDir := path.Join(".", "configuration")
	conf, err := NewAppConf(pathToConfDir, "base.yml", fmt.Sprintf("%s.yml", envName))
	if err != nil {
		log.Fatalf("failed to retrieve app config: %v", err)
	}
	log.Printf("Config: %+v", conf)

	tp, err := traceProvider(envName)
	if err != nil {
		log.Fatalf("failed to initialize trace provider: %v", err)
	}
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	address := fmt.Sprintf("%s:%d", conf.ServerConf.Address, conf.ServerConf.Port)

	tokenService := &TokenService{}
	computeProvisioningService := &ComputeProvisioningService{}

	r := chi.NewRouter()

	r.Route("/api", func(r1 chi.Router) {
		r1.Post("/token", tokenService.CreateToken)

		// Basic handler func to make sure we can deploy!
		r1.Get("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		r1.Group(func(r2 chi.Router) {
			r2.Use(AuthMiddleWare)
			r2.Post("/compute", computeProvisioningService.Provision)
		})
	})

	log.Println("Printing routes...")
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	server := http.Server{
		Addr:         address,
		Handler:      r,
		ReadTimeout:  time.Duration(conf.ServerConf.TimeoutConf.Read) * time.Second,
		WriteTimeout: time.Duration(conf.ServerConf.TimeoutConf.Write) * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("failed to start up: %v", err)
		}
	}()

	<-shutdown
	if err = server.Shutdown(context.Background()); err != nil {
		log.Fatalf("failed to shut down properly: %v", err)
	}

	log.Println("shutdown complete!")
}
