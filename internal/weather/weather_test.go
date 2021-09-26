package weather

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type TestData struct {
	Input  map[string]interface{} `json:"input"`
	Topic  string                 `json:"topic"`
	Output map[string]interface{} `json:"output"`
}

func getTestData(path string) (TestData, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return TestData{}, err
	}

	var data TestData

	if err := json.Unmarshal(file, &data); err != nil {
		return TestData{}, err
	}

	return data, nil
}

func TestBuildTopicStringEmptyMap(t *testing.T) {
	if _, err := buildTopicString(make(map[string]interface{})); err == nil {
		t.Errorf("expected error")
	}
}

func TestBuildTopicStringTestInput(t *testing.T) {
	test, err := getTestData("test.json")
	if err != nil {
		t.Error("failed to load test data")
	}

	if _, err := buildTopicString(test.Input); err != nil {
		t.Errorf("unexpected error, err: %s", err)
	}
}

func TestNormalizeDataEmptyMap(t *testing.T) {
	if _, err := normalizeData(make(map[string]interface{})); err == nil {
		t.Errorf("expected error")
	}
}

func TestNormalizeDataTestInput(t *testing.T) {
	test, err := getTestData("test.json")
	if err != nil {
		t.Errorf("failed to load test data")
	}

	data, err := normalizeData(test.Input)
	if err != nil {
		t.Errorf("failed to normalize test data")
	}
	output, err := json.Marshal(data)
	if err != nil {
		t.Error("unexpected error")
	}

	input, err := json.Marshal(test.Output)
	if err != nil {
		t.Error("unexpected error")
	}

	if string(input) != string(output) {
		t.Errorf("unexpected error, output: %v, expected: %v", string(output), test.Output)
	}
}
