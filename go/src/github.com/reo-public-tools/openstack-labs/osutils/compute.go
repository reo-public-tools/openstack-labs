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

    // Conver the flavor name to id
    flavorID, err := flavors.IDFromName(computeClient, flavorName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    // Get extra spec for capability value of a given key
    flavorExtraSpecs := flavors.GetExtraSpec(computeClient, flavorID, "capabilities:" + key)

    // This is weird
    return flavorExtraSpecs.Body.(map[string]interface{})["capabilities:" + key].(string), nil
}
