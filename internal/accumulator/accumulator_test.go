//  https://blog.codecentric.de/en/2017/08/gomock-tutorial
package accumulator

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/geoff-coppertop/weather-sensor-bridge/internal/mocks"
)

type realClock struct{}

func (realClock) Now() time.Time { return time.Unix(0, 0) }

func TestRollingWindow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := realClock{}.Now()
	idx := 0

	testData := []struct {
		input  float64
		delay  time.Duration
		output Stats
	}{
		{1.0, 0 * time.Millisecond, Stats{Maximum: 1.0, Minimum: 1.0, Average: 1.0, PeriodDelta: 0.0}},
		{2.0, 15 * time.Second, Stats{Maximum: 2.0, Minimum: 1.0, Average: 1.5, PeriodDelta: 1.0}},
		{2.0, 15 * time.Second, Stats{Maximum: 2.0, Minimum: 2.0, Average: 2.0, PeriodDelta: 0.0}},
	}

	clk := mocks.NewMockClock(ctrl)
	clk.
		EXPECT().
		Now().
		DoAndReturn(
			func() time.Time {
				now = now.Add(testData[idx].delay)
				idx++
				return now
			},
		).
		Times(len(testData))

	period, _ := time.ParseDuration("16s")
	acc := New(period, clk, ROLLING)

	for _, test := range testData {
		stat, err := acc.Accumulate(test.input)

		if err != nil {
			t.Error("")
		}

		if stat.Minimum != test.output.Minimum {
			t.Error("")
		}
		if stat.Maximum != test.output.Maximum {
			t.Error("")
		}
		if stat.PeriodDelta != test.output.PeriodDelta {
			t.Error("")
		}
		if stat.Average != test.output.Average {
			t.Error("")
		}
	}
}

func TestConsecutiveWindow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := realClock{}.Now()
	idx := 0

	testData := []struct {
		input  float64
		delay  time.Duration
		output Stats
	}{
		{1.0, 0 * time.Millisecond, Stats{Maximum: 1.0, Minimum: 1.0, Average: 1.0, PeriodDelta: 0.0}},
		{2.0, 15 * time.Second, Stats{Maximum: 2.0, Minimum: 1.0, Average: 1.5, PeriodDelta: 1.0}},
		{2.0, 15 * time.Second, Stats{Maximum: 2.0, Minimum: 2.0, Average: 2.0, PeriodDelta: 0.0}},
	}

	clk := mocks.NewMockClock(ctrl)
	clk.
		EXPECT().
		Now().
		DoAndReturn(
			func() time.Time {
				now = now.Add(testData[idx].delay)
				idx++
				return now
			},
		).
		Times(len(testData))

	period, _ := time.ParseDuration("16s")
	acc := New(period, clk, CONSECUTIVE)

	for _, test := range testData {
		stat, err := acc.Accumulate(test.input)

		if err != nil {
			t.Error("")
		}

		if stat.Minimum != test.output.Minimum {
			t.Error("")
		}
		if stat.Maximum != test.output.Maximum {
			t.Error("")
		}
		if stat.PeriodDelta != test.output.PeriodDelta {
			t.Error("")
		}
		if stat.Average != test.output.Average {
			t.Error("")
		}
	}
}
