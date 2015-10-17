package epaecho

import (
	"encoding/json"
	"fmt"
)

const (
	// facilitiesAPIURI is the unique URL path for a Get Facilities API request.
	facilitiesAPIURI      = "echo_rest_services.get_facilities"
	defaultResultsPerPage = "100"
)

var (
	// echoFaclitiesAPIOptions represents a human-readable map of facilities
	// query options.
	echoFacilitiesAPIOptions = map[string]string{
		"resultsPerPage":                    "responseset",
		"passthrough":                       "passthrough",
		"name":                              "p_fn",
		"address":                           "p_sa",
		"city":                              "p_ct",
		"county":                            "p_co",
		"fipsCode":                          "p_fips",
		"state":                             "p_st",
		"zip":                               "p_zip",
		"permitId":                          "p_pid",
		"region":                            "p_reg",
		"naicsCode":                         "p_ncs",
		"penaltyTime":                       "p_pen",
		"latitude":                          "p_c1lat",
		"longitude":                         "p_c1lon",
		"latitude2":                         "p_c2lat",
		"longitude2":                        "p_c2lon",
		"USMexicoBorderFacility":            "p_usmex",
		"SICCode":                           "p_sic2",
		"SICCode4":                          "p_sic4",
		"federalAgencyCode":                 "p_fa",
		"federalFacility":                   "p_ff",
		"activeFacility":                    "p_act",
		"majorFacility":                     "p_maj",
		"cleanAirMACT":                      "p_mact",
		"formalEnforcementWithinYearsAgo":   "p_fea",
		"formalEnforcementActionYearsAgo":   "p_feay",
		"FEAAAgencyCode":                    "p_feaa",
		"FEAARestrictedToICIS":              "p_feac",
		"informalEnforcementWithinYearsAgo": "p_iea",
		"informalEnforcementActionYearsAgo": "p_ieay",
		"informalEnforcementAgency":         "p_ieaa",
		"currentComplianceStatus":           "p_cs",
		"quartersInViolation":               "p_qiv",
		"nonAttainmentArea":                 "p_naa",
		"impairedWatersCategory":            "p_impw",
		"TRIReporterCurrent":                "p_trep",
		"onsiteChemicalReleases":            "p_ocr",
		"offsiteChemicalTransfers":          "p_oct",
		"percentMinority":                   "p_pm",
		"populationDensity":                 "p_pd",
		"lowIncome":                         "p_li",
		"nativeAmericanTerritory":           "p_ico",
		"watershedHUC":                      "p_huc",
		"mediaType":                         "p_med",
		"lastInspectionWithinYearsAgo":      "p_ysl",
		"lastINspectionYearsAgo":            "p_ysly",
		"inspectionAgency":                  "p_ysla",
		"quickSearch":                       "p_qs",
		"singleFacility":                    "p_sfs",
		"tribalId":                          "p_tribeid",
		"tribalDistance":                    "p_tribedist",
	}

	defaultFacilitiesParams = map[string]string{
		"resultsPerPage": defaultResultsPerPage,
	}
)

type EchoFacilitiesAPI struct {
	Config   *EchoAPIConfig
	response echoFacilitiesResponse
}

type echoFacilitiesResponse struct {
	Results struct {
		Error                          EchoAPIError   `json:"Error"`
		Message                        string         `json:"Message"`
		QueryRows                      echoInt        `json:",omitempty"`
		NonComplianceResults           echoInt        `json:"CVRows,omitempty"`
		QuartersNonComplianceResults   echoInt        `json:"V3Rows,omitempty"`
		FormalEnforcementActionResults echoInt        `json:"FEARows,omitempty"`
		CleanAirActRecordResults       echoInt        `json:"CAARows,omitempty"`
		CleanWaterActRecordResults     echoInt        `json:"CWARows,omitempty"`
		HazardousWasteRecordResults    echoInt        `json:"RCRRows,omitempty"`
		ToxicReleaseInventoryResults   echoInt        `json:"TRIRows,omitempty"`
		QueryID                        string         `json:",omitempty"`
		PageNo                         string         `json:"PageNo,omitempty"`
		Facilities                     []EchoFacility `json:",omitempty"`
		MapOutput                      struct {
			QueryID      string
			IconBaseURL  string
			PopUpBaseURL string
			MapData      []EchoMapData
		}
	}
}

// NewFacilitiesAPI is a constructor function which returns a pointer
// to an EchoFacilitiesAPI struct.
// Any values defined in the defaultFacilitiesParams map will be set during
// construction.
func NewFacilitiesAPI() *EchoFacilitiesAPI {
	facilitiesAPI := &EchoFacilitiesAPI{
		Config: newEchoAPIConfig(),
	}

	facilitiesAPI.Config.endpointURI = facilitiesAPIURI

	for key, val := range defaultFacilitiesParams {
		facilitiesAPI.SetParam(key, val)
	}

	return facilitiesAPI
}

func (api *EchoFacilitiesAPI) SetParam(key string, value string) error {
	err := api.Config.addParam(key, value, echoFacilitiesAPIOptions)
	if err != nil {
		return err
	}

	return nil
}

// Facilities returns the response from an EchoFacilities API request.
func (api *EchoFacilitiesAPI) Facilities() (echoFacilitiesResponse, error) {
	err := api.Config.doRequest()
	if err != nil {
		return echoFacilitiesResponse{}, err
	}

	unmarshalErr := api.unmarshalResponse()
	if unmarshalErr != nil {
		return echoFacilitiesResponse{}, unmarshalErr
	}

	return api.response, nil
}

func (api *EchoFacilitiesAPI) unmarshalResponse() error {
	response := echoFacilitiesResponse{}

	err := json.Unmarshal(api.Config.rawResponse, &response)
	if err != nil {
		return err
	}

	if len(response.Results.Error.ErrorMessage) > 0 {
		return fmt.Errorf("%s. Message was: %s", echoAPIError, response.Results.Error.ErrorMessage)
	}

	api.response = response

	return nil
}
