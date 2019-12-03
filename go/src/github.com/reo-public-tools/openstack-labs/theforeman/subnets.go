package theforeman

import (
    "fmt"
    "bytes"
    "errors"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
)

// Structures used for creating subnets
type SubnetPostData struct {
	OrganizationID int    `json:"organization_id"`
	LocationID     int    `json:"location_id"`
	Subnet         Subnet `json:"subnet"`
}
type Subnet struct {
	Name                       string                       `json:"name"`
	NetworkType                string                       `json:"network_type"`
	Network                    string                       `json:"network"`
	Mask                       string                       `json:"mask"`
	Gateway                    string                       `json:"gateway"`
	From                       string                       `json:"from"`
	To                         string                       `json:"to"`
	Ipam                       string                       `json:"ipam"`
	BootMode                   string                       `json:"boot-mode"`
	DomainIds                  []int                        `json:"domain_ids"`
	SubnetParametersAttributes []SubnetParametersAttributes `json:"subnet_parameters_attributes"`
}
type SubnetParametersAttributes struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func CreateSubnet(url string, session string, subnetData SubnetPostData) (Subnet, error) {

    sysLogPrefix := "theforeman(package).subnets(file).CreateSubnet(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating subnet %s.", sysLogPrefix, subnetData.Subnet.Name))

    // Init the query results
    var queryResults Subnet

    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/subnets/", url)

    // Convert data to json
    postData, err := json.Marshal(subnetData)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Subnet{}, err
    }

    // Set up the basic request from the url and body
    req, err := http.NewRequest("POST", requesturl, bytes.NewBufferString(string(postData)))
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Subnet{}, err
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
        return Subnet{}, err
    }
    defer resp.Body.Close()


    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 201 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return Subnet{}, errors.New(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Subnet{}, err
    }

    return queryResults, nil
}

func RemoveSubnetFromDomain(url string, session string, subnetID interface{}) (error) {

    sysLogPrefix := "theforeman(package).subnets(file).RemoveSubnetFromDomain(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Clearing subnet %v domain list to prep for removal.", sysLogPrefix, subnetID))

    requesturl := fmt.Sprintf("%s/api/subnets/%v", url, subnetID)
    data := "{\"subnet\": {\"domain_ids\": []}}"

    // Set up the basic request from the url and body
    req, err := http.NewRequest("PUT", requesturl, bytes.NewBufferString(data))
    if err != nil {
        return err
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
        return err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return errors.New(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
    }

    return nil
}

func DeleteSubnet(url string, session string, subnetID interface{}) (error) {

    sysLogPrefix := "theforeman(package).subnets(file).DeleteSubnet(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Deleting subnet %v.", sysLogPrefix, subnetID))

    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/subnets/%v", url, subnetID)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("DELETE", requesturl, nil)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
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
        return err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return errors.New(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
    }

    return nil
}
