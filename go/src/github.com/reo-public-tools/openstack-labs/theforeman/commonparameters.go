package theforeman

import (
    "fmt"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
)


// Structure for holding global paramter api call data
type CommonParameterQueryResults struct {
	Total    int               `json:"total"`
	Subtotal int               `json:"subtotal"`
	Page     int               `json:"page"`
	PerPage  int               `json:"per_page"`
	Search   string            `json:"search"`
	Sort     Sort              `json:"sort"`
	Results  []CommonParameter `json:"results"`
}

// Declared in domains.go
/*
type Sort struct {
	By    string `json:"by"`
	Order string `json:"order"`
}
*/

type CommonParameter struct {
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	IsHiddenValue bool        `json:"hidden_value?"`
	HiddenValue   string      `json:"hidden_value"`
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	ParameterType string      `json:"parameter_type"`
	Value         string      `json:"value"`
}

type GlobalParameters struct {
        LabOrgName          string
        LabLocationName     string
        LabBaseDomainName   string
        MulticastGroupBase  string
        VXLANNetworkPrefix  string
        VXLANNetmask        string
        VXLANThirdOctetStep int
        VXLANNetworks       []string
}

func GetGlobalParameters(url string, session string) (GlobalParameters, error) {

    sysLogPrefix := "theforeman(package).commonparameters(file).GetGlobalParameters(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Get globalParameters struct based on commonparameter attributes", sysLogPrefix))


    // Get a list of common parameters that will be used when creating a domain
    commonParameters, err := GetCommonParameters(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return GlobalParameters{}, err
    }

    // Set some variables based on the parameters
    globalParams := GlobalParameters{}
    for _, parameter := range commonParameters {
        switch parameter.Name {
            case "lab_base_domain_name":
                globalParams.LabBaseDomainName = parameter.Value
            case "lab_org_name":
                globalParams.LabOrgName = parameter.Value
            case "lab_location_name":
                globalParams.LabLocationName = parameter.Value
            case "multicast_group_base":
                globalParams.MulticastGroupBase = parameter.Value
            case "vxlan_network_prefix":
                globalParams.VXLANNetworkPrefix = parameter.Value
            case "vxlan_netmask":
                globalParams.VXLANNetmask = parameter.Value
            case "vxlan_third_octet_step":
                globalParams.VXLANThirdOctetStep, err = strconv.Atoi(parameter.Value)
                if err != nil {
                    return GlobalParameters{}, err
                }
            case "vxlan_networks":
                globalParams.VXLANNetworks = strings.Split(parameter.Value, ",")
            default:
                continue
        }
    }

    return globalParams, nil
}

func GetCommonParameters(url string, session string) ([]CommonParameter, error) {

    sysLogPrefix := "theforeman(package).commonparameters(file).GetCommonParameters(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting a list of global(common) parameters.", sysLogPrefix))

    // Initialize the page and entries per page
    var entries_per_page int = 20
    var curpage int = 1
    var processed int = 0

    // Track current query and overall domain array
    commonParameterList := []CommonParameter{}
    var queryResults CommonParameterQueryResults


    // Pager loop
    for {

        // Set the query url
        var requesturl string = fmt.Sprintf("%s/api/common_parameters?per_page=%d&page=%d", url, entries_per_page, curpage)

        // Set up the basic request from the url and body
        req, err := http.NewRequest("GET", requesturl, nil)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return commonParameterList, err
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
            return commonParameterList, err
        }
        defer resp.Body.Close()

        // Read in the body and check status
        body, _ := ioutil.ReadAll(resp.Body)
        if resp.StatusCode != 200 {
            return commonParameterList, fmt.Errorf("%s %s", sysLogPrefix, string(body))
        }

        // Convert the body to a byte array
        bytes := []byte(body)

        // Unmarshall the json byte array into a struct
        err = json.Unmarshal(bytes, &queryResults)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return commonParameterList, err
        }

        // Append the results to the domain list
        commonParameterList = append(commonParameterList, queryResults.Results...)

        // Do the pager calculations. 
        processed += len(queryResults.Results)
        if (processed < queryResults.Subtotal) {
            curpage += 1
            continue
        } else {
            break
        }

    }

    return commonParameterList, nil
}
