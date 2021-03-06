package main

import (
    "fmt"
    "sync"
//    "time"
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


    /* #### Ironic-On-Ironic node allocation #### */

    // Init the request struct
    IOIRequest := osutils.IronicOnIronicRequest{}

    // Loop through the config and update/append to the request for ironic nodes.
    for _, node := range labConfig.Environment.IronicOnIronicHosts {

        // Append to list
        curIOINode := osutils.IronicOnIronicNodeRequest{
            Flavor: node.Flavor,
            Count: node.Count,
        }

        IOIRequest.IronicOnIronicNodeRequests = append(IOIRequest.IronicOnIronicNodeRequests, curIOINode)

    }


    if IOIRequest.IronicOnIronicNodeRequests != nil {

        // Check out any ironic on ironic nodes
        _ = sysLog.Info(fmt.Sprintf("Reserving out ironic-on-ironic nodes\n"))
        fmt.Printf("Reserving out ironic-on-ironic nodes\n")
        IOINodeList, err := osutils.CheckOutIronicNodes(&provider, IOIRequest, curdomaininfo.Name)
        if err != nil {
            return err
        }

        // Marshal the node data struct to a base64 encoded json string
        _ = sysLog.Info(fmt.Sprintf("Converting request to base64 encoded json\n"))
        fmt.Printf("Converting request to base64 encoded json\n")
        jsonString, err := osutils.NodeDataToJSONString(IOINodeList, true)
        if err != nil {
            return err
        }

        // Set the domain parameter to the encoded value
        _ = sysLog.Info(fmt.Sprintf("Updating ioi_data domain parameter\n"))
        fmt.Printf("Updating ioi_data domain parameter\n")
        err = theforeman.SetDomainParameter(labConfig.ForemanURL, session, curdomaininfo.Name, "ioi_data", jsonString)
        if err != nil {
            return err
        }

    }


    /* #### Add ssh keys for the admin user as host level parameters are not seen by cloud-init #### */
    _ = sysLog.Info(fmt.Sprintf("Updating ssh_authorized_keys domain parameter\n"))
    fmt.Printf("Updating ssh_authorized_keys domain parameter\n")
    err = theforeman.SetDomainParameter(labConfig.ForemanURL, session, curdomaininfo.Name, "ssh_authorized_keys", labConfig.PubSSHKey)
    if err != nil {
        return err
    }


    /* #### Start foreman host buildout #### */

    // Setting up a type to use for the goroutine return
    type hostCreateRet struct {
        hostInfo theforeman.Host
        err      error
    }

    // Set up the channel for communications
    hostInfoChan := make(chan hostCreateRet)

    // Track the number of working goroutines
    var wg sync.WaitGroup

    // Loop through the foreman hosts in the lab config and create
    for _, foremanHost := range labConfig.Environment.ForemanHosts {

        // Name the host based on role and increment it.
        for i := 1; i <= foremanHost.Count; i++ {

            // Set up the hostname
            hostName := fmt.Sprintf("%s%0.2d", foremanHost.Role, i)

            // Add to the wait group
            wg.Add(1)

            go func(url string, session string, labName string, hostName string, hostGroup string, role string) {

                // Defer wait group finish
                defer wg.Done()

                // Init the return
                var hostRet hostCreateRet

                // Create the host and send the results back to the main goroutine
                fmt.Printf("Creating new node %s.%s in host group %s\n", hostName, labName, hostGroup)
                hostRet.hostInfo, hostRet.err = theforeman.CreateHost(url, session, labName, hostName, hostGroup, role)
                hostInfoChan <- hostRet

            } (labConfig.ForemanURL, session, curdomaininfo.Name, hostName, foremanHost.ForemanHostgroup, foremanHost.Role)

            // Sleep for a few seconds to give dhcp time to grab an ip.
            //fmt.Printf("Sleeping for %d seconds to give dhcp enough time to pull sequential ip addresses\n", 60)
            //time.Sleep(60 * time.Second)

        }
    }

    // Run the closer goroutine
    go func() {
        wg.Wait()
        close(hostInfoChan)
    }()

    // Loop over results
    for hostRetInfo := range hostInfoChan {
        if hostRetInfo.err != nil {
            fmt.Println(hostRetInfo.err)
        } else {
            fmt.Printf("Finished creating host %s with primay ip of %s.\n", hostRetInfo.hostInfo.Name, hostRetInfo.hostInfo.IP)
        }
    }



    /* #### Wrap things up and display the results to the user #### */

    // Print out the updated results
    err = Show(curdomaininfo.Name)
    if err != nil {
        return err
    }

    return nil

}



