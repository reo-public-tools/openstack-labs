package theforeman

import (
    "fmt"
    "time"
    "strconv"
    "strings"
    "math/rand"
)

const domainNamePrefix = "lab"
const maxLabs = 50

func CreateDynamicLab(url string, session string) (DomainInfo, error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).CreateDynamicLab(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating a dynamic(vxlan) backed lab.", sysLogPrefix))


    // Get a struct populated with common parameter data
    globalParams, err := GetGlobalParameters(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return DomainInfo{}, err
    }

    // Make call to create the domain in foreman
    domainInfo, err := CreateDynamicLabDomain(url, session, globalParams)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return DomainInfo{}, err
    }

    // Make the call to create the vxlan networks for the new domain
    err = CreateVXLANSubnets(url, session, domainInfo, globalParams)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return DomainInfo{}, err
    }

    // Make the call to create an internal floating ip
    err = AssignInternalFloatingIP(url, session, domainInfo.Name, "MGMT", 10, "")
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return DomainInfo{}, err
    }

    // Get detailed domain info to return
    domainDetails, err := GetDomainDetails(url, session, domainInfo.Name)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return DomainInfo{}, err
    }

    return domainDetails, nil
}

func DeleteDynamicLab(url string, session string, domainName string) (error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).DeleteDynamicLab(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Deleting a dynamic(vxlan) backed lab %s.", sysLogPrefix, domainName))

    // Delete the subnets for this domain
    err := DeleteVXLANSubnets(url, session, domainName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    // Delete the domain
    err = DeleteDomain(url, session, domainName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    return nil
}

func CreateDynamicLabDomain(url string, session string, globalParams GlobalParameters) (Domain, error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).CreateDynamicLabDomain(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating domain for a dynamic(vxlan) backed lab.", sysLogPrefix))

    // Find an available dynamic lab domain slot
    domainName, domainIndex, err := FindAvailableLabSlot(url, session, globalParams.LabBaseDomainName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Domain{}, err
    }

    // Find a free multicast group
    multicastGroup, err := FindAvailableMulticastGroup(url, session, globalParams.MulticastGroupBase)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Domain{}, err
    }

    // Request new external and internal haproxy/keepalived router ids
    // These should all be unique on a flat network.
    externalVRID, internalVRID, err := FindAvailableVRIDs(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Domain{}, err
    }

    // Set the description
    description := fmt.Sprintf("%s%d Dynamic Domain", strings.ToUpper(domainNamePrefix), domainIndex)

    // Convert the location name to id for the domain creation
    locationID, err := ConvLocNameToID(url, session, globalParams.LabLocationName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Domain{}, err
    }

    // Convert the organization name to id for the domain creation
    organizationID, err := ConvOrgNameToID(url, session, globalParams.LabOrgName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Domain{}, err
    }

    // Fill out the data structure to be used for creating the domain
    newDomainStruct := DomainPostData{
        OrganizationID: organizationID,
        LocationID: locationID,
        Domain: NewDomainData {
            Name: domainName,
            Fullname: description,
            DomainParametersAttributes: []DomainParametersAttributes {
                { Name: "type", Value: "vxlan" },
                { Name: "in-use", Value: "yes" },
                { Name: "multicast_group", Value: multicastGroup },
                { Name: "external_vrid", Value: strconv.Itoa(externalVRID) },
                { Name: "internal_vrid", Value: strconv.Itoa(internalVRID) },
            },
        },
    }

    // Make the call to actually create the domain in theForeman
    domainInfo, err := CreateNewDomain(url, session, newDomainStruct)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Domain{}, err
    }

    return domainInfo, nil
}

