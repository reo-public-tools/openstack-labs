package main

import (
    "fmt"
    "github.com/reo-public-tools/openstack-labs/theforeman"
    "github.com/reo-public-tools/openstack-labs/osutils"
)

func Create(configFile string) (error) {

    // Populate struct from config file
    _ = sysLog.Info(fmt.Sprintf("Pulling in config: %s\n", configFile))
    fmt.Printf("Pulling in config: %s\n", configFile)
    labConfig, err := PopulateConfigFile(configFile)
    if err != nil {
        return err
    }

    // Connect to theforeman
    _ = sysLog.Info(fmt.Sprintf("Logging into the foreman at: %s\n", labConfig.ForemanURL))
    fmt.Printf("Logging into the foreman at: %s\n", labConfig.ForemanURL)
    session, err := theforeman.TheForemanLogin(labConfig.ForemanURL)
    if err != nil {
        return err
    }

    // Check the ironic capacity for the request
    _ = sysLog.Info(fmt.Sprintf("Checking ironic capacity\n"))
    fmt.Printf("Checking ironic capacity\n")
    capCheck, err := CheckIronicCapacity(labConfig, session)
    if err != nil || ! capCheck {
        return err
    }
    _ = sysLog.Info(fmt.Sprintf("Ironic capacity is good.\n"))
    fmt.Printf("Ironic capacity is good.\n")


    /* #### Set up the domain in theforeman static | dynamic #### */

    // Track the domaininfo
    var curdomaininfo theforeman.DomainInfo

    // Create a dynamic(vxlan based) domain or check out a vlan domain
    if labConfig.LabType == "static" {

        // Find an availalbe vlan backed domain.
        _ = sysLog.Info(fmt.Sprintf("Looking for a free static domain\n"))
        fmt.Printf("Looking for a free static domain.\n")
        curdomaininfo, err = theforeman.FindAvailableVLANDomain(labConfig.ForemanURL, session)
        if err != nil {
            return err
        }
        _ = sysLog.Info(fmt.Sprintf("Found available domain %s to check out\n", curdomaininfo.Name))
        fmt.Printf("Found available domain %s to check out\n", curdomaininfo.Name)

        // Check out the static vlan backed domain
        _ = sysLog.Info(fmt.Sprintf("Checking out static lab %s\n", curdomaininfo.Name))
        fmt.Printf("Checking out lab %s\n", curdomaininfo.Name)
        err = theforeman.CheckOutVLANDomain(labConfig.ForemanURL, session, curdomaininfo.ID)
        if err != nil {
            return err
        }

    } else {

        _ = sysLog.Info("Creating a new dynamic lab.")
        fmt.Printf("Creating a new dynamic lab\n")
        curdomaininfo, err = theforeman.CreateDynamicLab(labConfig.ForemanURL, session)
        if err != nil {
            return err
        }

       // Print out domain details here(domain show)
       fmt.Println(curdomaininfo)

    }


    /* #### Set up the external_floating_ip #### */

    // Connect to openstack
    provider, err := osutils.OpenstackLogin(labConfig.OpenstackCloud)
    if err != nil {
        return err
    }

    // Create a neutron port and set the external floating ip on the new domain from the port.
    _ = sysLog.Info(fmt.Sprintf("Creating the neutron port for the external_floating_ip parameter\n"))
    fmt.Printf("Creating the neutron port for the external_floating_ip parameter\n")
    portInfo, err := osutils.CreateNetworkPort(&provider, labConfig.OpenstackNetwork, curdomaininfo.Name)
    if err != nil {
        return err
    }

    // Set the domain parameter to the encoded value
    _ = sysLog.Info(fmt.Sprintf("Setting the external_floating_ip parameter to the new port ip\n"))
    fmt.Printf("Setting the external_floating_ip parameter to the new port ip\n")
    err = theforeman.SetDomainParameter(labConfig.ForemanURL, session, curdomaininfo.Name, "external_floating_ip", portInfo.FixedIPs[0].IPAddress)
    if err != nil {
        return err
    }


    return nil

}



