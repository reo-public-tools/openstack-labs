package osutils

import (
        "fmt"
        "github.com/gophercloud/gophercloud"
        "github.com/gophercloud/gophercloud/openstack"
        "github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
//        "github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
        utilnetworks "github.com/gophercloud/utils/openstack/networking/v2/networks"
)


func NetworkNameToID(provider *gophercloud.ProviderClient, networkName string) (string, error ) {

    sysLogPrefix := "osutils(package).networking(file).NetworkNameToID(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting network id from name %s.", sysLogPrefix, networkName))

    // Get the network service client
    networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    // Use build-in function to get the id from the name
    netID, err := utilnetworks.IDFromName(networkClient, networkName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return "", err
    }

    return netID, nil
}

func CreateNetworkPort(provider *gophercloud.ProviderClient, networkName string, labName string) (ports.Port, error) {

    sysLogPrefix := "osutils(package).networking(file).CreateNetworkPort(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating a new network port for network %s lab: %s.", sysLogPrefix, networkName, labName))


    // Pull the project ID
    projectID, err := GetMyProjectID(provider)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return ports.Port{}, err
    }

    // Pull the network ID from the provided name
    netID, err := NetworkNameToID(provider, networkName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return ports.Port{}, err
    }

    // Get the network service client
    networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return ports.Port{}, err
    }

    // See if this port exists and return the info.
    listOpts := ports.ListOpts {
        Name: fmt.Sprintf("%s_external_floating_ip", labName),
    }
    allPages, err := ports.List(networkClient, listOpts).AllPages()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return ports.Port{}, err
    }

    isEmpty, err := allPages.IsEmpty()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return ports.Port{}, err
    }

    if ! isEmpty {
        allPorts, err := ports.ExtractPorts(allPages)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return ports.Port{}, err
        }
        return allPorts[0], nil
    }

    // Create a port if we have made it this far
    createOpts := ports.CreateOpts{
            Name: fmt.Sprintf("%s_external_floating_ip", labName),
            Description: fmt.Sprintf("%s_external_floating_ip", labName),
            NetworkID: netID,
            ProjectID: projectID,
    }

    port, err := ports.Create(networkClient, createOpts).Extract()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return ports.Port{}, err
    }

    return *port, nil
}


func DeleteNetworkPort(provider *gophercloud.ProviderClient, portID string) (error) {

    sysLogPrefix := "osutils(package).networking(file).DeleteNetworkPort(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Delete network port with port id: %s", sysLogPrefix, portID))


    // Get the network service client
    networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    // Delete the port
    err = ports.Delete(networkClient, portID).ExtractErr()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    return nil
}


func DeleteExtNetworkPortByLabName(provider *gophercloud.ProviderClient, labName string) (error) {

    sysLogPrefix := "osutils(package).networking(file).DeleteExtNetworkPortByLabName(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Delete network port for lab: %s", sysLogPrefix, labName))


    // Get the network service client
    networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    // See if this port exists and return the info.
    listOpts := ports.ListOpts {
        Name: fmt.Sprintf("%s_external_floating_ip", labName),
    }
    allPages, err := ports.List(networkClient, listOpts).AllPages()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    isEmpty, err := allPages.IsEmpty()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    if isEmpty {
        return nil
    }

    allPorts, err := ports.ExtractPorts(allPages)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    // Delete the port
    err = ports.Delete(networkClient, allPorts[0].ID).ExtractErr()
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    return nil
}
