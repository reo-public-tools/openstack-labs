package main

import (
    "fmt"
    "strings"
    "github.com/reo-public-tools/openstack-labs/theforeman"
)

func Scratch(configFile string) (error) {

    labName := "lab2.phobos.rpc.rackspace.com"


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

    // Get domain details from provided name
    _ = sysLog.Info(fmt.Sprintf("Pulling info for domain: %s\n", labName))
    fmt.Printf("Pulling info for domain: %s.\n", labName)
    curdomaininfo, err := theforeman.GetDomainDetails(labConfig.ForemanURL, session, labName)
    if err != nil {
        return err
    }

    // Get host listing for the domain
    _ = sysLog.Info(fmt.Sprintf("Pulling host listing for domain for deletion: %s\n", curdomaininfo.Name))
    fmt.Printf("Pulling host listing for domain: %s.\n", curdomaininfo.Name)
    detailedHostList, err := theforeman.GetHostsDetailsByDomainID(labConfig.ForemanURL, session, curdomaininfo.ID)
    if err != nil {
        return err
    }

    // Set up some subnet name -> id mappings 
    subnetmap := make(map[string]int)
    for _, curSubnet := range curdomaininfo.Subnets {
        subnetmap[curSubnet.Name] = curSubnet.ID
    }

    // Pull the lab name from the fqdn of the lab name
    labShortName := strings.ToUpper(strings.Split(labName, ".")[0])

    // Loop over hosts and add interfaces to each
    for _, curHost := range detailedHostList {

        _ = sysLog.Info(fmt.Sprintf("Adding interface to host: %s\n", curHost.Name))
        fmt.Printf("Deleting host: %s.\n", curdomaininfo.Name)
/*
        interfaceInfo, err := theforeman.AddBridgeInterface(labConfig.ForemanURL, session, "br-mgmt", curHost.Name, curdomaininfo.ID, subnetmap[fmt.Sprintf("%s-MGMT",labShortName)])
        if err != nil {
            return err
        }
*/
        interfaceInfo, err := theforeman.AddBridgeInterface(labConfig.ForemanURL, session, "br-lbaas", curHost.Name, curdomaininfo.ID, subnetmap[fmt.Sprintf("%s-LBAAS",labShortName)])
        if err != nil {
            return err
        }
        fmt.Println(interfaceInfo)
    }

    return nil

}



