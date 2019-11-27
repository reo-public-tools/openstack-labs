package theforeman

import (
    "fmt"
    "errors"
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

// Structures for pulling a domain listing
func GetCommonParameters(url string, session string) ([]CommonParameter, error) {

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
            return commonParameterList, err
        }
        defer resp.Body.Close()

        // Read in the body and check status
        body, _ := ioutil.ReadAll(resp.Body)
        if resp.StatusCode != 200 {
            return commonParameterList, errors.New(string(body))
        }

        // Convert the body to a byte array
        bytes := []byte(body)

        // Unmarshall the json byte array into a struct
        err = json.Unmarshal(bytes, &queryResults)
        if err != nil {
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
