package main

import (
	"context"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/ozonva/ova-algorithm-api/internal/config"
	"github.com/ozonva/ova-algorithm-api/internal/repo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"

	api "github.com/ozonva/ova-algorithm-api/internal/api"
	desc "github.com/ozonva/ova-algorithm-api/pkg/ova-algorithm-api"
	"google.golang.org/grpc/reflection"

	"database/sql"
	_ "github.com/jackc/pgx/stdlib"

	"github.com/Shopify/sarama"
	"github.com/ozonva/ova-algorithm-api/internal/tracer"
)

func newNotificationProducer(brokerList []string) (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, cfg)
	if err != nil {
		return producer, err
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Warn().Err(err).Msg("failed to write to notification")
		}
	}()

	return sarama.NewAsyncProducer(brokerList, cfg)
}

type GrpcApp struct {
	grpcServer *grpc.Server
	producer   sarama.AsyncProducer
	db         *sql.DB
}

func (a *GrpcApp) Stop() {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
		a.grpcServer = nil
	}
	if a.producer != nil {
		if err := a.producer.Close(); err != nil {
			log.Error().Err(err).Msg("error occurred while closing kafka producer")
			// error dropped intentionally
		}
		a.producer = nil
	}
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Error().Err(err).Msg("error occurred while closing db connection")
			// error dropped intentionally
		}
		a.db = nil
	}
}

func (a *GrpcApp) Start(cfg *config.OvaAlgorithm) error {
	var err error
	a.db, err = sql.Open("pgx", cfg.Dsn.MakeStr())
	if err != nil {
		return fmt.Errorf("cannot open database connection: %w", err)
	}

	a.grpcServer = grpc.NewServer()
	reflection.Register(a.grpcServer)

	brokerAddr := fmt.Sprintf("%v:%v", cfg.Broker.Hostname, cfg.Broker.Port)
	p, err := newNotificationProducer([]string{brokerAddr})
	if err != nil {
		return fmt.Errorf("cannot create kafka producer: %w", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		return fmt.Errorf("cannot listen to port: %w", err)
	}

	r := repo.NewRepo(a.db)
	desc.RegisterOvaAlgorithmApiServer(a.grpcServer, api.NewOvaAlgorithmApi(r, p))

	go func() {
		if err := a.grpcServer.Serve(listener); err != nil {
			log.Error().Err(err).Msg("grpcServer stopped with error")
		}
	}()

	return nil
}

func (a *GrpcApp) applyCfg(cfg *config.OvaAlgorithm) error {
	a.Stop()
	if err := a.Start(cfg); err != nil {
		log.Error().Err(err).Msg("failed to start OvaAlgorithm Service")
		a.Stop() //Clean-up
		return fmt.Errorf("failed to start OvaAlgorithm Service: %w", err)
	}
	return nil
}

type PrometheusService struct {
	server *http.Server
}

func (p *PrometheusService) Stop() {
	if p.server != nil {
		if err := p.server.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("error occurred while stopping prometheus server")
			// error dropped intentionally
		}
		p.server = nil
	}
}

func (p *PrometheusService) Start(cfg *config.Prometheus) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	withPort := fmt.Sprintf(":%v", cfg.Port)
	p.server = &http.Server{
		Addr:    withPort,
		Handler: mux,
	}

	go func() {
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("prometheus")
		}
	}()
}

func (p *PrometheusService) applyCfg(cfg *config.Prometheus) {
	p.Stop()
	p.Start(cfg)
}

func main() {
	// Configure and enable tracer
	tracer, closer, err := tracer.NewTracer()
	if err != nil {
		log.Fatal().Msg("cannot start tracer")
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	configMonitor := config.NewMonitorConfig(nil)
	defer configMonitor.Stop()

	app := GrpcApp{}
	prometheus := PrometheusService{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigs:
			log.Info().Msgf("%v signal received. terminating...", sig)
			return
		case cfg := <-configMonitor.Updates():
			log.Debug().Msg("new config received")
			if err := app.applyCfg(&cfg.OvaAlgorithm); err != nil {
				log.Fatal().Err(err).Msg("cannot apply Ova Algorithm config")
			}
			prometheus.applyCfg(&cfg.Prometheus)

		case err := <-configMonitor.Errors():
			log.Fatal().Err(err).Msg("error occurred in config reader")
		}
	}
}
