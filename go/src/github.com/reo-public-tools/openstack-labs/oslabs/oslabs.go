package main

import (
  "fmt"
  "log"
  "github.com/reo-public-tools/openstack-labs/theforeman"
)

func main() {

    // Set the url to your local forman server
    var url string = "https://172.20.41.28"

    // Get the session id string
    session, err := theforeman.TheForemanLogin(url)
    if err != nil {
        log.Fatal(err)
    }
    //fmt.Printf("Cookie: %s\n", session)


    // Get a foreman domain listing
    domains, err := theforeman.GetDomains(url, session)
    if err != nil {
        log.Fatal(err)
    }
    //fmt.Println(domains)


    // Init curdomaininfo outside of range
    curdomaininfo := theforeman.DomainInfo{}

    // Parse through the listing
    for _, domain := range domains {

        // Ignore the default domain
        if domain.Name == "localdomain" {
            continue
        }

        // Let the user know which domain was found
        fmt.Printf("Found domain: id: %d name: %s fullname: %s\n", domain.ID, domain.Name, domain.Fullname)

        // Get the domain details
        curdomaininfo, err = theforeman.GetDomainDetails(url, session, domain.ID)
        if err != nil {
            log.Fatal(err)
        }

        // init a couple of check vars and set them based off of the domain parameters
        is_free_for_use := 0
        is_type_vlan := 0
        for _, parameter := range curdomaininfo.Parameters {
            if parameter.Name == "in-use" && parameter.Value == "no" {
                is_free_for_use = 1
            }
            if parameter.Name == "type" && parameter.Value == "vlan" {
                is_type_vlan = 1
            }
        }

        // Check if the domain is free fore use
        if is_free_for_use == 1 && is_type_vlan == 1 {
            fmt.Printf("Domain %s is free to check out\n", domain.Name)
            break
        } else {
            fmt.Printf("Domain %s is either the wrong type or in use\n", domain.Name)
            continue
        }

    } // End "for_, domain := range domains {"

    theforeman.CheckOutVLANDomain(url, session, curdomaininfo.ID)
    theforeman.ReleaseVLANDomain(url, session, curdomaininfo.ID)

}
