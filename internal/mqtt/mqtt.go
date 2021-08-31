package mqtt

import (
	"context"
	"sync"

	wx "github.com/geoff-coppertop/weather-sensor-bridge/internal/weather"
)

func Start(ctx context.Context, wg *sync.WaitGroup, in <-chan wx.WeatherData) error {
	wg.Add(1)

	go func() {
		for {
			select {
			case data, ok := <-in:
				if !ok {
					continue
				}

				if err := handleData(data); err != nil {
					continue
				}

			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	return nil
}

func handleData(data wx.WeatherData) error {
	return nil
}
