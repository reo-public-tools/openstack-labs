package theforeman

import (
    "fmt"
    "errors"
    "strings"
)

const domainNamePrefix = "lab"
const maxLabs = 50

func CreateDynamicLab(url string, session string) (error) {

    // Get a list of common parameters that will be used when creating a domain
    commonParameters, err := GetCommonParameters(url, session)
    if err != nil {
        return err
    }

    // Set some variables based on the parameters
    baseDomain := ""
    labOrg := ""
    labLoc := ""
    for _, parameter := range commonParameters {
        switch parameter.Name {
            case "lab_base_domain_name":
                baseDomain = parameter.Value
            case "lab_org_name":
                labOrg = parameter.Value
            case "lab_location_name":
                labLoc = parameter.Value
            default:
                continue
        }
    }

    // Make call to create the domain in foreman
    err = CreateDynamicLabDomain(url, session, baseDomain, labOrg, labLoc)
    if err != nil {
        return err
    }

    return nil
}

func CreateDynamicLabDomain(url string, session string, baseDomain string, labOrg string, labLoc string) (error) {

    // Find an available dynamic lab domain slot
    domainName, domainIndex, err := FindAvailableLabSlot(url, session, baseDomain)
    if err != nil {
        return err
    }

    fmt.Printf("Creating domain %s with index of %d\n", domainName, domainIndex)
    description := fmt.Sprintf("%s%d Dynamic Domain", strings.ToUpper(domainNamePrefix), domainIndex)
    fmt.Printf("Desc: %s\n", description)

    return nil
}

func FindAvailableLabSlot(url string, session string, baseDomain string) (string, int, error) {

    // Get a full foreman domain listing
    domains, err := GetDomains(url, session)
    if err != nil {
        return "", 0, err
    }

    // Track the index
    domainIndex := 1

    // loop until we find something free
    for domainIndex <= maxLabs {

        // Track if a match is found
        matchFound := 0

        // Set up match string
        domainMatch := fmt.Sprintf("%s%d.%s", domainNamePrefix, domainIndex, baseDomain)

        // Parse through the listing
        for _, domain := range domains {

            // Ignore the default domain
            if domain.Name == "localdomain" {
                continue
            }

            // Set up the name match
            if domain.Name == domainMatch {
                matchFound = 1
                break
            }

        }

        // No match was found if we made it this far
        if matchFound == 0 {
            return domainMatch, domainIndex, nil
        } else {
            domainIndex += 1
        }

    } // End "for domainIndex <= maxLabs {"

    return "", 0, errors.New("We hit the max lab limit without finding an available slow")

}
