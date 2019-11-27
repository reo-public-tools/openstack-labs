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
	Priority      int    `json:"priority"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ID            int    `json:"id"`
	Name          string `json:"name"`
	ParameterType string `json:"parameter_type"`
	Value         string `json:"value"`
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


func GetDomains(url string, session string) ([]Domain, error) {

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
            return domainList, err
        }
        defer resp.Body.Close()

        // Read in the body and check status
        body, _ := ioutil.ReadAll(resp.Body)
        if resp.StatusCode != 200 {
            return domainList, errors.New(string(body))
        }

        // Convert the body to a byte array
        bytes := []byte(body)

        // Unmarshall the json byte array into a struct
        err = json.Unmarshal(bytes, &queryResults)
        if err != nil {
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


func GetDomainDetails(url string, session string, domainid int) (DomainInfo, error) {

    // Init some vars
    domainInfo := DomainInfo{}

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/domains/%d", url, domainid)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", requesturl, nil)
    if err != nil {
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
        return domainInfo, err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        return domainInfo, errors.New(string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &domainInfo)
    if err != nil {
        return domainInfo, err
    }

    return domainInfo, nil

}

func SetDomainParameter(url string, session string, domainid int, paramkey string, paramvalue interface{}) (error) {

    // Get a fresh set of domain details
    curdomaininfo, err := GetDomainDetails(url, session, domainid)
    if err != nil {
        return err
    }


    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/domains/%d/parameters", url, domainid)
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
        requesturl = fmt.Sprintf("%s/api/domains/%d/parameters/%s", url, domainid, paramkey)
        data = fmt.Sprintf("{\"value\": \"%v\"}", paramvalue)
        action = "PUT"
    }

    // Set up the basic request from the url and body
    req, err := http.NewRequest(action, requesturl, bytes.NewBufferString(data))
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
        return err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        return errors.New(string(body))
    }

    return nil

}

