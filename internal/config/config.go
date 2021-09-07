package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

const (
	defaultConfigPath = "configs/config.json"
)

type DSN struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	SslMode  string `json:"sslmode"`
}

func (dsn *DSN) MakeStr() string {
	return fmt.Sprintf("user=%v password=%v database=%v sslmode=%v",
		dsn.User, dsn.Password, dsn.Database, dsn.SslMode)
}

type Broker struct {
	Hostname string `json:"hostname"`
	Port     uint16 `json:"port"`
	Topic    string `json:"topic"`
}

type OvaAlgorithm struct {
	Port   uint16 `json:"port"`
	Broker Broker `json:"broker"`
}

type Prometheus struct {
	Port uint16 `json:"port"`
}

type Config struct {
	Dsn          DSN          `json:"dsn"`
	OvaAlgorithm OvaAlgorithm `json:"ovaAlgorithm"`
	Prometheus   Prometheus   `json:"prometheus"`
}

func (c *Config) Validate() error {
	return nil
}

type MonitorConfig interface {
	Stop()
	Updates() <-chan *Config
	Errors() <-chan error
}

func NewMonitorConfig(path *string) MonitorConfig {
	monitor := &monitorConfig{
		configPath: defaultConfigPath,
		updates:    make(chan *Config),
		stop:       make(chan struct{}),
	}

	if path != nil {
		monitor.configPath = *path
	}

	go monitor.watch()

	return monitor
}

type monitorConfig struct {
	configPath   string
	activeConfig *Config
	updates      chan *Config
	errors       chan error
	stop         chan struct{}
}

func (m *monitorConfig) Stop() {
	m.stop <- struct{}{}
	close(m.stop)
}

func (m *monitorConfig) Updates() <-chan *Config {
	return m.updates
}

func (m *monitorConfig) Errors() <-chan error {
	return m.errors
}

func (m *monitorConfig) readConfig() (*Config, error) {
	jsonFile, err := ioutil.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}
	config := new(Config)
	if err = json.Unmarshal(jsonFile, config); err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("indalid config: %w", err)
	}

	return config, nil
}

func (m *monitorConfig) onNewConfig(config *Config) {
	if m.activeConfig == nil || *m.activeConfig != *config {
		m.activeConfig = config
		m.updates <- config
	}
}

func (m *monitorConfig) tryToReadConfigAndNotifyResult() {
	config, err := m.readConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read config")
		m.errors <- err
	} else {
		m.onNewConfig(config)
	}
}

func (m *monitorConfig) watch() {
	defer close(m.updates)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create config watcher")
	}
	defer watcher.Close()

	err = watcher.Add(m.configPath)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create config watcher")
	}

	//read initially without any notification
	m.tryToReadConfigAndNotifyResult()

	for {
		select {
		case <-m.stop:
			log.Info().Msg("config monitor has been stopped")
			return
		case event, ok := <-watcher.Events:
			if !ok {
				log.Error().Msg("watcher error: stopped watching")
				return
			}
			log.Debug().Str("name", event.Name).Str("op", event.Op.String()).Msg("new event")
			if event.Op&fsnotify.Write == fsnotify.Write {
				m.tryToReadConfigAndNotifyResult()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error().Err(err).Msg("error occurred while watching")
		}
	}
}
