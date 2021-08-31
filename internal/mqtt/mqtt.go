package mqtt

import (
	"context"
	"sync"
)

type Data struct {
	Topic string
	Data  []byte
}

func Start(ctx context.Context, wg *sync.WaitGroup, in <-chan Data) error {
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

func SanitizeTopic(topic string) string {
	return ""
}

func JoinTopic(topics ...string) string {
	return ""
}

func handleData(data Data) error {
	return nil
}
