package osutils

import (
    "fmt"
    "github.com/gophercloud/gophercloud"
    "github.com/gophercloud/gophercloud/openstack"
    "github.com/gophercloud/utils/openstack/clientconfig"
)

func OpenstackLogin(cloud string) (gophercloud.ProviderClient, error) {

    sysLogPrefix := "osutils(package).auth(file).OpenstackLogin(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Setting up and testing ProviderClient using clouds.yml based config.", sysLogPrefix))


    // Set up the client opts from your clouds.yml file
    clientOpts := &clientconfig.ClientOpts{
        Cloud: cloud,
    }

    // Get an authenticated client 
    provider, err := clientconfig.AuthenticatedClient(clientOpts)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return gophercloud.ProviderClient{}, err
    }

    // Do a quick identity test
    _, err = openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{ Region: "RegionOne", })
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return gophercloud.ProviderClient{}, err
    }

    return *provider, nil
}
