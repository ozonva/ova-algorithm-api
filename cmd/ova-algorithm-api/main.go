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

	"github.com/Shopify/sarama"
	"github.com/ozonva/ova-algorithm-api/internal/tracer"
)

const (
	grpcPort = ":44555"
)

func newNotificationProducer(brokerList []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokerList, config)
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

	go func() {
		for suc := range producer.Successes() {
			key, _ := suc.Key.Encode()
			value, _ := suc.Key.Encode()
			log.Debug().Bytes("key", key).Bytes("value", value).Msg("notification")
		}
	}()

	return sarama.NewAsyncProducer(brokerList, config)
}

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

	p, err := newNotificationProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatal().Msg("Failed to create a kafka producer")
	}

	if err != nil {
		log.Fatal().Msg("Failed to create Kafka producer")
	}

	r := repo.NewRepo(db)
	s.GracefulStop()
	desc.RegisterOvaAlgorithmApiServer(s, api.NewOvaAlgorithmApi(r, p))

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
