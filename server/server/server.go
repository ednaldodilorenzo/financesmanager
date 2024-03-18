package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/modules/routes"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	App *fiber.App
}

func (s *Server) Setup() {
	if s.App == nil {
		log.Fatalln("Server Incorrectly setup")
	}

	api := s.App.Group("/api")
	routes.SetRoutes(&api)
	config.Connect()
}

func (s *Server) BasicSetup(prefix string, f func(router fiber.Router)) {
	if s.App == nil {
		log.Fatalln("Server Incorrectly setup")
	}

	api := s.App.Group("/api")
	api.Route(prefix, f)
}

func (s *Server) Start() <-chan os.Signal {
	s.Setup()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.App.Listen(":5000"); err != nil {
			log.Fatal(err)
		}
	}()

	return quit
}

func (s *Server) ShutdownGracefully() {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		// Release resources like Database connections
		cancel()
	}()

	shutdownChan := make(chan error, 1)
	go func() { shutdownChan <- s.App.Shutdown() }()

	select {
	case <-timeout.Done():
		log.Fatal("Server Shutdown Timed out before shutdown.")
	case err := <-shutdownChan:
		if err != nil {
			log.Fatal("Error while shutting down server", err)
		} else {
			log.Println("Server Shutdown Successful")
		}
	}
}
