package epaecho

import (
	"io/ioutil"
	"testing"
)

func TestUnmarshalQueryData(t *testing.T) {
	dataMock, err := ioutil.ReadFile("./echotest/qid_test_data.json")
	if err != nil {
		t.Error(err)
	}

	api := NewQueryAPI("1")
	api.Config.rawResponse = dataMock

	unmarshalErr := api.unmarshalResponse()
	if unmarshalErr != nil {
		t.Error(unmarshalErr)
	}

	if api.response.Results.QueryRows != 574 {
		t.Error("Invalid number of QueryRows.")
	}

	if api.response.Results.PageNo != 1 {
		t.Error("Invalid PageNo.")
	}

	fFacility := api.response.Results.Facilities[0]

	if fFacility.RegistryID != "000000" {
		t.Error("Invalid facility registry ID.")
	}
	if fFacility.Name != "ELOHSSA CORPORATION MINE" {
		t.Error("Invalid facility name.")
	}
}
