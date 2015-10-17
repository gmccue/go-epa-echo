package epaecho

import (
	"io/ioutil"
	"testing"
)

func TestSetInvalidQueryParam(t *testing.T) {
	api := NewFacilitiesAPI()

	err := api.SetParam("invalid", "val")
	if err == nil {
		t.Error("No error returned for invalid query parameter.")
	}
}

func TestUnmarshalFacilities(t *testing.T) {
	dataMock, err := ioutil.ReadFile("./echoTest/facilities_test_data.json")
	if err != nil {
		t.Error(err)
	}

	api := NewFacilitiesAPI()
	api.Config.rawResponse = dataMock

	unmarshalErr := api.unmarshalResponse()
	if unmarshalErr != nil {
		t.Error(unmarshalErr)
	}

	if api.response.Results.Message != "Success" {
		t.Error("Invalid success message.")
	}

	if len(api.response.Results.Facilities) != 1 {
		t.Error("Invalid number of facilities.")
	}

	if len(api.response.Results.MapOutput.MapData) != 1 {
		t.Error("Invalid number of MapData results.")
	}
}
