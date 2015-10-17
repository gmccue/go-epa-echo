# go-epa-echo - A Go wrapper for the EPA's ECHO "All Data Facility Search" API.

[![Build Status](https://api.travis-ci.org/gmccue/go-epa-echo.png?branch=master)](https://travis-ci.org/gmccue/go-epa-echo)
[![GoDoc](https://godoc.org/github.com/gmccue/go-epa-echo?status.svg)](https://godoc.org/github.com/gmccue/go-epa-echo)

go-epa-echo provides programmatic access to the [EPA's Enforcement and Compliance Online (ECHO) All Data Facility Search API](http://echo.epa.gov/).

## Installation

```
go get github.com/gmccue/go-epa-echo
```

## Usage

By default, retrieving full data from an API request requires three steps:

1. Send a "Get Facilities" API request. This is the top-level API request, and starting point for additional queries. You can set query parameters for this request by using the `SetParam()` method. A detailed list of available query parameters is available on the [Get Facilities query parameters wikie page](https://github.com/gmccue/go-epa-echo/wiki/Get-Facilities-Query-Parameters). Currently only the parameters listed in this document are supported by go-epa-echo.<br>You can optionally set the `passthrough` parameter to `"Y"` (`api.SetParam("passthrough", "Y")`) which will return all matching facility and map data in a single response. Pagination is not supported in "passthrough" mode.<br>This request will return an [echoFacilitiesRespones data structure](https://github.com/gmccue/go-epa-echo/wiki/Get-Facilities-API-Result-Fields).

2. Send a "Get QID" API request, using the query ID returned from step 1. This API returns more detailed, paginated facility information related to the base query. You can set query parameters for this request (such as page number) by using the `SetParam()` method. A detailed list of available query parameters is available on the [Get QID query parameters wiki page](https://github.com/gmccue/go-epa-echo/wiki/Get-QID-Query-Parameters).<br>This request will return an [echoQueryResponse data structure](https://github.com/gmccue/go-epa-echo/wiki/Get-QID-API-Result-Fields).

3. If desired, send a "Get Map Data" API request. This API returns detailed mapping data for all facilities. There is an upper-limit of 500 results per request for this API.<br>This request will return an [echoMapResponse data structure](https://github.com/gmccue/go-epa-echo/wiki/Get-Map-Data-API-Result-Fields).

An example using all three API requests might look something like this:

```
import (
	"log"

	echo github.com/gmccue/go-epa-echo
)

func test() {
	fAPI := echo.NewFacilitiesAPI()
	fAPI.Config.Debug = true
	fAPI.SetParam("resultsPerPage", "1")
	fAPI.SetParam("city", "Baltimore")
	fAPI.SetParam("state", "MD")

	res, err := fAPI.Facilities()
	if err != nil {
		log.Println(err)
	}

	qAPI := echo.NewQueryAPI(res.Results.QueryID)
	qAPI.Config.Debug = true
	qAPI.SetParam("pageNumber", "1")

	qres, qerr := qAPI.Results()
	if qerr != nil {
		log.Println(qerr)
	}

	mAPI := echo.NewMapAPI(res.Results.QueryID)
	mAPI.Config.Debug = true

	mres, merr := mAPI.Maps()
	if merr != nil {
		log.Println(merr)
	}
}
```

### Configurable fields
| Field   | Description                                                            | Example |
|---------|------------------------------------------------------------------------|---------|
| Debug   | Output detailed information related to an API request. Uses pkg `log`. | ap.Debug(true)
| Timeout | The HTTP request timeout (in seconds).                                 | api.Timeout(30)

### Get Facilities API

Create a new Facilities API request with `api := echo.NewFacilitiesAPI()`

Set a query parameters with 
```
api.SetParam("city", "baltimore")
api.SetParam("state", "MD")
api.SetParam("resultsPerPage", "10")
```

Get the facilities API results:
```
res, err := fAPI.Facilities()
if err != nil {
    log.Println(err)
}
```

A list of returned data structure fields is available on the [facilities response struct wiki page](https://github.com/gmccue/go-epa-echo/wiki/Get-Facilities-API-Result-Fields).
A list of all available query parameters is available on the [Get Facilities query parameter wiki page](https://github.com/gmccue/go-epa-echo/wiki/Get-Facilities-Query-Parameters).


### Get QID API

You can retrieve detailed facilities results from a previously run Get Facilties query using the QID API.

```
queryID := "1"
qAPI := echo.NewQueryAPI(queryID)

qres, qerr := qAPI.Results()
if qerr != nil {
	log.Println(qerr)
}

```

A list of returned data structure fields is available on the [Get QID response struct wiki page](https://github.com/gmccue/go-epa-echo/wiki/Get-QID-API-Result-Fields).
A list of all available query parameters is available on the [Get QID query parameter wiki page](https://github.com/gmccue/go-epa-echo/wiki/Get-QID-Query-Parameters).


### Get Map Data API

You can retrieve detailed mapping data for a previously run Get Facilities query using the Map Data API.

```
queryID := "1"
mAPI := echo.NewMapAPI(queryID)

mres, merr := mAPI.Maps()
if merr != nil {
	log.Println(merr)
}
```

A list of returned data structure fields is available on the [map data response struct wiki page](https://github.com/gmccue/go-epa-echo/wiki/Get-Map-Data-API-Result-Fields).

## Detailed API Documentation

Detailed documentation for the All Data Facility Search API is [available online here](http://echo.epa.gov/system/files/ECHO%20All%20Data%20Search%20Services_v3.pdf).


## FAQ

FAQ related to the Echo API can be found here: [Frequently Asked Questions](https://echo.epa.gov/resources/general-info/echo-faq) 


## Running tests

All tests can be run with the command `go test ./... -v`
