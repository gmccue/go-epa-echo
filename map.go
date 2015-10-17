package epaecho

import (
	"encoding/json"
)

// mapAPIURI is the unique URL path for a Get Map Data API request.
const mapAPIURI = "echo_rest_services.get_map"

// echoMapAPIOptions represents a human-readable map of map query options.
var echoMapAPIOptions = map[string]string{
	"queryID":    "QID",
	"Latitude1":  "C1Lat",
	"Longitude1": "C1Long",
	"Latitude2":  "C2Lat",
	"Longitude2": "C2Long",
}

type EchoMapAPI struct {
	Config   *EchoAPIConfig
	response echoMapResponse
}

type echoMapResponse struct {
	MapOutput struct {
		QueryID         string
		IconBaseURL     string
		PopUpBaseURL    string
		CenterLatitude  string
		CenterLongitude string
		MapData         []EchoMapData
	}
}

type EchoMapData struct {
	Latitude                string `json:"LAT"`
	Longitude               string `json:"LON"`
	Icon                    string `json:"ICON,omitempty"`
	Type                    string `json:"TYPE,omitempty"`
	Name                    string `json:"NAME,omitempty"`
	PUV                     string `json:"PUV,omitempty"`
	SDWAStatus              string `json:"SDWAstatus,omitempty"`
	RCRAStatus              string `json:"RCRAstatus,omitempty"`
	CWAStatus               string `json:"CWAstatus,omitempty"`
	CAAStatus               string `json:"CAAstatus,omitempty"`
	LastInspection          string `json:"LastInsp,omitempty"`
	FormalInspectionCount   string `json:"FormalCount,omitempty"`
	InformalInspectionCount string `json:"InformalCount,omitempty"`
}

// NewMapAPI is a constructor function which returns a pointer
// to an EchoFacilitiesAPI struct.
// Any values defined in the defaultQueryParameters map will be set during
// construction.
func NewMapAPI(qid string) *EchoMapAPI {
	mapAPI := &EchoMapAPI{
		Config: newEchoAPIConfig(),
	}

	mapAPI.Config.endpointURI = mapAPIURI
	mapAPI.SetParam("queryID", qid)

	return mapAPI
}

func (api *EchoMapAPI) SetParam(key string, value string) error {
	err := api.Config.addParam(key, value, echoMapAPIOptions)
	if err != nil {
		return err
	}

	return nil
}

// Maps returns the response from a Get Map Data API request.
func (api *EchoMapAPI) Maps() (echoMapResponse, error) {
	err := api.Config.doRequest()
	if err != nil {
		return echoMapResponse{}, err
	}

	unmarshalErr := api.unmarshalResponse()
	if unmarshalErr != nil {
		return echoMapResponse{}, unmarshalErr
	}

	return api.response, nil
}

func (api *EchoMapAPI) unmarshalResponse() error {
	response := echoMapResponse{}

	err := json.Unmarshal(api.Config.rawResponse, &response)
	if err != nil {
		return err
	}

	api.response = response

	return nil
}
