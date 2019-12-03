package theforeman

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
)

type LocationsQueryResults struct {
	Total    int                 `json:"total"`
	Subtotal int                 `json:"subtotal"`
	Page     int                 `json:"page"`
	PerPage  int                 `json:"per_page"`
	Search   string              `json:"search"`
	Sort     Sort                `json:"sort"`
	Results  []Location          `json:"results"`
}
// Defined in domains.go
/*
type Sort struct {
	By    interface{} `json:"by"`
	Order interface{} `json:"order"`
}
*/
type Location struct {
	Ancestry    interface{} `json:"ancestry"`
	ParentID    int         `json:"parent_id"`
	ParentName  string      `json:"parent_name"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
}

// Get the id of a location
func ConvLocNameToID(url string, session string, locName string) (int, error) {

    sysLogPrefix := "theforeman(package).locations(file).ConvLocNameToID(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting id for location name \"%s\"", sysLogPrefix, locName))

    locations, err := GetLocations(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return 1, err
    }
    for _, location := range locations {
        if location.Name == locName {
            _ = sysLog.Debug(fmt.Sprintf("%s Found id %d for location name %s", sysLogPrefix, location.ID, locName))
            return location.ID, nil
        }
    }

    return 1, fmt.Errorf("%s Location %s not found in ConvLocNameToID", sysLogPrefix, locName)
}

// Structures for pulling a domain listing
func GetLocations(url string, session string) ([]Location, error) {

    sysLogPrefix := "theforeman(package).locations(file).GetLocations(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting a list of locations.", sysLogPrefix))

    // Initialize the page and entries per page
    var entries_per_page int = 20
    var curpage int = 1
    var processed int = 0

    // Track current query and overall domain array
    locationList := []Location{}
    var queryResults LocationsQueryResults


    // Pager loop
    for {

        // Set the query url
        var requesturl string = fmt.Sprintf("%s/api/locations?per_page=%d&page=%d", url, entries_per_page, curpage)

        // Set up the basic request from the url and body
        req, err := http.NewRequest("GET", requesturl, nil)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return locationList, err
        }

        // Make sure we are using the proper content type for the configs api
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json,version=2")

        // Set the session Cookie header
        req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", session))

        // Disable tls verify
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        // Set up the http client and do the request
        client := &http.Client{Transport: tr}
        resp, err := client.Do(req)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return locationList, err
        }
        defer resp.Body.Close()

        // Read in the body and check status
        body, _ := ioutil.ReadAll(resp.Body)
        if resp.StatusCode != 200 {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
            return locationList, fmt.Errorf("%s %s", sysLogPrefix, string(body))
        }

        // Convert the body to a byte array
        bytes := []byte(body)

        // Unmarshall the json byte array into a struct
        err = json.Unmarshal(bytes, &queryResults)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return locationList, err
        }

        // Append the results to the domain list
        locationList = append(locationList, queryResults.Results...)

        // Do the pager calculations. 
        processed += len(queryResults.Results)
        if (processed < queryResults.Subtotal) {
            curpage += 1
            continue
        } else {
            break
        }

    }

    return locationList, nil
}
