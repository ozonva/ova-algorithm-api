package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/ozonva/ova-algorithm-api/internal/repo"
	"github.com/rs/zerolog/log"
	"io"
	"net"
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

	"github.com/ozonva/ova-algorithm-api/internal/tracer"
)

const (
	grpcPort = ":44555"
)

func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen:")
	}

	dsn := "user=melkozer password=melkozer dbname=ova sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("cannot open database connection: %w", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	r := repo.NewRepo(db)
	desc.RegisterOvaAlgorithmApiServer(s, api.NewOvaAlgorithmApi(r))

	if err := s.Serve(listen); err != nil {
		log.Err(err).Msg("failed to listen:")
	}

	return nil
}

type Config struct{}

func readConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file \"%v\"", path)
	}
	defer file.Close()

	data := make([]byte, 256)

	for {
		count, err := file.Read(data)

		if count > 0 {
			fmt.Printf("%s", data[:count])
		}

		if err != nil {
			if err == io.EOF {

				break
			} else {
				return nil, fmt.Errorf("error occured while config \"%v\": %w", path, err)
			}
		}
	}

	// TODO: remove stub
	return &Config{}, nil
}

func monitorConfig() <-chan *Config {
	c := make(chan *Config)

	go func() {
		for {
			newConfig, err := readConfig("configs/config.json")

			if err != nil {
				fmt.Printf("update Config failed: %v\n", err.Error())
			}
			// TODO: add comparison configs as soon as settings added to Config
			// if newConfig != oldConfig {
			c <- newConfig
			// }

			time.Sleep(3600 * time.Second)
		}
	}()

	return c
}

func main() {
	// Configure and enable tracer
	tracer, closer, err := tracer.NewTracer()
	if err != nil {
		log.Fatal().Msg("cannot start tracer")
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	configUpdates := monitorConfig()

	go func() {
		if err := run(); err != nil {
			log.Error().Err(err).Msg("error during service")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigs:
			log.Info().Msgf("%v signal received. terminating...", sig)
			return
		case <-configUpdates:
			log.Log().Msg("new config received")
		}
	}
}
