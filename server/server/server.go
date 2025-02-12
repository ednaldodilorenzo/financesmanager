package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ednaldo-dilorenzo/iappointment/config"
	"github.com/ednaldo-dilorenzo/iappointment/middleware"
	"github.com/ednaldo-dilorenzo/iappointment/model"
	"github.com/ednaldo-dilorenzo/iappointment/modules/account"
	"github.com/ednaldo-dilorenzo/iappointment/modules/auth"
	"github.com/ednaldo-dilorenzo/iappointment/modules/budget"
	"github.com/ednaldo-dilorenzo/iappointment/modules/category"
	"github.com/ednaldo-dilorenzo/iappointment/modules/generic"
	"github.com/ednaldo-dilorenzo/iappointment/modules/planning"
	"github.com/ednaldo-dilorenzo/iappointment/modules/routes"
	"github.com/ednaldo-dilorenzo/iappointment/modules/tag"
	"github.com/ednaldo-dilorenzo/iappointment/modules/transaction"
	"github.com/ednaldo-dilorenzo/iappointment/util"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/dig"
)

type Server struct {
	App      *fiber.App
	db       *config.Database
	mb       util.EmailSender
	settings *config.Settings
}

type ServerDependencies struct {
	dig.In
	AuthController        auth.AuthController
	AccountController     generic.GenericController[*model.Account]
	TransactionController transaction.TransactionController
	CategoryController    generic.GenericController[*model.Category]
	TagController         tag.TagController
	DB                    config.Database
}

func NewServer(authController auth.AuthController,
	accountController generic.GenericController[*model.Account],
	transactionController transaction.TransactionController,
	categoryController generic.GenericController[*model.Category],
	tagController tag.TagController,
	planningControler planning.PlanningController,
	budgetController budget.BudgetController,
	deserializer *middleware.Deserializer,
	db *config.Database,
	settings *config.Settings,
	mb util.EmailSender) *Server {
	server := &Server{
		App:      InitFiberApplication(),
		db:       db,
		mb:       mb,
		settings: settings,
	}

	api := server.App.Group("/api")
	api.Route(auth.GetRoutes(authController, deserializer))
	api.Route(category.GetRoutes(categoryController, deserializer))
	api.Route(account.GetRoutes(accountController, deserializer))
	api.Route(transaction.GetRoutes(transactionController, deserializer))
	api.Route(planning.GetRoutes(planningControler, deserializer))
	api.Route(tag.GetRoutes(tagController, deserializer))
	api.Route(budget.GetRoutes(budgetController, deserializer))
	return server
}

func (s *Server) Setup() {
	if s.App == nil {
		log.Fatalln("Server Incorrectly setup")
	}

	api := s.App.Group("/api")
	routes.SetRoutes(&api)
	s.db.Connect(&s.settings.Database)
	s.mb.Config(&s.settings.MessageBroker)
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
