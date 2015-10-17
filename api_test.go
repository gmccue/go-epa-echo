package epaecho

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestBuildEndpoint(t *testing.T) {
	testURI := "apitest"

	testQueryParams := []string{
		"a=one",
		"b=two",
		"c=three",
	}

	api := &EchoAPIConfig{
		endpointURI: testURI,
		params:      make(url.Values),
	}

	for _, val := range testQueryParams {
		keyVals := strings.Split(val, "=")
		api.params.Add(keyVals[0], keyVals[1])
	}

	api.buildEndpoint()

	expected := fmt.Sprintf("http://ofmpub.epa.gov/echo/%s?%s", testURI, strings.Join(testQueryParams, "&"))

	if api.endpoint.String() != expected {
		t.Error("Invalid URL built: ", api.endpoint.String())
	}
}

func TestUnmarshalEchoBool(t *testing.T) {
	echoBoolJSON := `
	{
        "TestBool": "Yes"
	}`

	echoBoolInstance := struct {
		TestBool echoBool
	}{}

	err := json.Unmarshal([]byte(echoBoolJSON), &echoBoolInstance)
	if err != nil {
		t.Error(err)
	}

	if echoBoolInstance.TestBool != true {
		t.Error("Invalid EchoBool unmarshal value.")
	}
}

func TestUnmarshalEchoDate(t *testing.T) {
	echoDateJSON := `
	{
		"TestDate": "01/30/2016"
	}`

	echoDateInstance := struct {
		TestDate echoDate
	}{}

	err := json.Unmarshal([]byte(echoDateJSON), &echoDateInstance)
	if err != nil {
		t.Error(err)
	}

	date := time.Time(echoDateInstance.TestDate)

	if date.Day() != 30 {
		t.Error("Invalid Time.Day() value returned.")
	}

	if date.Month() != 1 {
		t.Error("Invalid Time.Month() value returned.")
	}

	if date.Year() != 2016 {
		t.Error("Invalid Time.Year() value returned.")
	}
}

func TestUnmarshalEchoFloat(t *testing.T) {
	echoFloatJSON := `
	{
		"TestFloat": "-7.95723"
	}`

	echoFloatInstance := struct {
		TestFloat echoFloat
	}{}

	err := json.Unmarshal([]byte(echoFloatJSON), &echoFloatInstance)
	if err != nil {
		t.Error(err)
	}

	if echoFloatInstance.TestFloat != echoFloat(float64(-7.95723)) {
		t.Error("Invalid echoFloat conversion.")
	}
}

func TestUnmarshalEchoInt(t *testing.T) {
	echoIntJSON := `
	{
		"TestInt": "100"
	}`

	echoIntInstance := struct {
		TestInt echoInt
	}{}

	err := json.Unmarshal([]byte(echoIntJSON), &echoIntInstance)
	if err != nil {
		t.Error(err)
	}

	if echoIntInstance.TestInt != echoInt(100) {
		t.Error("Invalid echoInt conversion")
	}
}
