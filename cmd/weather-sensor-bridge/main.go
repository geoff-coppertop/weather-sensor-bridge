package main

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	cfg "github.com/geoff-coppertop/weather-sensor-bridge/internal/config"
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

	cmd := exec.CommandContext(ctx, "./cmd/weather-sensor-bridge/test.sh")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stdoutDone := make(chan interface{})

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		r := bufio.NewReader(stdout)
		s, e := Readln(r)
		for e == nil {
			log.Println(s)
			s, e = Readln(r)
		}

		log.Info(e)
		log.Info("done")
		close(stdoutDone)
	}()

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

	<-stdoutDone

	log.Info("done for real")
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}
