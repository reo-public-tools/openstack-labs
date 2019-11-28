package theforeman

import (
    "fmt"
    "time"
    "errors"
    "strconv"
    "strings"
    "math/rand"
)

const domainNamePrefix = "lab"
const maxLabs = 50


type GlobalParameters struct {
        LabOrgName          string
        LabLocationName     string
        LabBaseDomainName   string
        MulticastGroupBase  string
        VXLANNetworkPrefix  string
        VXLANNetmask        string
        VXLANThirdOctetStep int
}

func CreateDynamicLab(url string, session string) (error) {

    // Get a list of common parameters that will be used when creating a domain
    commonParameters, err := GetCommonParameters(url, session)
    if err != nil {
        return err
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
                    return err
                }
            default:
                continue
        }
    }

    // Make call to create the domain in foreman
    err = CreateDynamicLabDomain(url, session, globalParams)
    if err != nil {
        return err
    }

    return nil
}

func CreateDynamicLabDomain(url string, session string, globalParams GlobalParameters) (error) {

    // Find an available dynamic lab domain slot
    domainName, domainIndex, err := FindAvailableLabSlot(url, session, globalParams.LabBaseDomainName)
    if err != nil {
        return err
    }

    // Find a free multicast group
    multicastGroup, err := FindAvailableMulticastGroup(url, session, globalParams.MulticastGroupBase)
    if err != nil {
        return err
    }
    fmt.Printf("Multicast group %s\n", multicastGroup)

    // Request new external and internal haproxy/keepalived router ids
    // These should all be unique on a flat network.

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


func FindAvailableMulticastGroup(url string, session string, multicastGroupBase string) (string, error) {

    mcreturn := ""

    // Get a full foreman domain listing
    domains, err := GetDomainsWithDetails(url, session)
    if err != nil {
        return "", err
    }

    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())
    min := 1
    max := 254

    // Generate random group numbers and check it against the domain parameters
    isFreeForUse := false
    for !isFreeForUse  {

        // Generate a random number and define a new mc group
        randomNum := rand.Intn(max - min + 1) + min
        mcreturn = fmt.Sprintf("%s.%d", multicastGroupBase, randomNum)

        // Set it free here as it will be marked false if a domain has it.
        isFreeForUse = true

        // Parse through the domain listing to see if its in use
        for _, domain := range domains {

            // Skip any with no parameters
            if len(domain.Parameters) == 0 {
                continue
            }

            // Loop over the parameters
            for _, parameter := range domain.Parameters {

                // Skip over non-vxlan domains
                if parameter.Name == "type" && parameter.Value != "vxlan"{
                    break
                }

                // Mark it as not free for use so a new number is generated
                if parameter.Name == "multicast_group" && parameter.Value == mcreturn {
                    isFreeForUse = false
                }
            }

        } // end "for _, domain := range domains {"

    } // end "for isFreeForUse  {"

    return mcreturn, nil

}
