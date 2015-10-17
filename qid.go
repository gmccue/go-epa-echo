package epaecho

import (
	"encoding/json"
	"fmt"
)

// queryAPIURI is the unique URL path for a QID API request.
const queryAPIURI = "echo_rest_services.get_qid"

type EchoQueryAPI struct {
	Config   *EchoAPIConfig
	response echoQueryResponse
}

var (
	// echoQueryAPIOptions represents a human-readable map of QID query options.
	echoQueryAPIOptions = map[string]string{
		"queryID":    "qid",
		"pageNumber": "pageno",
	}

	defaultQueryParams = map[string]string{
		"pageNumber": "1",
	}
)

type echoQueryResponse struct {
	Message    string       `json:"Message"`
	Error      EchoAPIError `json:"Error"`
	QueryRows  echoInt      `json:"QueryRows"`
	QueryID    echoInt      `json:"QueryID"`
	PageNo     echoInt      `json:"PageNo"`
	Facilities []EchoFacility
}

// NewQueryAPI is a constructor function which returns a pointer
// to an EchoQueryAPI struct.
// Any values defined in the defaultQueryParams map will be set during
// construction.
func NewQueryAPI(queryID string) *EchoQueryAPI {
	queryAPI := &EchoQueryAPI{
		Config: newEchoAPIConfig(),
	}

	queryAPI.Config.endpointURI = queryAPIURI
	queryAPI.SetParam("queryID", queryID)

	for key, val := range defaultQueryParams {
		queryAPI.SetParam(key, val)
	}

	return queryAPI
}

func (api *EchoQueryAPI) SetParam(key string, value string) error {
	err := api.Config.addParam(key, value, echoQueryAPIOptions)
	if err != nil {
		return err
	}

	return nil
}

func (api *EchoQueryAPI) Results() (echoQueryResponse, error) {
	err := api.Config.doRequest()
	if err != nil {
		return echoQueryResponse{}, err
	}

	unmarshalErr := api.unmarshalResponse()
	if unmarshalErr != nil {
		return echoQueryResponse{}, unmarshalErr
	}

	return api.response, nil
}

func (api *EchoQueryAPI) unmarshalResponse() error {
	response := echoQueryResponse{}

	err := json.Unmarshal(api.Config.rawResponse, &response)
	if err != nil {
		return err
	}

	if len(response.Error.ErrorMessage) > 0 {
		return fmt.Errorf("%s. Message was: %s", echoAPIError, response.Error.ErrorMessage)
	}

	api.response = response

	return nil
}