func FindAvailableLabSlot(url string, session string, baseDomain string) (string, int, error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).FindAvailableLabSlot(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Finding free domain for a dynamic(vxlan) for base domain %s", sysLogPrefix, baseDomain))

    // Get a full foreman domain listing
    domains, err := GetDomains(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
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

    return "", 0, fmt.Errorf("%s We hit the max lab limit without finding an available slot.", sysLogPrefix)

}


func FindAvailableMulticastGroup(url string, session string, multicastGroupBase string) (string, error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).FindAvailableMulticastGroup(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Finding free multicast group using base %s for a new dynamic(vxlan) domain.", sysLogPrefix, multicastGroupBase))

    // Init the return value
    mcreturn := ""

    // Keep track of the used multicast group last oct
    usedMulticastGroups := []string{}

    // Get a full foreman domain listing
    domains, err := GetDomainsWithDetails(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

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
            if parameter.Name == "multicast_group" {
                usedMulticastGroups = append(usedMulticastGroups, parameter.Value)
            }
        }

    } // end "for _, domain := range domains {"

    // Loop through to find a free multicast group oct
    for i := 10; i < 254; i++ {
        contains := false
        matchTest := fmt.Sprintf("%s.%d", multicastGroupBase, i)
        for _, existing := range usedMulticastGroups {
            if existing == matchTest {
                contains = true
            }
        }
        if ! contains {
            mcreturn = fmt.Sprintf("%s.%d", multicastGroupBase, i)
            break
        }
    }

    _ = sysLog.Debug(fmt.Sprintf("%s Found free multicast group %s.", sysLogPrefix, mcreturn))
    return mcreturn, nil

}


// Checks both internal and external used vrids to find a couple that are free for use
// and returns two integers that can be used for either.
func FindAvailableVRIDs(url string, session string) (int, int, error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).FindAvailableVRIDs(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Finding free VRIDs for a new dynamic(vxlan) domain.", sysLogPrefix))

    externalVRID := 0
    internalVRID := 0
    usedVRIDs := []int{}

    // Get a full foreman domain listing
    domains, err := GetDomainsWithDetails(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return 1, 1, err
    }

    // Get a listing of used vrids(both external and internal as they all need to be unique)
    for _, domain := range domains {

        // Skip any with no parameters
        if len(domain.Parameters) == 0 {
            continue
        }

        // Loop over the parameters
        for _, parameter := range domain.Parameters {

            // Mark it as not free for use so a new number is generated
            if strings.Contains(parameter.Name, "ternal_vrid") {
                intValue, err := strconv.Atoi(parameter.Value)
                if err != nil {
                    _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
                    return 1, 1, err
                }
                usedVRIDs = append(usedVRIDs, intValue)
            }
        }

    }

    // Find the external(OSA keepalived config defaults to 10 & 11, so lets start at 12)
    for i := 12; i < 255; i++ {
        contains := false
        for _, existing := range usedVRIDs {
            if existing == i {
                contains = true
            }
        }
        if ! contains {
            externalVRID = i
            usedVRIDs = append(usedVRIDs, i)
            break
        }
    }
    for i := 12; i < 255; i++ {
        contains := false
        for _, existing := range usedVRIDs {
            if existing == i {
                contains = true
            }
        }
        if ! contains {
            internalVRID = i
            break
        }
    }

    _ = sysLog.Debug(fmt.Sprintf("%s Found externalVRID %d and internalVRID %d.", sysLogPrefix, externalVRID, internalVRID))
    return externalVRID, internalVRID, nil

}


