package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	cfg "github.com/geoff-coppertop/weather-sensor-bridge/internal/config"
	pub "github.com/geoff-coppertop/weather-sensor-bridge/internal/publisher"
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

	pubCh := pub.Start(ctx, &wg, cfg, wxCh)

	WaitProcess(&wg, pubCh, cancel)
}

func WaitProcess(wg *sync.WaitGroup, ch <-chan error, cancel context.CancelFunc) {
	log.Info("Waiting")

	select {
	case <-OSExit():
		log.Info("signal caught - exiting")

	case <-ch:
		log.Errorf("uh-oh")
	}

	cancel()

	log.Info("cancelled")

	wg.Wait()

	log.Info("goodbye")
}

func OSExit() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	return sig
}
