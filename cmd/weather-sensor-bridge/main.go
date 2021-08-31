package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	cfg "github.com/geoff-coppertop/weather-sensor-bridge/internal/config"
	"github.com/geoff-coppertop/weather-sensor-bridge/internal/mqtt"
	sns "github.com/geoff-coppertop/weather-sensor-bridge/internal/sensor"
	wx "github.com/geoff-coppertop/weather-sensor-bridge/internal/weather"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := cfg.GetConfig()
	if err != nil {
		log.Panic(err)
	}

	log.SetLevel(cfg.Debug)
	log.Info("Starting")

	log.Debug(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	snsCh := sns.Start(ctx, &wg)

	wxCh := wx.Start(ctx, &wg, snsCh)

	mqtt.Start(ctx, &wg, wxCh)

	// Messages will be handled through the callback so we really just need to wait until a shutdown
	// is requested
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	log.Info("Waiting")

	<-sig

	log.Info("signal caught - exiting")

	cancel()

	log.Info("cancelled")

	wg.Wait()

	log.Info("goodbye")
}
