package config

import (
	"fmt"
	"net/url"
	"os"

	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Configuration will be pulled from the environment using the following keys
const (
	envDebug = "DEBUG" // Controls debug level for the whole program, should be one of panic, fatal, error, warn, info, debug, or trace

	envServerURL         = "SERVER_URL" // MQTT server URL
	envKeepAlive         = "KA_TIME"    // seconds between keepalive packets
	envConnectRetryDelay = "CRD_TIME"   // milliseconds to delay between connection attempts
)

// Config holds the configuration
type Config struct {
	Debug log.Level

	// MQTT connection details
	ServerURL         *url.URL      // MQTT server URL
	KeepAlive         uint16        // seconds between keepalive packets
	ConnectRetryDelay time.Duration // Period between connection attempts
}

// GetConfig - Retrieves the configuration from the environment
func GetConfig() (Config, error) {
	var cfg Config
	var err error
	var debugLevel string
	var serverURL string
	var iKA int

	if debugLevel, err = stringFromEnv(envDebug); err != nil {
		return Config{}, err
	}
	if cfg.Debug, err = log.ParseLevel(debugLevel); err != nil {
		return Config{}, err
	}

	if serverURL, err = stringFromEnv(envServerURL); err != nil {
		return Config{}, err
	}
	if cfg.ServerURL, err = url.Parse(serverURL); err != nil {
		return Config{}, fmt.Errorf("environmental variable %s must be a valid URL (%w)", envServerURL, err)
	}

	if iKA, err = intFromEnv(envKeepAlive); err != nil {
		return Config{}, err
	}
	cfg.KeepAlive = uint16(iKA)

	if cfg.ConnectRetryDelay, err = milliSecondsFromEnv(envConnectRetryDelay); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// stringFromEnv - Retrieves a string from the environment and ensures it is not blank (ort non-existent)
func stringFromEnv(key string) (string, error) {
	s := os.Getenv(key)
	if len(s) == 0 {
		return "", fmt.Errorf("environmental variable %s must not be blank", key)
	}
	return s, nil
}

// intFromEnv - Retrieves an integer from the environment (must be present and valid)
func intFromEnv(key string) (int, error) {
	s := os.Getenv(key)
	if len(s) == 0 {
		return 0, fmt.Errorf("environmental variable %s must not be blank", key)
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environmental variable %s must be an integer", key)
	}
	return i, nil
}

// milliSecondsFromEnv - Retrieves milliseconds (as time.Duration) from the environment (must be present and valid)
func milliSecondsFromEnv(key string) (time.Duration, error) {
	var i int
	var err error

	if i, err = intFromEnv(key); err != nil {
		return 0, err
	}
	return time.Duration(i) * time.Millisecond, nil
}
