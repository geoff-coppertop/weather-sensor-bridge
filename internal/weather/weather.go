package weather

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

type WeatherData struct {
	topic string
	data  []byte
}

func Start(ctx context.Context, wg *sync.WaitGroup, in <-chan map[string]interface{}) <-chan WeatherData {
	out := make(chan WeatherData)

	wg.Add(1)

	go func() {
		for {
			select {
			case data, ok := <-in:
				if !ok {
					continue
				}

				wxData, err := handleData(data)
				if err != nil {
					continue
				}

				out <- wxData

			case <-ctx.Done():
				close(out)
				wg.Done()
				return
			}
		}
	}()

	return out
}

func handleData(data map[string]interface{}) (WeatherData, error) {
	log.Debug(data)

	return WeatherData{}, nil
}
