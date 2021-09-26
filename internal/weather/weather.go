package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	acc "github.com/geoff-coppertop/weather-sensor-bridge/internal/accumulator"
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

type dataSynth func(data acc.Stats) float64

type synthesizer struct {
	outKey   string
	acc      *acc.Accumulator
	dataFunc dataSynth
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

func Start(ctx context.Context, wg *sync.WaitGroup, in <-chan map[string]interface{}) <-chan mqtt.Data {
	out := make(chan mqtt.Data)

	wg.Add(1)

	synthMap := map[string][]synthesizer{
		"wspd_2m":   {synthesizer{"wspd", acc.New(2*time.Minute, realClock{}, acc.ROLLING), getAverage}},
		"rain_1hr":  {synthesizer{"rain_acc", acc.New(1*time.Hour, realClock{}, acc.ROLLING), getPeriodDelta}},
		"rain_24hr": {synthesizer{"rain_acc", acc.New(24*time.Hour, realClock{}, acc.CONSECUTIVE), getPeriodDelta}},
		"wdir_2m":   {synthesizer{"wdir", acc.New(2*time.Minute, realClock{}, acc.ROLLING), getAverage}},
	}

	go func() {
		for {
			select {
			case data, ok := <-in:
				if !ok {
					continue
				}

				wxData, err := handleData(synthMap, data)
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

func handleData(synthMap map[string][]synthesizer, data map[string]interface{}) (mqtt.Data, error) {
	log.Debug(data)

	topic, err := buildTopicString(data)
	if err != nil {
		return mqtt.Data{}, err
	}

	normalizedData, err := normalizeData(data)
	if err != nil {
		return mqtt.Data{}, err
	}

	synthesizedData, err := synthesizeData(synthMap, normalizedData)
	if err != nil {
		return mqtt.Data{}, err
	}

	txData, err := json.Marshal(synthesizedData)
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
		normalizedData["batt"] = val
	}

	// Wind
	if val, ok := mh.GetFloatValue(data, "avewindspeed"); ok {
		// 0+, needs to be in m/s
		normalizedData["wspd"] = math.Round(val/10, 2)
	}
	if val, ok := mh.GetFloatValue(data, "gustwindspeed"); ok {
		// 0+, needs to be in m/s
		normalizedData["wspd_gust"] = math.Round(val/10, 2)
	}
	if val, ok := mh.GetIntValue(data, "winddirection"); ok {
		// 0 - 359, needs to be in degrees
		val = val % 360
		normalizedData["wdir"] = val
	}

	// Rain
	if val, ok := mh.GetFloatValue(data, "cumulativerain"); ok {
		// 0+, needs to be in mm
		normalizedData["rain_acc"] = math.Round(val/10, 2)
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
			normalizedData["temp"] = math.Round(unit.FromFahrenheit(float64(val-400)/10).Celsius(), 2)

			// dewpoint =  state.OutdoorTemperature - ((100.0 - state.OutdoorHumidity) / 5.0);
		}
	}
	if val, ok := mh.GetIntValue(data, "humidity"); ok {
		// 0 - 100%
		switch val {
		case HumidityError:
		case HumidityInvalid:
			break

		default:
			normalizedData["hum"] = val
			break
		}
	}

	// Sun
	if val, ok := mh.GetIntValue(data, "light"); ok {
		// 0 - 200k lux
		if (val >= 0) && (val < SunlightInvalid) {
			normalizedData["light"] = val
		}
	}
	if val, ok := mh.GetFloatValue(data, "uv"); ok {
		// 0+?, it's a unitless quantity
		if (val >= 0) && (val < UVIndexInvalid) {
			normalizedData["uv"] = math.Round(val/10, 2)
		}
	}

	if len(normalizedData) == 0 {
		return normalizedData, fmt.Errorf("no data to normalize from input: %v", data)
	}

	return normalizedData, nil
}

func synthesizeData(synthMap map[string][]synthesizer, data map[string]interface{}) (map[string]interface{}, error) {
	/* Generate dewpoint since it requires two fields of data */
	hValue, hOk := mh.GetFloatValue(data, "hum")
	tValue, tOk := mh.GetFloatValue(data, "temp")
	if tOk && hOk {
		data["dewpoint"] = tValue - ((100.0 - hValue) / 5.0)
	}

	/* Solar radiation is a function of incident light, it's a little bit black magic
	 * https://help.ambientweather.net/help/why-is-the-lux-to-w-m-2-conversion-factor-126-7 */
	if sValue, sOk := mh.GetFloatValue(data, "light"); sOk {
		data["solar"] = sValue / 126.7
	}

	if wValue, wOk := mh.GetFloatValue(data, "wdir"); wOk {
		data["wdir_gust"] = wValue
	}

	/* Generate statistical data */
	for key, synths := range synthMap {
		key = strings.ToLower(key)

		dataValue, ok := mh.GetFloatValue(data, key)
		if !ok {
			log.Errorf("unknown field %s", key)
			continue
		}

		for _, synth := range synths {
			stats, err := synth.acc.Accumulate(dataValue)
			if err != nil {
				log.Error(err)
				continue
			}

			outValue := synth.dataFunc(stats)

			log.Debugf("%s: %.2f", synth.outKey, outValue)

			data[synth.outKey] = outValue
		}
	}

	return data, nil
}

func getAverage(s acc.Stats) float64 {
	return s.Average
}

func getPeriodDelta(s acc.Stats) float64 {
	return s.PeriodDelta
}
