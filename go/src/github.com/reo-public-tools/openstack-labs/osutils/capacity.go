package osutils

import (
        "fmt"
        "github.com/gophercloud/gophercloud"
)

// structs used to hold capacity info for user
type IronicCapacityData struct {
        NodesByType    []IronicNodeCapInfo
        CountByProject []IronicProjectCount
        CountByUser    []IronicUserCount
}

type IronicNodeCapInfo struct {
        CapacityType string
        Flavor       string
        Used         int
        Free         int
}

type IronicProjectCount struct {
        Project string
        Total   int
}

type IronicUserCount struct {
        User    string
        Total   int
}



// Methods
func (icd *IronicCapacityData) UpdateIronicNodeCapInfo(capType string, flavor string, status string) {

    // Check if it exists so we can init if needed
    foundItem := false
    for _, curCapInfo := range icd.NodesByType {

        // Update if it exists
        if curCapInfo.CapacityType == capType {
            foundItem = true
        }
    }

    // Append a new entry if needed
    if ! foundItem {

        // Init the struct to append
        appendData := IronicNodeCapInfo{
                CapacityType: capType,
                Flavor:  flavor,
                Used: 0,
                Free: 0,
        }

        // Append it to the array
        icd.NodesByType = append(icd.NodesByType, appendData)

    }

    // Increment values based on status
    for index, curCapInfo := range icd.NodesByType {
        // Update if it exists
        if curCapInfo.CapacityType == capType {
            if status == "free" {
                icd.NodesByType[index].Free++
            } else {
                icd.NodesByType[index].Used++
            }
        }
    }

}


// Pull and format the data for display use
func GetIronicCapacity(provider *gophercloud.ProviderClient) (IronicCapacityData, error) {

    sysLogPrefix := "osutils(package).capacity(file).GetIronicCapacity(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting ironic capacity information for user display.", sysLogPrefix))


    // init the return data stor
    ironicCapacityDataRet := IronicCapacityData{}

    // Get a detailed ironic node listing
    curNodeList, err := GetDetailedIronicNodeList(provider)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return IronicCapacityData{}, err
    }

    // Init a local system_type to flavor name mapping to cut down on unneeded lookups
    stofmap, err := GetFlavorToSystemTypeMappings(provider)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return IronicCapacityData{}, err
    }

    // Loop through the nodes and update return data struct
    for _, curnode := range curNodeList {

        // Break out the capabilities for this node
        capabilities := mapCapabilityString(curnode.Properties["capabilities"].(string))
        systemType := capabilities["system_type"]

        // Get the flavor
        flavor := ""
        if _, ok := stofmap[systemType]; ok {
            flavor = stofmap[systemType]
        }

        // Get capability=>flavor mappings
        status := "free"
        if curnode.ProvisionState != "" || curnode.InstanceUUID != "" || curnode.Maintenance == true {
            status = "used"
        }
        ironicCapacityDataRet.UpdateIronicNodeCapInfo(
                                             capabilities["system_type"],
                                             flavor,
                                             status)

    }

    return ironicCapacityDataRet, nil

}
