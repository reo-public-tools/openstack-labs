package theforeman

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
)


// Struct for holding the api query results for compute profiles
type HostGroupQueryResults struct {
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	ID                  int    `json:"id"`
	Name                string `json:"name"`
        ComputeProfileName  string `json:"compute_profile_name"`
        ComputeProfileID    int    `json:"compute_profile_id"`
        ComputeResourceName string `json:"compute_resource_name"`
        ComputeResourceID   int    `json:"compute_resource_id"`
        ArchitectureName    string `json:"architecture_name"`
        ArchitectureID      int    `json:"architecture_id"`
}

func GetHostgroupInfo(url string, session string, hostGroup interface{}) (HostGroupQueryResults, error) {

    sysLogPrefix := "theforeman(package).hostgroups(file).ConvHostGroupNameToID(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting host info for \"%v\"", sysLogPrefix, hostGroup))

    // var for holding the query results
    var queryResults HostGroupQueryResults

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/hostgroups/%v", url, hostGroup)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", requesturl, nil)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return HostGroupQueryResults{}, err
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
        return HostGroupQueryResults{}, err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return HostGroupQueryResults{}, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return HostGroupQueryResults{}, err
    }

    return queryResults, nil

}

func ConvHostGroupNameToID(url string, session string, hostGroupName string) (int, error) {

    sysLogPrefix := "theforeman(package).hostgroups(file).ConvHostGroupNameToID(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting host group id from name \"%s\"", sysLogPrefix, hostGroupName))

    // var for holding the query results
    var queryResults HostGroupQueryResults

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/hostgroups/%s", url, hostGroupName)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", requesturl, nil)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return 0, err
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
        return 0, err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return 0, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return 0, err
    }

    return queryResults.ID, nil

}


