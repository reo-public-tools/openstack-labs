package theforeman

import (
    "fmt"
)

func CheckOutVLANDomain(url string, session string, domainid int) (error) {

    sysLogPrefix := "theforeman(package).staticlabs(file).CheckOutVLANDomain(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Checking out existing static VLAN domain %d", sysLogPrefix, domainid))

    err := SetDomainParameter(url, session, domainid, "in-use", "yes")
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }
    return nil
}

func ReleaseVLANDomain(url string, session string, domainid int) (error) {

    sysLogPrefix := "theforeman(package).staticlabs(file).ReleaseVLANDomain(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Releasing existing static VLAN domain %d", sysLogPrefix, domainid))

    err := SetDomainParameter(url, session, domainid, "in-use", "no")
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }
    return nil
}

func FindAvailableVLANDomain(url string, session string) (DomainInfo, error) {

    sysLogPrefix := "theforeman(package).staticlabs(file).FindAvailableVLANDomain(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Finding an existing static VLAN domain.", sysLogPrefix))

    // Get a full foreman domain listing
    domains, err := GetDomains(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return DomainInfo{}, err
    }

    // Init curdomaininfo outside of range
    curdomaininfo := DomainInfo{}

    // Parse through the listing
    for _, domain := range domains {

        // Ignore the default domain
        if domain.Name == "localdomain" {
            continue
        }

        // Get the domain details
        curdomaininfo, err = GetDomainDetails(url, session, domain.ID)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return DomainInfo{}, err
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
            _ = sysLog.Debug(fmt.Sprintf("%s Found static VLAN domain %s that is not in use.", sysLogPrefix, curdomaininfo.Name))
            return curdomaininfo, nil
        }

    } // End "for_, domain := range domains {"


    return DomainInfo{}, fmt.Errorf("%s No static vlan backed domains available for use.", sysLogPrefix)
}