func CreateVXLANSubnets(url string, session string, domainInfo Domain, globalParams GlobalParameters) (error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).CreateVXLANSubnets(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating vxlan subnets for domain %s.", sysLogPrefix, domainInfo.Name))

    // Create a random vxlan_id just in case someone external to the labs is using the same group
    rand.Seed(time.Now().UnixNano())
    vxlanID := rand.Intn(16000000 - 50000 + 1) + 50000

    // Convert the location name to id for the domain creation
    locationID, err := ConvLocNameToID(url, session, globalParams.LabLocationName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return err
    }

    // Convert the organization name to id for the domain creation
    organizationID, err := ConvOrgNameToID(url, session, globalParams.LabOrgName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return err
    }

    // Use the same domain prefix for the subnets
    subnetPrefix := strings.ToUpper(strings.Split(domainInfo.Name, ".")[0])

    // Start creating the subnets
    curOffset := 0
    for _, network := range globalParams.VXLANNetworks {
        subnetName := fmt.Sprintf("%s-%s", subnetPrefix, strings.ToUpper(network))
        curNetwork := fmt.Sprintf("%s.%d.0", globalParams.VXLANNetworkPrefix, curOffset)
        curGateway := fmt.Sprintf("%s.%d.1", globalParams.VXLANNetworkPrefix, curOffset)
        curFrom := fmt.Sprintf("%s.%d.50", globalParams.VXLANNetworkPrefix, curOffset)
        curTo := fmt.Sprintf("%s.%d.255", globalParams.VXLANNetworkPrefix, curOffset)
        curOffset += globalParams.VXLANThirdOctetStep

        // Set up the post data struct to create the subnet
        newSubnetStruct := SubnetPostData{
            OrganizationID: organizationID,
            LocationID: locationID,
            Subnet: Subnet {
                Name: subnetName,
                NetworkType: "IPv4",
                Network: curNetwork,
                Mask: globalParams.VXLANNetmask,
                Gateway: curGateway,
                From: curFrom,
                To: curTo,
                Mtu: 9000,
                Ipam: "Random DB",
                BootMode: "Static",
                DomainIds: []int{ domainInfo.ID },
                SubnetParametersAttributes: []SubnetParametersAttributes {
                    { Name: "type", Value: "vxlan" },
                    { Name: "vxlan-id", Value: strconv.Itoa(vxlanID) },
                },
            },
        }

        // Create the subnet
        _, err := CreateSubnet(url, session, newSubnetStruct)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return err
        }

        // Increase the vxlan ID for each network by one
        vxlanID += 1

    }

    return nil
}

func DeleteVXLANSubnets(url string, session string, domainName string) (error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).DeleteVXLANSubnets(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Deleting vxlan subnets for domain %s.", sysLogPrefix, domainName))

    // Get a fresh set of domain details
    domainDetails, err := GetDomainDetails(url, session, domainName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    // Loop over subnets and remove domain associations
    for _, subnet := range domainDetails.Subnets {

        // Remove domain associations
        err := RemoveSubnetFromDomain(url, session, subnet.Name)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return err
        }

        // Delete the subnet
        err = DeleteSubnet(url, session, subnet.Name)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return err
        }

    }

    return nil
}


func AssignInternalFloatingIP(url string, session string, domainName string, netSuffix string, netOffset int, overrideIP string) (error) {

    sysLogPrefix := "theforeman(package).dynamiclabs(file).AssignInternalFloatingIP(func):"

    // init some vars
    internalFloatingIP := ""

    // Get a fresh set of domain details
    domainDetails, err := GetDomainDetails(url, session, domainName)
    if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    if overrideIP != "" {

        _ = sysLog.Debug(fmt.Sprintf("%s Assiging internal floating ip with override %s to domain %s.", sysLogPrefix, overrideIP, domainName))
        internalFloatingIP = overrideIP

    } else {

        _ = sysLog.Debug(fmt.Sprintf("%s Creating internal floating ip for domain %s from subnet suffix %s.", sysLogPrefix, domainName, netSuffix))

        // Set up string to match
        subnetPrefix := strings.ToUpper(strings.Split(domainDetails.Name, ".")[0])
        subnetName := fmt.Sprintf("%s-%s", subnetPrefix, netSuffix)

        // Loop over subnets to find a match
        for _, subnet := range domainDetails.Subnets {
            if subnet.Name == subnetName {
                netPart := strings.Split(subnet.NetworkAddress, "/")[0]
                octList := strings.Split(netPart, ".")
                netStart := strings.Join(octList[0:3], ".")
                netEnd, err := strconv.Atoi(octList[3])
                if err != nil {
                    _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
                    return err
                }
                internalFloatingIP = fmt.Sprintf("%s.%d", netStart, netEnd + netOffset)
                break
            }
        }

    }

    if internalFloatingIP == "" {
        return fmt.Errorf("%s Unable to assign internal floating ip for domain %s\n", sysLogPrefix, domainName)
    }

    _ = sysLog.Debug(fmt.Sprintf("%s Assiging internal floating ip %s to domain %s.", sysLogPrefix, internalFloatingIP, domainName))
    err = SetDomainParameter(url, session, domainDetails.ID, "internal_floating_ip", internalFloatingIP)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    return nil

}
