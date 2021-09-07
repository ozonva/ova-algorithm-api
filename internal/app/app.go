package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/etherlabsio/healthcheck/v2"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "github.com/ozonva/ova-algorithm-api/internal/api"
	"github.com/ozonva/ova-algorithm-api/internal/config"
	"github.com/ozonva/ova-algorithm-api/internal/repo"
	desc "github.com/ozonva/ova-algorithm-api/pkg/ova-algorithm-api"
)

// App is the interface application. Application can be started or restarted
// using ApplyCfg function. Stop stops an application. Use Errors() to get a
// channel reporting errors
type App interface {
	// ApplyCfg dynamically reconfigures configuration and restarts modules
	// Which are affected by configuration change
	ApplyCfg(cfg *config.Config) error
	// Stop synchronously stops an application
	Stop() error
	// Errors returns a channel for getting errors. Please be aware that channel
	// cannot be caches as it might change during restart of the application
	Errors() <-chan error
}

// NewApp creates new empty App  which does nothing unless ApplyCfg is called
func NewApp() App {
	return &app{}
}

type common struct {
	wg    sync.WaitGroup
	db    *sql.DB
	errCh chan error
}

type app struct {
	common     *common
	grpcApp    grpcApp
	monitoring monitoringService
}

func (a *app) ApplyCfg(cfg *config.Config) error {
	if err := a.Stop(); err != nil {
		return fmt.Errorf("error occured while stopping active cfg: %w", err)
	}

	deleter, err := a.initCommon(&cfg.Dsn)
	if err != nil {
		return fmt.Errorf("cannot init database: %w", err)
	}
	defer deleter()

	if err := a.grpcApp.Init(a.common, &cfg.OvaAlgorithm); err != nil {
		return fmt.Errorf("cannot start database: %w", err)
	}
	a.monitoring.Init(a.common, &cfg.Prometheus)

	return nil
}


func (a *app) Stop() error {
	var asyncErrors []error
	asyncErrorsDone := make(chan struct{})

	if a.common != nil {

		go func(c *common) {
			defer close(asyncErrorsDone)
			for err := range c.errCh {
				log.Error().Err(err).Msg("error detected while stopping service")
				asyncErrors = append(asyncErrors, err)
			}
		}(a.common)
	}

	var syncErrors []error

	if err := a.monitoring.Stop(); err != nil {
		syncErrors = append(syncErrors, err)
		log.Error().Err(err).Msg("while stopping monitoring")
	}
	if err := a.grpcApp.Stop(); err != nil {
		syncErrors = append(syncErrors, err)
		log.Error().Err(err).Msg("while stopping grpc")
	}

	var reportBuilder strings.Builder
	if len(syncErrors) != 0 {
		reportBuilder.WriteString("sync errors caught while stopping:\n")
		for i := 0; i < len(syncErrors); i++ {
			fmt.Fprintf(&reportBuilder, "%s;", syncErrors[i].Error())
		}
		reportBuilder.WriteString("\n")
	}

	if a.common != nil {
		<-asyncErrorsDone

		if len(asyncErrors) != 0 {
			reportBuilder.WriteString("async errors caught while stopping:\n")
			for i := 0; i < len(asyncErrors); i++ {
				fmt.Fprintf(&reportBuilder, "%s;", asyncErrors[i].Error())
			}
			reportBuilder.WriteString("\n")
		}
	}

	report := reportBuilder.String()
	if len(report) > 0 {
		return errors.New(report)
	}

	return nil
}

func (a *app) Errors() <-chan error {
	if a.common == nil {
		return nil
	}
	return a.common.errCh
}

func (a *app) initCommon(cfg *config.DSN) (func(), error) {
	a.common = &common{
		errCh: make(chan error),
	}

	var err error
	a.common.db, err = sql.Open("pgx", cfg.MakeStr())
	if err != nil {
		return nil, fmt.Errorf("cannot open database connection: %w", err)
	}

	return func() {
		go func(c *common) {
			c.wg.Wait() // wait all dependent service stop
			defer close(c.errCh)

			if err := c.db.Close(); err != nil {
				log.Error().Err(err).Msg("error occurred while closing db connection")
				c.errCh <- err
			}
		}(a.common)
	}, nil
}

type grpcApp struct {
	common     *common
	grpcServer *grpc.Server
	producer   sarama.AsyncProducer
}

func (a *grpcApp) Stop() error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}
	if a.producer != nil {
		if err := a.producer.Close(); err != nil {
			log.Error().Err(err).Msg("error occurred while closing kafka producer")
			return fmt.Errorf("error occured while closing kafla procuder, %w", err)
		}
	}
	return nil
}

func (a *grpcApp) Init(c *common, cfg *config.OvaAlgorithm) error {
	a.common = c

	var err error

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

	r := repo.NewRepo(a.common.db)
	desc.RegisterOvaAlgorithmApiServer(a.grpcServer, api.NewOvaAlgorithmAPI(r, p))

	c.wg.Add(1)
	go func(c *common) {
		defer c.wg.Done()

		if err := a.grpcServer.Serve(listener); err != nil {
			log.Error().Err(err).Msg("grpcServer stopped with error")
			c.errCh <- fmt.Errorf("grpc server stopped with error: %w", err)
		}
	}(c)

	return nil
}

type monitoringService struct {
	common *common
	server *http.Server
}

func (p *monitoringService) Stop() error {
	var err error
	if p.server != nil {
		if err = p.server.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("error occurred while stopping prometheus server")
		}
	}
	return err
}

func createHealthCheckHandler(db *sql.DB) http.Handler {
	return healthcheck.Handler(

		healthcheck.WithTimeout(5*time.Second),

		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
			),
		),
	)
}

func (p *monitoringService) Init(c *common, cfg *config.Prometheus) {
	p.common = c

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/healthcheck", createHealthCheckHandler(c.db))

	withPort := fmt.Sprintf(":%v", cfg.Port)
	p.server = &http.Server{
		Addr:    withPort,
		Handler: mux,
	}

	c.wg.Add(1)
	go func(c *common) {
		defer c.wg.Done()

		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("prometheus")
			c.errCh <- fmt.Errorf("monitoring service stopped with error: %w", err)
		}
	}(c)
}

func newNotificationProducer(brokerList []string) (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(brokerList, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	go func() {
		for err := range producer.Errors() {
			log.Warn().Err(err).Msg("failed to deliver to notification")
		}
	}()

	return producer, nil
}
