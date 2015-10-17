// Package epaecho provides programmatic access to the EPA's Compliance
// History "All Data Facility Search" API (http://echo.epa.gov).
package epaecho

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultHttpTimeout    = 60
	echoAPIBaseScheme     = "http"
	echoAPIBaseHost       = "ofmpub.epa.gov"
	echoAPIBasePath       = "echo"
	echoDefaultDateFormat = "01/02/2006"
	echoNullDateValue     = "01/01/1900"
)

var (
	echoAPIError           = errors.New("The Echo API request returned an error.")
	invalidQueryParamError = errors.New("Invalid query parameter provided.")
	invalidResponseError   = errors.New("The HTTP request failed.")
	missingArgumentError   = errors.New("Required argument missing.")
)

// echoBool is a custom boolean type used for properly unmarshaling JSON.
type echoBool bool

// echoDate is a custom time.Time type used for properly unmarshaling JSON.
type echoDate time.Time

// echoFloat is a custom float type used for properly unmarshaling JSON.
type echoFloat float64

// echoInt is a custom int type used for properly unmarshaling JSON.
type echoInt int

// EchoAPIConfig holds common configuration for building API requests.
type EchoAPIConfig struct {
	Debug       bool
	client      *http.Client
	endpoint    url.URL
	endpointURI string
	params      url.Values
	rawResponse []byte
}

// EchoAPIError reports an error and the associated message.
type EchoAPIError struct {
	ErrorMessage string
}

type echoAPI interface {
	SetParam(key string, value string) error
	unmarshalResponse() error
}

func newEchoAPIConfig() *EchoAPIConfig {
	return &EchoAPIConfig{
		client: &http.Client{
			Timeout: time.Duration(defaultHttpTimeout) * time.Second,
		},
		params: make(url.Values),
	}
}

// SetTimeout sets the HTTP request timeout in seconds
func (api *EchoAPIConfig) SetTimeout(timeout int) {
	api.client.Timeout = time.Duration(timeout) * time.Second
}

func (api *EchoAPIConfig) buildEndpoint() error {
	api.endpoint.Scheme = echoAPIBaseScheme
	api.endpoint.Host = echoAPIBaseHost
	api.endpoint.Path = fmt.Sprintf("%s/%s", echoAPIBasePath, api.endpointURI)
	api.endpoint.RawQuery = api.params.Encode()

	return nil
}

// addParam checks a provided query parameter against a provided human-readable map.
// If a valid parameter is found, a lookup is performed against the provided map,
// and a new query parameter is added to the EchoAPIConfig struct.
func (api EchoAPIConfig) addParam(key string, value string, optionsMap map[string]string) error {
	if apiVal, ok := optionsMap[key]; ok {
		// If this parameter is already set, remove it.
		if len(api.params.Get(apiVal)) > 0 {
			api.params.Del(apiVal)
		}

		api.params.Add(apiVal, value)
	} else {
		return fmt.Errorf("%s Parameter provided was: %s.", invalidQueryParamError, key)
	}

	return nil
}

func (api *EchoAPIConfig) doRequest() error {
	api.buildEndpoint()

	if api.Debug {
		log.Printf("API endpoint URL: %s", api.endpoint.String())
	}

	req, err := http.NewRequest("GET", api.endpoint.String(), nil)
	if err != nil {
		return err
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if api.Debug {
		log.Printf("Response body: %s", string(body))
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s HTTP status code returned was: %s.", invalidResponseError, resp.StatusCode)
	}

	api.rawResponse = body

	return nil
}

func isEchoNullValue(s string) bool {
	return s == "null" || s == "ul"
}

// The ECHO API returns JSON which only contains string values.
// The following methods are implementations of the UnmarshalJSON interface
// which unmarshal string results to proper data types.
func (eb *echoBool) UnmarshalJSON(b []byte) error {
	boolStr := string(b[1 : len(b)-1])
	boolVal := boolStr == "Y" || boolStr == "Yes"

	*eb = echoBool(boolVal)

	return nil
}

func (ed *echoDate) UnmarshalJSON(b []byte) error {
	dateStr := string(b[1 : len(b)-1])

	if isEchoNullValue(dateStr) {
		dateStr = echoNullDateValue
	}

	d, err := time.Parse(echoDefaultDateFormat, dateStr)
	if err != nil {
		return err
	}

	*ed = echoDate(d)

	return nil
}

func (ef *echoFloat) UnmarshalJSON(b []byte) error {
	floatStr := string(b[1 : len(b)-1])

	if isEchoNullValue(floatStr) {
		*ef = echoFloat(0)
	} else {
		f, err := strconv.ParseFloat(floatStr, 64)

		if err != nil {
			return err
		}

		*ef = echoFloat(f)
	}

	return nil
}

func (ei *echoInt) UnmarshalJSON(b []byte) error {
	intStr := string(b[1 : len(b)-1])

	if isEchoNullValue(intStr) {
		*ei = echoInt(0)
	} else {
		i, err := strconv.Atoi(intStr)
		if err != nil {
			return err
		}

		*ei = echoInt(i)
	}

	return nil
}
