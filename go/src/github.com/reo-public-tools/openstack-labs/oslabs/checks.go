package main

import (
    "strings"
    "github.com/reo-public-tools/openstack-labs/theforeman"
    "github.com/reo-public-tools/openstack-labs/osutils"
)

func CheckIronicCapacity(labConfig LabConfig, session string) (bool, error) {

    // Init the ironic request to be used for the capacity check
    IOIRequest := osutils.IronicOnIronicRequest{}

    // Loop through the config and get the flavor name for each foremanhost request part
    for _, node := range labConfig.Environment.ForemanHosts {

        // Get the compute profile from the hostgroup.
        hostgroupInfo, err := theforeman.GetHostgroupInfo(labConfig.ForemanURL, session, node.ForemanHostgroup)
        if err != nil {
            return false, err
        }

        if strings.Contains(hostgroupInfo.ComputeProfileName, "ironic") {

            flavorName, err := theforeman.ConvProfileNameToFlavorName(labConfig.ForemanURL, session, hostgroupInfo.ComputeProfileName)
            if err != nil {
                return false, err
            }

            curIOINode := osutils.IronicOnIronicNodeRequest{
                Flavor: flavorName,
                Count: node.Count,
            }

            IOIRequest.IronicOnIronicNodeRequests = append(IOIRequest.IronicOnIronicNodeRequests, curIOINode)

        }

    }

    // Loop through the config and update/append to the request for ironic nodes.
    for _, node := range labConfig.Environment.IronicOnIronicHosts {

        existing := false

        // Update the existing flavor if more are requeted through ioi nodes.
        for ioiindex, ioirnode := range IOIRequest.IronicOnIronicNodeRequests {
            if ioirnode.Flavor == node.Flavor {
                IOIRequest.IronicOnIronicNodeRequests[ioiindex].Count += node.Count
                existing = true
            }
        }

        // Append to list if the flavor is not already there
        if ! existing {
            curIOINode := osutils.IronicOnIronicNodeRequest{
                Flavor: node.Flavor,
                Count: node.Count,
            }

            IOIRequest.IronicOnIronicNodeRequests = append(IOIRequest.IronicOnIronicNodeRequests, curIOINode)
        }

    }

    // Connect to openstack
    provider, err := osutils.OpenstackLogin(labConfig.OpenstackCloud)
    if err != nil {
        return false, err
    }

    // Check the ironic capacity
    check, err := osutils.CheckIronicCapacity(&provider, IOIRequest)
    if err != nil {
        return false, err
    }

    return check, nil
}
