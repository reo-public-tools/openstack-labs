package osutils

import (
        "fmt"
        "strings"
        "encoding/json"
        "encoding/base64"
        "github.com/gophercloud/gophercloud"
        "github.com/gophercloud/gophercloud/openstack"
        "github.com/gophercloud/gophercloud/pagination"
        "github.com/gophercloud/gophercloud/openstack/baremetal/v1/nodes"
        "github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
)

// Structures for the ironic-on-ironic request
type IronicOnIronicRequest struct {
        IronicOnIronicNodeRequests []IronicOnIronicNodeRequest
}

type IronicOnIronicNodeRequest struct {
        Flavor string
        Count  int
}

// Structures for checked-out node data converted to json for theforeman domain parameter
type IronicNodeDetailsList struct {
        IronicNodeDetails []IronicNodeDetails
}

type IronicNodeDetails struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	CPUArch      string   `json:"cpu_arch"`
	Cpus         float64  `json:"cpus"`
	MemoryMb     float64  `json:"memory_mb"`
	LocalGb      float64  `json:"local_gb"`
	Size         float64  `json:"size"`
	IpmiAddress  string   `json:"impi_address"`
	IpmiPassword string   `json:"impi_password"`
	IpmiUsername string   `json:"impi_username"`
	Macs         []string `json:"macs"`
        Flavor       string   `json:"flavor"`
        SystemType   string   `json:"system_type"`
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


func CheckOutIronicNodes(provider *gophercloud.ProviderClient, request IronicOnIronicRequest, reason string) (IronicNodeDetailsList, error) {

    sysLogPrefix := "osutils(package).ironic(file).CheckOutIronicNodes(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Checking out ironic-on-ironic nodes.", sysLogPrefix))

    retList := IronicNodeDetailsList{}

    // Get a baremetal serviceclient
    baremetalClient, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return IronicNodeDetailsList{}, err
    }


    for _, nodeRequest := range request.IronicOnIronicNodeRequests {

        // Get capability for system_type from provided flavor.
        flavorCapability, err := GetFlavorCapability(provider, nodeRequest.Flavor, "system_type")
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return IronicNodeDetailsList{}, err
        }

        // Pull a list of available nodes by capability
        curNodeList, err := GetAvailableIronicNodeListByCapability(provider, "system_type", flavorCapability)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return IronicNodeDetailsList{}, err
        }

        // Check the availability of the current flavor
        if nodeRequest.Count <= len(curNodeList) {
            for _, curNode := range curNodeList[:nodeRequest.Count] {

                _ = sysLog.Debug(fmt.Sprintf("%s Checking out node %s(%s)", sysLogPrefix, curNode.UUID, curNode.Name))

                // Move node into maintenance mode here.
                updateOpts := nodes.UpdateOpts{
                        nodes.UpdateOperation{
                                Op:    nodes.ReplaceOp,
                                Path:  "/maintenance",
                                Value: "true",
                        },
                }

                // gophercloud doesn't allow for setting /maintenance_reason. This would be nice
                // to have, but not needed.

                _, err := nodes.Update(baremetalClient, curNode.UUID , updateOpts).Extract()
                if err != nil {
                    _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
                    return IronicNodeDetailsList{}, err
                }


                // Get mac list for current node
                macList, err := GetIronicPXEMacs(provider, curNode.UUID)
                if err != nil {
                    _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
                    return IronicNodeDetailsList{}, err
                }


                // Create detail data and append it to the return list
                curDetails := IronicNodeDetails{
                        ID: curNode.UUID,
                        Name: curNode.Name,
                        CPUArch: curNode.Properties["cpu_arch"].(string),
                        Cpus: curNode.Properties["cpus"].(float64),
                        MemoryMb: curNode.Properties["memory_mb"].(float64),
                        LocalGb: curNode.Properties["local_gb"].(float64),
                        Size: curNode.Properties["size"].(float64),
                        IpmiUsername: curNode.DriverInfo["ipmi_username"].(string),
                        IpmiPassword: curNode.DriverInfo["ipmi_password"].(string),
                        IpmiAddress: curNode.DriverInfo["ipmi_address"].(string),
                        Macs: macList,
                        Flavor: nodeRequest.Flavor,
                        SystemType: flavorCapability,
                    }
                retList.IronicNodeDetails = append(retList.IronicNodeDetails, curDetails)

            } // End "for _, curNode := range curNodeList[:nodeRequest.Count] {"
        } // End "if nodeRequest.Count <= len(curNodeList) {"
    } // End "for _, nodeRequest := range request.IronicOnIronicNodeRequests {"

    return retList, nil
}

func ReleaseIronicNodes(provider *gophercloud.ProviderClient, nodeList IronicNodeDetailsList, cycleClean bool) (error) {

    sysLogPrefix := "osutils(package).ironic(file).ReleaseIronicNodes(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Releasing ironic-on-ironic nodes.", sysLogPrefix))

    // Get a baremetal serviceclient
    baremetalClient, err := openstack.NewBareMetalV1(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }
    baremetalClient.Microversion = "1.22"

    for _, curNode := range nodeList.IronicNodeDetails {

        _ = sysLog.Debug(fmt.Sprintf("%s Releasing node %s(%s)", sysLogPrefix, curNode.ID, curNode.Name))

        // Set node provision state to manage if cleaning
        if cycleClean {

            provStateOpts := nodes.ProvisionStateOpts{
                    Target: nodes.TargetManage,
            }

            err := nodes.ChangeProvisionState(baremetalClient, curNode.ID, provStateOpts).ExtractErr()
            if err != nil {
                _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
                return err
            }
        }

        // Move node into maintenance mode here.
        updateOpts := nodes.UpdateOpts{
                nodes.UpdateOperation{
                        Op:    nodes.ReplaceOp,
                        Path:  "/maintenance",
                        Value: "false",
                },
        }

        _, err := nodes.Update(baremetalClient, curNode.ID, updateOpts).Extract()
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return err
        }

        // Set node provision state to provide to kick off cleaning
        if cycleClean {

            provStateOpts := nodes.ProvisionStateOpts{
                    Target: nodes.TargetProvide,
            }

            err = nodes.ChangeProvisionState(baremetalClient, curNode.ID, provStateOpts).ExtractErr()
            if err != nil {
                _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
                return err
            }
        }


    } // End "for _, curNode := range request.IronicOnIronicNodeRequests {"

    return nil
}

func NodeDataToJSONString(nodeData IronicNodeDetailsList, b64encode bool) (string, error) {

    sysLogPrefix := "osutils(package).ironic(file).NodeDataToJSONString(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Converting ironic node data to json string.", sysLogPrefix))

    jsonRet, err := json.Marshal(nodeData)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    if b64encode {
        jsonRet = []byte(base64.StdEncoding.EncodeToString(jsonRet))
    }

    return string(jsonRet), err

}

func JSONStringToNodeData(nodeJSON string, b64decode bool) (IronicNodeDetailsList, error) {

    sysLogPrefix := "osutils(package).ironic(file).JSONStringToNodeData(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Converting json string into ironic node data", sysLogPrefix))


    var retData IronicNodeDetailsList
    bytes := []byte(nodeJSON)

    var err error
    if b64decode {
        bytes, err = base64.StdEncoding.DecodeString(string(bytes))
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return IronicNodeDetailsList{}, err
        }
    }

    err = json.Unmarshal(bytes, &retData)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return IronicNodeDetailsList{}, err
    }

    return retData, err

}
