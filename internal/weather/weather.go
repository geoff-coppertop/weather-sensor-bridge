package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	mh "github.com/geoff-coppertop/weather-sensor-bridge/internal/maphelper"
	"github.com/geoff-coppertop/weather-sensor-bridge/internal/math"
	"github.com/geoff-coppertop/weather-sensor-bridge/internal/mqtt"
	"github.com/martinlindhe/unit"
	log "github.com/sirupsen/logrus"
)

const (
	BaseTopic = "sensor/rtl_433"

	TemperatureError        = 0x0FFF
	TemperatureInvalid      = 0x07FA
	TemperatureBelowMinimum = 0x07FC
	TemperatureAboveMaximum = 0x07FD

	HumidityError   = 0xFF
	HumidityInvalid = 0x7A

	SunlightInvalid = 0x1FFFA
	UVIndexInvalid  = 0xFA
)

func Start(ctx context.Context, wg *sync.WaitGroup, in <-chan map[string]interface{}) <-chan mqtt.Data {
	out := make(chan mqtt.Data)

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

func handleData(data map[string]interface{}) (mqtt.Data, error) {
	log.Debug(data)

	topic, err := buildTopicString(data)
	if err != nil {
		return mqtt.Data{}, err
	}

	normalizedData, err := normalizeData(data)
	if err != nil {
		return mqtt.Data{}, err
	}

	txData, err := json.Marshal(normalizedData)
	if err != nil {
		return mqtt.Data{}, err

	}

	return mqtt.Data{
		Topic: topic,
		Data:  txData,
	}, nil
}

func buildTopicString(data map[string]interface{}) (string, error) {
	var err error
	topic := BaseTopic

	if val, ok := mh.GetStringValue(data, "model"); ok {
		topic = mqtt.JoinTopic(topic, val)
	}

	if val, ok := mh.GetStringValue(data, "channel"); ok {
		topic = mqtt.JoinTopic(topic, val)
	}

	if val, ok := mh.GetStringValue(data, "id"); ok {
		topic = mqtt.JoinTopic(topic, val)
	}

	if topic == BaseTopic {
		err = fmt.Errorf("data has no topic information")
	}

	return topic, err
}

func normalizeData(data map[string]interface{}) (map[string]interface{}, error) {
	// https://www.switchdoc.com/wp-content/uploads/2021/04/WeatherRack2Installation1.3.pdf - page 20
	normalizedData := make(map[string]interface{})

	// Battery
	if val, ok := mh.GetBoolValue(data, "batterylow"); ok {
		normalizedData["batterylow"] = val
	}

	// Wind
	if val, ok := mh.GetFloatValue(data, "avewindspeed"); ok {
		// 0+, needs to be in m/s
		normalizedData["windspeedaverage"] = math.Round(val/10, 2)
	}
	if val, ok := mh.GetFloatValue(data, "gustwindspeed"); ok {
		// 0+, needs to be in m/s
		normalizedData["windspeedgust"] = math.Round(val/10, 2)
	}
	if val, ok := mh.GetIntValue(data, "winddirection"); ok {
		// 0 - 359, needs to be in degrees
		val = val % 360
		normalizedData["winddirection"] = val
	}

	// Rain
	if val, ok := mh.GetFloatValue(data, "cumulativerain"); ok {
		// 0+, needs to be in mm
		normalizedData["rainaccumulation"] = math.Round(val/10, 2)
	}

	// Temperature
	if val, ok := mh.GetIntValue(data, "temperature"); ok {
		// Needs to be in C, because we aren't heathens
		switch val {
		case TemperatureError:
		case TemperatureInvalid:
		case TemperatureAboveMaximum:
		case TemperatureBelowMinimum:
			break

		default:
			normalizedData["temperature"] = math.Round(unit.FromFahrenheit(float64(val-400)/10).Celsius(), 2)
		}
	}
	if val, ok := mh.GetIntValue(data, "humidity"); ok {
		// 0 - 100%
		switch val {
		case HumidityError:
		case HumidityInvalid:
			break

		default:
			normalizedData["humidity"] = val
			break
		}
	}

	// Sun
	if val, ok := mh.GetIntValue(data, "light"); ok {
		// 0 - 200k lux
		if (val >= 0) && (val < SunlightInvalid) {
			normalizedData["sunlight"] = val
		}
	}
	if val, ok := mh.GetFloatValue(data, "uv"); ok {
		// 0+?, its a unitless quantity
		if (val >= 0) && (val < UVIndexInvalid) {
			normalizedData["uvindex"] = math.Round(val/10, 2)
		}
	}

	if len(normalizedData) == 0 {
		return normalizedData, fmt.Errorf("no data to normalize from input: %v", data)
	}

	return normalizedData, nil
}
