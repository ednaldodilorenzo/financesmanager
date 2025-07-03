package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ansrivas/fiberprometheus/v2"
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
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/dig"
)

type Server struct {
	App      *fiber.App
	db       *config.Database
	mb       *config.Broker
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

func initTracer() func() {
	exporter, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("fiber-app"),
		)),
	)
	otel.SetTracerProvider(tp)
	return func() { _ = tp.Shutdown(context.Background()) }
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
	mb *config.Broker) *Server {
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
	shutdown := initTracer()
	defer shutdown()

	tracer := otel.Tracer("fiber-tracer")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if s.App == nil {
		log.Fatal().Msg("Server Incorrectly setup")
	}

	prometheus := fiberprometheus.New("financecockpit")
	prometheus.RegisterAt(s.App, "/metrics")
	prometheus.SetSkipPaths([]string{"/ping"})            // Optional: Remove some paths from metrics
	prometheus.SetIgnoreStatusCodes([]int{401, 403, 404}) // Optional: Skip metrics for these status codes
	s.App.Use(func(c *fiber.Ctx) error {
		_, span := tracer.Start(c.Context(), c.Path())
		defer span.End()
		return c.Next()
	})
	s.App.Use(prometheus.Middleware)

	s.App.Use(recover.New())
	s.App.Use(middleware.LogRequests)
	api := s.App.Group("/api")
	routes.SetRoutes(&api)
	s.settings.LoadSettings()
	s.db.Connect(&s.settings.Database)
	s.mb.Connect(&s.settings.MessageBroker)
	//s.mb.Config(&s.settings.MessageBroker)
}

func (s *Server) BasicSetup(prefix string, f func(router fiber.Router)) {
	if s.App == nil {
		log.Fatal().Msg("Server Incorrectly setup")
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
			log.Fatal().AnErr("Error", err)
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
		log.Fatal().Msg("Server Shutdown Timed out before shutdown.")
	case err := <-shutdownChan:
		if err != nil {
			log.Debug().AnErr("Error while shutting down server", err)
		} else {
			log.Debug().Msg("Server Shutdown Successful")
		}
	}

	// Close message broker client (asynq.Client)
	if err := s.mb.Close(); err != nil {
		log.Error().AnErr("Failed to close broker client: %v", err)
	} else {
		log.Debug().Msg("Broker client closed successfully")
	}
}
