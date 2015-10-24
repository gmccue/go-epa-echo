package epaecho

import (
	"io/ioutil"
	"testing"
)

func TestUnmarshalMaps(t *testing.T) {
	dataMock, err := ioutil.ReadFile("./echotest/map_test_data.json")
	if err != nil {
		t.Error(err)
	}

	api := NewMapAPI("1")
	api.Config.rawResponse = dataMock

	unmarshalErr := api.unmarshalResponse()
	if unmarshalErr != nil {
		t.Error(unmarshalErr)
	}

	if api.response.MapOutput.QueryID != "537" {
		t.Error("Invalid QueryID.")
	}

	if len(api.response.MapOutput.MapData) != 1 {
		t.Error("Invalid number of map results.")
	}

	if api.response.MapOutput.MapData[0].Latitude != "3.188739" {
		t.Error("Invalid map latitude returned")
	}

	if api.response.MapOutput.MapData[0].PUV != "000000" {
		t.Error("Invalid map PUV returned")
	}
}
