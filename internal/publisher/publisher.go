package publisher

import (
	"context"
	"sync"

	cfg "github.com/geoff-coppertop/weather-sensor-bridge/internal/config"
	"github.com/geoff-coppertop/weather-sensor-bridge/internal/mqtt"
	log "github.com/sirupsen/logrus"
)

func Start(ctx context.Context, wg *sync.WaitGroup, cfg cfg.Config, in <-chan mqtt.Data) <-chan error {
	out := make(chan error)

	var con *mqtt.Connection
	var err error = nil

	if con, err = mqtt.Connect(ctx, cfg); err != nil {
		log.Debug("connect failed")
		log.Debug(err)
		close(out)
		return out
	} else {
		log.Debug("connection started")
		wg.Add(1)
	}

	go func() {
		defer close(out)
		defer wg.Done()

		log.Debug("Waiting for the connection")
		select {
		case <-con.OnConnectionUp():
			log.Debug("<-con.OnConnectionUp")

		case err := <-con.OnError():
			log.Debug("<-con.OnError")
			log.Error(err)
			con.Disconnect()
			return

		case <-ctx.Done():
			log.Debug("<-ctx.Done")
			con.Disconnect()
			return
		}
		log.Debug("connected")

		for {
			select {
			case err := <-con.OnError():
				log.Debug("<-con.OnError")
				log.Error(err)
				con.Disconnect()
				return

			case <-con.OnServerDisconnect():
				log.Debug("disconnect sig")
				return

			case data, ok := <-in:
				if !ok {
					continue
				}

				con.Publish(data)

			case <-ctx.Done():
				log.Debug("<-ctx.Done")
				con.Disconnect()
				return
			}
		}
	}()

	return out
}
