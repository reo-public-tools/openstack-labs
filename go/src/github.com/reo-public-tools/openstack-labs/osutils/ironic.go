package osutils

import (
        "fmt"
        "strings"
        "github.com/gophercloud/gophercloud"
        "github.com/gophercloud/gophercloud/openstack"
        "github.com/gophercloud/gophercloud/pagination"
        "github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
        "github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)


type IronicOnIronicRequest struct {
        IronicOnIronicNodeRequests []IronicOnIronicNodeRequest
}

type IronicOnIronicNodeRequest struct {
        Flavor string
        Count  int
}

func mapCapabilityString(capabilities string) (map[string]string) {

    // Initialize the map
    retMap := make(map[string]string)

    // Split key:val pairs by ','
    capabilityList := strings.Split(capabilities, ",")

    // Loop through and split each key:val pair by ':' then assign to map
    for _, curCapability := range capabilityList {
        curKeyVal := strings.Split(curCapability, ":")
        retMap[curKeyVal[0]] = curKeyVal[1]
    }

    return retMap

}

func GetDetailedIronicNodeList(provider *gophercloud.ProviderClient) ([]nodes.Node, error) {

    sysLogPrefix := "osutils(package).ironic(file).GetDetailedIronicNodeList(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting detailed ironic node list. ", sysLogPrefix))

    retNodeList := []nodes.Node{}

    // Get a baremetal serviceclient
    baremetalClient, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []nodes.Node{}, err
    }

    // List nodes with detail and append only those that match to the return list
    pages := 0
    err = nodes.ListDetail(baremetalClient, nodes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
        pages += 1

        results, err := nodes.ExtractNodes(page)
        if err != nil {
            return false, err
        }

        retNodeList = append(retNodeList, results...)

        return true, nil
    })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []nodes.Node{}, err
    }

    return retNodeList, nil
}


func GetIronicNodeListByCapability(provider *gophercloud.ProviderClient, ckey string, cval string) ([]nodes.Node, error) {

    sysLogPrefix := "osutils(package).ironic(file).GetIronicNodeListByCapability(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting ironic node list with capability %s=%s", sysLogPrefix, ckey, cval))

    retNodeList := []nodes.Node{}

    curNodeList, err := GetDetailedIronicNodeList(provider)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []nodes.Node{}, err
    }

    for _, curnode := range curNodeList {
        capabilities := mapCapabilityString(curnode.Properties["capabilities"].(string))
        if capabilities[ckey] == cval {
            retNodeList = append(retNodeList, curnode)
        }
    }

    return retNodeList, nil
}

func GetAvailableIronicNodeListByCapability(provider *gophercloud.ProviderClient, ckey string, cval string) ([]nodes.Node, error) {

    sysLogPrefix := "osutils(package).ironic(file).GetAvailableIronicNodeListByCapability(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting available ironic node list with capability %s=%s", sysLogPrefix, ckey, cval))

    retNodeList := []nodes.Node{}

    curNodeList, err := GetIronicNodeListByCapability(provider, ckey, cval)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []nodes.Node{}, err
    }

    for _, curnode := range curNodeList {

        // Filter out nodes in maintenance mode
        if curnode.Maintenance {
            continue
        }

        // Filter out only those available(var looks to be empty if available)
        if curnode.ProvisionState != "" {
            continue
        }

        // Filter out any with a defined instance uuid
        if curnode.InstanceUUID != "" {
            continue
        }

        // Append if we made it this far
        retNodeList = append(retNodeList, curnode)
    }

    return retNodeList, nil
}

func CheckIronicCapacity(provider *gophercloud.ProviderClient, request IronicOnIronicRequest) (bool, error) {

    sysLogPrefix := "osutils(package).ironic(file).CheckIronicCapacity(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Checking ironic capacity for request", sysLogPrefix))

    /* Each IronicOnIronicNodeRequest has the flavor and count.  The
       flavor will need converted to the capability then counts checked
       for each */

    for _, nodeRequest := range request.IronicOnIronicNodeRequests {

        // Get capability for system_type from provided flavor.
        flavorCapability, err := GetFlavorCapability(provider, nodeRequest.Flavor, "system_type")
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return false, err
        }

        // Pull a list of available nodes by capability
        curNodeList, err := GetAvailableIronicNodeListByCapability(provider, "system_type", flavorCapability)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return false, err
        }

        // Check the availability of the current flavor
        if nodeRequest.Count > len(curNodeList) {
            curError := fmt.Errorf("%s Requesting %d of type %s, but %d available.\n", sysLogPrefix, nodeRequest.Count, nodeRequest.Flavor, len(curNodeList))
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, curError))
            return false, curError
        }

    }

    return true, nil
}

func GetDetailedIronicPortList(provider *gophercloud.ProviderClient) ([]ports.Port, error) {

    sysLogPrefix := "osutils(package).ironic(file).GetDetailedIronicPortList(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting detailed ironic port listing.", sysLogPrefix))


    retPortList := []ports.Port{}

    // Get a baremetal serviceclient
    baremetalClient, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []ports.Port{}, err
    }

    // List ports with detail and append only those that match to the return list
    pages := 0
    err = ports.ListDetail(baremetalClient, ports.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
        pages += 1

        results, err := ports.ExtractPorts(page)
        if err != nil {
            return false, err
        }

        retPortList = append(retPortList, results...)

        return true, nil
    })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []ports.Port{}, err
    }

    return retPortList, nil

}

func GetIronicPXEMacs(provider *gophercloud.ProviderClient, nodeID string) ([]string, error) {

    sysLogPrefix := "osutils(package).ironic(file).GetIronicPXEMacs(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Get a port mac list for a ironic node id %s.", sysLogPrefix, nodeID))

    macList := []string{}

    portArray, err := GetDetailedIronicPortList(provider)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []string{}, err
    }
    for _, curPort := range portArray {
        if curPort.NodeUUID == nodeID {
            macList = append(macList, curPort.Address)
        }
    }

    return macList, nil
}
