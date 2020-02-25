package osutils

import (
        "fmt"
        "github.com/gophercloud/gophercloud"
        "github.com/gophercloud/gophercloud/openstack"
        "github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)


func GetFlavorCapability(provider *gophercloud.ProviderClient, flavorName string, key string) (string, error) {

    sysLogPrefix := "osutils(package).compute(file).GetFlavorCapabilities(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting capabilities from flavorName %s.", sysLogPrefix, flavorName))

    // Get a baremetal serviceclient
    computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    // Convert the flavor name to id
    flavorID, err := flavors.IDFromName(computeClient, flavorName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    // Get extra spec for capability value of a given key
    flavorExtraSpecs := flavors.GetExtraSpec(computeClient, flavorID, "capabilities:" + key)

    if flavorExtraSpecs.Body == nil {
        _ = sysLog.Err(fmt.Sprintf("%s capability %s not found for flavor %s", sysLogPrefix, key, flavorName))
        return "notfound", err
    } else {
        // This is weird
        return flavorExtraSpecs.Body.(map[string]interface{})["capabilities:" + key].(string), nil
    }

}

func GetFlavorToSystemTypeMappings(provider *gophercloud.ProviderClient) (map[string]string, error) {

    sysLogPrefix := "osutils(package).compute(file).GetFlavorToSystemTypeMappings(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting flavor name to system type mappings for type %s.", sysLogPrefix))

    stofmap := make(map[string]string)

    // Get a baremetal serviceclient
    computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return stofmap, err
    }

    // Loop through all of the flavors until we find a match
    listOpts := flavors.ListOpts{ AccessType: flavors.PublicAccess, }

    allPages, err := flavors.ListDetail(computeClient, listOpts).AllPages()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return stofmap, err
    }

    allFlavors, err := flavors.ExtractFlavors(allPages)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return stofmap, err
    }

    // Init a local system_type to flavor name mapping to cut down on unneeded lookups

    for _, flavor := range allFlavors {

        curSystemType, err := GetFlavorCapability(provider, flavor.Name, "system_type") 
        if err != nil && curSystemType != "notfound" {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return stofmap, err
        }
        if curSystemType == "notfound" {
            continue
        }

        stofmap[curSystemType] = flavor.Name
    }

    return stofmap, nil

}
