package theforeman

import (
    "fmt"
    "bytes"
    "net/http"
    "io/ioutil"
    "crypto/tls"
    "encoding/json"
)


// Structures for pulling a domain listing
type Sort struct {
	By    string `json:"by"`
	Order string `json:"order"`
}

type DomainQueryResults struct {
	Total    int         `json:"total"`
	Subtotal int         `json:"subtotal"`
	Page     int         `json:"page"`
	PerPage  int         `json:"per_page"`
	Search   string      `json:"search"`
	Sort     Sort        `json:"sort"`
	Results  []Domain    `json:"results"`
}

type Domain struct {
	Fullname  string      `json:"fullname"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	DNSID     int         `json:"dns_id"`
	DNS       string      `json:"dns"`
}


// Structures for pulling info for a single domain
type DomainInfo struct {
	Fullname      string          `json:"fullname"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	DNSID         int             `json:"dns_id"`
	DNS           string          `json:"dns"`
	Subnets       []Subnets       `json:"subnets"`
	Interfaces    []interface{}   `json:"interfaces"`
	Parameters    []Parameters    `json:"parameters"`
	Locations     []Locations     `json:"locations"`
	Organizations []Organizations `json:"organizations"`
}
type Subnets struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	NetworkAddress string `json:"network_address"`
}
type Parameters struct {
	Priority      int         `json:"priority"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	ParameterType string      `json:"parameter_type"`
	Value         interface{} `json:"value"`
}
type Locations struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
type Organizations struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Structures used when creating a new domain
type DomainPostData struct {
        OrganizationID int           `json:"organization_id"`
        LocationID     int           `json:"location_id"`
        Domain         NewDomainData `json:"domain"`
}
type DomainParametersAttributes struct {
        Name  string `json:"name"`
        Value string `json:"value"`
}
type NewDomainData struct {
        Name                       string                       `json:"name"`
        Fullname                   string                       `json:"fullname"`
        DomainParametersAttributes []DomainParametersAttributes `json:"domain_parameters_attributes"`
}


func GetDomains(url string, session string) ([]Domain, error) {

    sysLogPrefix := "theforeman(package).domains(file).GetDomains(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting a non-detailed list of domains.", sysLogPrefix))


    // Initialize the page and entries per page
    var entries_per_page int = 20
    var curpage int = 1
    var processed int = 0

    // Track current query and overall domain array
    domainList := []Domain{}
    var queryResults DomainQueryResults


    // Pager loop
    for {

        // Set the query url
        var requesturl string = fmt.Sprintf("%s/api/domains?per_page=%d&page=%d", url, entries_per_page, curpage)

        // Set up the basic request from the url and body
        req, err := http.NewRequest("GET", requesturl, nil)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return domainList, err
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
            return domainList, err
        }
        defer resp.Body.Close()

        // Read in the body and check status
        body, _ := ioutil.ReadAll(resp.Body)
        if resp.StatusCode != 200 {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
            return domainList, fmt.Errorf("%s %s", sysLogPrefix, string(body))
        }

        // Convert the body to a byte array
        bytes := []byte(body)

        // Unmarshall the json byte array into a struct
        err = json.Unmarshal(bytes, &queryResults)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return domainList, err
        }

        // Append the results to the domain list
        domainList = append(domainList, queryResults.Results...)

        // Do the pager calculations. 
        processed += len(queryResults.Results)
        if (processed < queryResults.Subtotal) {
            curpage += 1
            continue
        } else {
            break
        }

    }

    return domainList, nil

}


func GetDomainDetails(url string, session string, domainid interface{}) (DomainInfo, error) {

    sysLogPrefix := "theforeman(package).domains(file).GetDomainDetails(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting details for domain %v", sysLogPrefix, domainid))

    // Init some vars
    domainInfo := DomainInfo{}

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/domains/%v", url, domainid)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", requesturl, nil)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return domainInfo, err
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
        return domainInfo, err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return domainInfo, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &domainInfo)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return domainInfo, err
    }

    return domainInfo, nil

}

func SetDomainParameter(url string, session string, domainid interface{}, paramkey string, paramvalue interface{}) (error) {

    sysLogPrefix := "theforeman(package).domains(file).SetDomainParameter(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Setting domain(%v) parameter key(%s) value(%v)", sysLogPrefix, domainid, paramkey, paramvalue))

    // Get a fresh set of domain details
    curdomaininfo, err := GetDomainDetails(url, session, domainid)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }


    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/domains/%v/parameters", url, domainid)
    var data string = fmt.Sprintf("{\"parameter\": {\"name\": \"%s\", \"value\": \"%v\"}}", paramkey, paramvalue)
    var action string = "POST"

    // See if the parameter exists
    pExists := 0
    for _, parameter := range curdomaininfo.Parameters {
        if parameter.Name == paramkey {
            pExists = 1
        }
    }

    // Modify the request url to account for an existing parameter for this key
    if pExists == 1 {
        requesturl = fmt.Sprintf("%s/api/domains/%v/parameters/%s", url, domainid, paramkey)
        data = fmt.Sprintf("{\"value\": \"%v\"}", paramvalue)
        action = "PUT"
    }

    // Set up the basic request from the url and body
    req, err := http.NewRequest(action, requesturl, bytes.NewBufferString(data))
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
    if !(resp.StatusCode == 200 || resp.StatusCode == 201) {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    return nil

}

func DeleteDomainParameter(url string, session string, domainid interface{}, paramkey string) (error) {

    sysLogPrefix := "theforeman(package).domains(file).DeleteDomainParameter(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Deleting domain(%v) parameter key(%s)", sysLogPrefix, domainid, paramkey))

    // Get a fresh set of domain details
    curdomaininfo, err := GetDomainDetails(url, session, domainid)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }


    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/domains/%v/parameters/%s", url, domainid, paramkey)
    var action string = "DELETE"

    pExists := 0
    for _, parameter := range curdomaininfo.Parameters {
        if parameter.Name == paramkey {
            pExists = 1
        }
    }
    if pExists == 0 {
        return nil
    }

    // Set up the basic request from the url and body
    req, err := http.NewRequest(action, requesturl, nil)
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
    if !(resp.StatusCode == 200 || resp.StatusCode == 201) {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    return nil

}

func GetDomainParameter(url string, session string, domainid interface{}, paramkey string) (string, error) {

    sysLogPrefix := "theforeman(package).domains(file).GetDomainParameter(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting domain(%v) parameter key(%s)", sysLogPrefix, domainid, paramkey))

    // Get a fresh set of domain details
    curdomaininfo, err := GetDomainDetails(url, session, domainid)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    // See if the parameter exists
    for _, parameter := range curdomaininfo.Parameters {
        if parameter.Name == paramkey {
            return parameter.Value.(string), nil
        }
    }

    return "", fmt.Errorf("Parameter key %s for domain %s is not currently set", paramkey, domainid)

}

func GetDomainsWithDetails(url string, session string) ([]DomainInfo, error) {

    sysLogPrefix := "theforeman(package).domains(file).GetDomainsWithDetails(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting a list of domains including details.", sysLogPrefix))

    // List to return
    domainInfoList := []DomainInfo{}

    // Get a full foreman domain listing
    domains, err := GetDomains(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return domainInfoList, err
    }

    // Parse through the listing and get the detailed domain for each
    for _, domain := range domains {

        // Ignore the default domain
        if domain.Name == "localdomain" {
            continue
        }

        // Get the details for the current domain
        curDomainInfo, err := GetDomainDetails(url, session, domain.ID)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return domainInfoList, err
        }

        // Append the domain info to the list to return
        domainInfoList = append(domainInfoList, curDomainInfo)

    }

    return domainInfoList, nil
}

func CreateNewDomain(url string, session string, domainData DomainPostData) (Domain, error) {

    sysLogPrefix := "theforeman(package).domains(file).CreateNewDomain(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating a new domain %s.", sysLogPrefix, domainData.Domain.Name))

    // Init the query results var
    var queryResults Domain

    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/domains/", url)

    // Convert data to json
    postData, err := json.Marshal(domainData)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Domain{}, err
    }

    // Set up the basic request from the url and body
    req, err := http.NewRequest("POST", requesturl, bytes.NewBufferString(string(postData)))
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Domain{}, err
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
        return Domain{}, err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 201 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return Domain{}, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Domain{}, err
    }

    return queryResults, nil
}

func DeleteDomain(url string, session string, domainName string) (error) {

    sysLogPrefix := "theforeman(package).domains(file)DeleteDomain(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Deleting domain %s.", sysLogPrefix, domainName))

    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/domains/%s", url, domainName)

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
        return fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    return nil
}
