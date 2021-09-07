package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/ozonva/ova-algorithm-api/internal/app"
	"github.com/ozonva/ova-algorithm-api/internal/config"
	"github.com/ozonva/ova-algorithm-api/internal/tracer"
)

func main() {

	log.Logger = log.Output(&lumberjack.Logger{
		Filename:   "ova-algorithm.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
	})

	// Configure and enable tracer
	tracer, closer, err := tracer.NewTracer()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start tracer")
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	configMonitor := config.NewMonitorConfig(nil)
	defer configMonitor.Stop()

	app := app.NewApp()
	defer func() {
		if err := app.Stop(); err != nil {
			log.Error().Err(err).Msg("got error on deferred app stop")
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigs:
			log.Info().Msgf("%v signal received. terminating...", sig)
			return
		case cfg, ok := <-configMonitor.Updates():
			log.Debug().Bool("ok", ok).Msg("new config received")
			if ok {
				if err := app.ApplyCfg(cfg); err != nil {
					log.Fatal().Err(err).Msg("cannot apply Ova Algorithm config")
				}
			} else {
				log.Fatal().Msg("unexpected config monitor updates closed")
			}

		case err, ok := <-configMonitor.Errors():
			if ok {
				log.Fatal().Err(err).Msg("error occurred in config reader")
			} else {
				log.Fatal().Msg("unexpected config monitor errors channel closed")
			}

		case err, ok := <-app.Errors():
			if ok {
				log.Fatal().Err(err).Msg("error occurred in config reader")
			} else {
				log.Fatal().Msg("unexpected app errors channel closed")
			}
		}
	}
}
