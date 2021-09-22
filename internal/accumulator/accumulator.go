package accumulator

import (
	"container/list"
	"fmt"
	"math"
	"time"
)

//go:generate mockgen -destination=../mocks/mock_clock.go -package=mocks github.com/geoff-coppertop/weather-sensor-bridge/internal/accumulator Clock
type Clock interface {
	Now() time.Time
}

type timestampedValue struct {
	value     float64
	timestamp time.Time
}

type Accumulator struct {
	values *list.List
	period time.Duration
	clock  Clock
}

type Stats struct {
	Minimum     float64
	Maximum     float64
	PeriodDelta float64
	Average     float64
}

func New(period time.Duration, clock Clock) *Accumulator {
	acc := Accumulator{
		values: list.New(),
		period: period,
		clock:  clock,
	}

	return &acc
}

func (acc *Accumulator) Accumulate(val float64) (Stats, error) {
	newVal := timestampedValue{
		value:     val,
		timestamp: acc.clock.Now(),
	}

	acc.values.PushBack(newVal)

	/* Pop elements off of the front of the list until the list only goes back
	 * period time from the newest value */
	for {
		val, err := getValue(acc.values.Front())
		if err != nil {
			break
		}

		if val.timestamp.Before(newVal.timestamp.Add(-acc.period)) {
			acc.values.Remove(acc.values.Front())
		} else {
			break
		}
	}

	stat := Stats{
		Minimum:     math.MaxFloat64,
		Maximum:     -math.MaxFloat64,
		PeriodDelta: 0,
		Average:     0,
	}

	/* Iterate through the list to calculate, min, max, and average */
	for e := acc.values.Front(); e != nil; e = e.Next() {
		val, err := getValue(e)
		if err != nil {
			return Stats{}, err
		}

		if val.value > stat.Maximum {
			stat.Maximum = val.value
		}

		if val.value < stat.Minimum {
			stat.Minimum = val.value
		}

		stat.Average += val.value
	}

	stat.Average /= float64(acc.values.Len())

	/* Calculate the start -> end delta of the list by looking at the first and
	 * last elements that are left */
	start, err := getValue(acc.values.Front())
	if err != nil {
		return Stats{}, err
	}

	end, err := getValue(acc.values.Back())
	if err != nil {
		return Stats{}, err
	}

	stat.PeriodDelta = end.value - start.value

	return stat, nil
}

func getValue(e *list.Element) (timestampedValue, error) {
	val, ok := e.Value.(timestampedValue)
	if !ok {
		return timestampedValue{}, fmt.Errorf("unable to convert value %v", e.Value)
	}

	return val, nil
}
