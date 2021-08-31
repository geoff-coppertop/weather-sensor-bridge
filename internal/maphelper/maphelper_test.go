package maphelper

import (
	"testing"
)

func TestGetBoolValue(t *testing.T) {
	var tests = []struct {
		input map[string]interface{}
		key   string
		val   bool
		ok    bool
	}{
		{map[string]interface{}{"test": true}, "test", true, true},
		{map[string]interface{}{"test": false}, "test", false, true},
		{map[string]interface{}{"test": 1}, "test", true, true},
		{map[string]interface{}{"test": 0}, "test", false, true},
		{map[string]interface{}{"test": -1}, "test", true, true},
		{map[string]interface{}{"test": -1.0}, "test", true, true},
		{map[string]interface{}{"test": "test"}, "test", false, false},
		{map[string]interface{}{"test": true}, "banana", false, false},
		{map[string]interface{}{}, "test", false, false},
	}

	for _, test := range tests {
		val, ok := GetBoolValue(test.input, test.key)

		if ok != test.ok {
			t.Errorf("unexpected error")
		}

		if val != test.val {
			t.Errorf("unexpected error")
		}
	}
}

func TestGetStringValue(t *testing.T) {
	var tests = []struct {
		input map[string]interface{}
		key   string
		val   string
		ok    bool
	}{
		{map[string]interface{}{"test": "test"}, "test", "test", true},
		{map[string]interface{}{"test": 0}, "test", "0", true},
		{map[string]interface{}{"test": -1}, "test", "-1", true},
		{map[string]interface{}{"test": 0.000000}, "test", "0.000000", true},
		{map[string]interface{}{"test": 1.0}, "test", "1.000000", true},
		{map[string]interface{}{"test": true}, "test", "true", true},
		{map[string]interface{}{"test": false}, "test", "false", true},
		{map[string]interface{}{}, "test", "", false},
	}

	for _, test := range tests {
		val, ok := GetStringValue(test.input, test.key)

		if ok != test.ok {
			t.Errorf("unexpected error")
		}

		if val != test.val {
			t.Errorf("unexpected value")
		}
	}
}

func TestGetFloatValue(t *testing.T) {
	var tests = []struct {
		input map[string]interface{}
		key   string
		val   float64
		ok    bool
	}{
		{map[string]interface{}{"test": 0.00}, "test", 0.00, true},
		{map[string]interface{}{"test": 0}, "test", 0.00, true},
		{map[string]interface{}{"test": 0.0}, "test", 0.00, true},
		{map[string]interface{}{"test": "test"}, "test", 0.00, false},
		{map[string]interface{}{}, "test", 0.00, false},
	}

	for _, test := range tests {
		val, ok := GetFloatValue(test.input, test.key)

		if ok != test.ok {
			t.Errorf("unexpected error")
		}

		if val != test.val {
			t.Errorf("unexpected error")
		}
	}
}

func TestGetIntValue(t *testing.T) {
	var tests = []struct {
		input map[string]interface{}
		key   string
		val   int
		ok    bool
	}{
		{map[string]interface{}{"test": 1}, "test", 1, true},
		{map[string]interface{}{"test": 0}, "test", 0, true},
		{map[string]interface{}{"test": -1}, "test", -1, true},
		{map[string]interface{}{"test": -1.0}, "test", -1, true},
		{map[string]interface{}{"test": false}, "test", 0, false},
		{map[string]interface{}{"test": "test"}, "test", 0, false},
		{map[string]interface{}{}, "test", 0, false},
	}

	for _, test := range tests {
		val, ok := GetIntValue(test.input, test.key)

		if ok != test.ok {
			t.Errorf("unexpected error")
		}

		if val != test.val {
			t.Errorf("unexpected error")
		}
	}
}
