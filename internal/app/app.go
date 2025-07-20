package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/azoma13/marketplace-service/config"
	v1 "github.com/azoma13/marketplace-service/internal/controller/http/v1"
	"github.com/azoma13/marketplace-service/internal/repo"
	"github.com/azoma13/marketplace-service/internal/service"
	"github.com/azoma13/marketplace-service/pkg/hasher"
	"github.com/azoma13/marketplace-service/pkg/httpserver"
	"github.com/azoma13/marketplace-service/pkg/postgres"
	"github.com/azoma13/marketplace-service/pkg/validator"
	"github.com/labstack/echo/v4"
)

func Run() {
	err := config.NewConfig()
	if err != nil {
		log.Fatal("app - Run - config.NewConfig: %w", err)
		return
	}

	log.Println("Initializing postgres...")
	pg, err := postgres.New(config.Cfg.PG.URL, postgres.MaxPoolSize(config.Cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatal("app - Run - postgres.New: %w", err)
		return
	}
	defer pg.Close()

	log.Println("Initializing repositories...")
	repositories := repo.NewRepositories(pg)

	log.Println("Initializing service...")
	deps := service.ServicesDependencies{
		Repos:    repositories,
		Hasher:   hasher.NewSHA512Hasher(config.Cfg.Hasher.Salt),
		SignKey:  config.Cfg.JWT.SignKey,
		TokenTTL: config.Cfg.JWT.TokenTTL,
	}
	service := service.NewService(deps)

	log.Println("Initializing handlers and routes...")
	handler := echo.New()

	handler.Validator = validator.NewCustomValidator()
	v1.NewRouter(handler, service)

	log.Println("Start http server...")
	log.Printf("Server port: %s", config.Cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(config.Cfg.HTTP.Port))

	log.Println("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	log.Println("Shutting down...")
	err = os.RemoveAll("./content/")
	if err != nil {
		log.Println(fmt.Errorf("app - Run - os.RemoveAll: %w", err))
	}
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
