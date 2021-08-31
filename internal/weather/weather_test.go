//  https://blog.codecentric.de/en/2017/08/gomock-tutorial
package weather

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/geoff-coppertop/weather-sensor-bridge/internal/mocks"
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
	if _, err := BuildTopicString(make(map[string]interface{})); err == nil {
		t.Errorf("expected error")
	}
}

func TestBuildTopicStringTestInput(t *testing.T) {
	test, err := getTestData("test.json")
	if err != nil {
		t.Error("failed to load test data")
	}

	if _, err := BuildTopicString(test.Input); err != nil {
		t.Errorf("unexpected error, err: %s", err)
	}
}

func TestNormalizeDataEmptyMap(t *testing.T) {
	if _, err := NormalizeData(make(map[string]interface{})); err == nil {
		t.Errorf("expected error")
	}
}

func TestNormalizeDataTestInput(t *testing.T) {
	test, err := getTestData("test.json")
	if err != nil {
		t.Errorf("failed to load test data")
	}

	data, err := NormalizeData(test.Input)
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
		t.Errorf("unexpected error, output: %v, expected: %v", output, test.Output)
	}
}

func TestTopicStringFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	topicError := errors.New("topic error")

	wx := mocks.NewMockWeatherFuncs(ctrl)
	wx.
		EXPECT().
		BuildTopicString(gomock.Any()).
		Return("", topicError).
		Times(1)

	mqtt := mocks.NewMockDataFuncs(ctrl)
	mqtt.
		EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		Times(0)

	if err := PublishData(DataFuncs{Wx: wx, MQTT: mqtt}, map[string]interface{}{}); err == nil {
		t.Errorf("expected error")
	}
}

func TestNormalizeDataFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	normalizeDataError := errors.New("topic error")

	wx := mocks.NewMockWeatherFuncs(ctrl)
	wx.
		EXPECT().
		BuildTopicString(gomock.Any()).
		Return("topic", nil).
		Times(1)
	wx.
		EXPECT().
		NormalizeData(gomock.Any()).
		Return(map[string]interface{}{}, normalizeDataError).
		Times(1)

	mqtt := mocks.NewMockDataFuncs(ctrl)
	mqtt.
		EXPECT().
		Publish(gomock.Any(), gomock.Any()).
		Times(0)

	if err := PublishData(DataFuncs{Wx: wx, MQTT: mqtt}, map[string]interface{}{}); err == nil {
		t.Errorf("expected error")
	}
}

func TestPublishFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wx := mocks.NewMockWeatherFuncs(ctrl)
	wx.
		EXPECT().
		BuildTopicString(gomock.Any()).
		Return("topic", nil).
		Times(1)
	wx.
		EXPECT().
		NormalizeData(gomock.Any()).
		Return(map[string]interface{}{"test": "banana"}, nil).
		Times(1)

	publishError := errors.New("publish error")

	mqtt := mocks.NewMockDataFuncs(ctrl)
	mqtt.
		EXPECT().
		Publish("topic", gomock.Any()).
		Return(publishError).
		Times(1)

	if err := PublishData(DataFuncs{Wx: wx, MQTT: mqtt}, map[string]interface{}{}); err == nil {
		t.Errorf("expected error")
	}
}

func TestPublishSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wx := mocks.NewMockWeatherFuncs(ctrl)
	wx.
		EXPECT().
		BuildTopicString(gomock.Any()).
		Return("topic", nil).
		Times(1)
	wx.
		EXPECT().
		NormalizeData(gomock.Any()).
		Return(map[string]interface{}{"test": "banana"}, nil).
		Times(1)

	mqtt := mocks.NewMockDataFuncs(ctrl)
	mqtt.
		EXPECT().
		Publish("topic", gomock.Any()).
		Return(nil).
		Times(1)

	if err := PublishData(DataFuncs{Wx: wx, MQTT: mqtt}, map[string]interface{}{}); err != nil {
		t.Errorf("unexpected error")
	}
}
