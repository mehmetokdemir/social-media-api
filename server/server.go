package server

import (
	"fmt"
	"github.com/ansrivas/fiberprometheus/v2"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/mehmetokdemir/social-media-api/docs"
	"github.com/mehmetokdemir/social-media-api/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

type Server struct {
	app    *fiber.App
	config config.Config
	logger *zap.SugaredLogger
}

func New(handlers []Handler, config config.Config, logger *zap.SugaredLogger) *Server {
	app := fiber.New()
	app.Use(cors.New())

	for _, handler := range handlers {
		handler.RegisterRoutes(app)
	}
	server := &Server{
		app:    app,
		config: config,
		logger: logger,
	}
	server.AddRoutes()
	return server
}

func (s *Server) AddRoutes() {
	s.app.Get("/health", s.healthCheck)
	s.app.Get("/swagger/*", swagger.HandlerDefault)
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description Get the status of server.
// @Tags Health
// @Accept */*
// @Success 200
// @Failure 404
// @Router /health [get]
func (s *Server) healthCheck(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

// @title Social Media API
// @version 2.0
// @description This is a social media api
// @BasePath /
// @schemes http
func (s *Server) Start() error {
	promRegistry := prometheus.NewRegistry()
	// Register the metrics.
	requestsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "The total number of requests.",
	})
	cpuUsage := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "The current CPU usage.",
	})
	responseTimes := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "response_times",
		Help: "The response times for requests.",
	})
	requestSizes := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "request_sizes",
		Help: "The sizes of requests.",
	})
	promRegistry.MustRegister(requestsTotal, cpuUsage, responseTimes, requestSizes)

	fiberPrometheus := fiberprometheus.New("social-media-api")
	fiberPrometheus.RegisterAt(s.app, "/metrics")
	s.app.Use(fiberPrometheus.Middleware)

	address := fmt.Sprintf(":%s", s.config.ServerPort)
	shutDownChan := make(chan os.Signal, 1)
	signal.Notify(shutDownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutDownChan
		if err := s.app.Shutdown(); err != nil {
			s.logger.Error(err)
		}
	}()
	return s.app.Listen(address)
}
