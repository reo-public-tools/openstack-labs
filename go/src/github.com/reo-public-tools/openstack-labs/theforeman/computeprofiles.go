package theforeman

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
)


// Struct for holding the api query results for compute profiles
type ComputeProfileQueryResults struct {
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	ID                int    `json:"id"`
	Name              string `json:"name"`
	ComputeAttributes []struct {
		ID                   int    `json:"id"`
		Name                 string `json:"name"`
		ComputeResourceID    int    `json:"compute_resource_id"`
		ComputeResourceName  string `json:"compute_resource_name"`
		ProviderFriendlyName string `json:"provider_friendly_name"`
		ComputeProfileID     int    `json:"compute_profile_id"`
		ComputeProfileName   string `json:"compute_profile_name"`
	} `json:"compute_attributes"`
}

func ConvProfileNameToFlavorName(url string, session string, computeProfileName string) (string, error) {

    sysLogPrefix := "theforeman(package).computeprofiles(file).ConvProfileNameToFlavorName(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting flavor name for compute profile \"%s\"", sysLogPrefix, computeProfileName))

    // var for holding the query results
    var queryResults ComputeProfileQueryResults

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/compute_profiles/%s", url, computeProfileName)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", requesturl, nil)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
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
        return "", err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return "", fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    return queryResults.ComputeAttributes[0].Name, nil

}


