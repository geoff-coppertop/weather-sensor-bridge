package sensor

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"
	"sync"

	log "github.com/sirupsen/logrus"
)

func Start(ctx context.Context, wg *sync.WaitGroup) <-chan map[string]interface{} {
	out := make(chan map[string]interface{})

	wg.Add(1)

	cmd := exec.CommandContext(ctx, "./cmd/weather-sensor-bridge/test.sh", "10")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	go func() {
		r := bufio.NewReader(stdout)

		for {
			line, err := Readln(r)
			if err != nil {
				log.Info(err)
				break
			}

			var sensorData map[string]interface{}
			if err := json.Unmarshal([]byte(line), &sensorData); err != nil {
				log.Error(err)
				continue
			}

			out <- sensorData
		}

		close(out)
		wg.Done()
	}()

	return out
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
