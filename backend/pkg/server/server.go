package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xanuthatusu/tepia/internal/api"
	"github.com/xanuthatusu/tepia/internal/db"
	"github.com/xanuthatusu/tepia/internal/middleware"
	"github.com/xanuthatusu/tepia/internal/sessions"
)

type ServerConfig struct {
	Port string
	Addr string
}

type Server struct {
	Config ServerConfig
}

func ConfigFromEnv() ServerConfig {
	defaults := ServerConfig{
		Port: "8080",
		Addr: "localhost",
	}
	conf := ServerConfig{
		Port: os.Getenv("PORT"),
		Addr: os.Getenv("ADDR"),
	}

	if conf.Port == "" {
		conf.Port = defaults.Port
	}
	if conf.Addr == "" {
		conf.Addr = defaults.Addr
	}

	return conf
}

func New() *Server {
	return &Server{
		Config: ConfigFromEnv(),
	}
}

func (s *Server) Start() error {
	sessions.Init()

	pool := db.Connect()
	defer pool.Close()

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())

	api.RegisterRoutes(r, pool)
	api.RegisterAuthRoutes(r, pool)

	srv := &http.Server{
		Addr:    s.Config.Addr + ":" + s.Config.Port,
		Handler: r,
	}

	go func() {
		fmt.Println("Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	}

	fmt.Println("Server exited gracefully")

	return nil
}
