package main

import (
    "fmt"
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

    // Loop over hosts and delete each
    for _, curHost := range detailedHostList {

        _ = sysLog.Info(fmt.Sprintf("Deleting host: %s\n", curHost.Name))
        fmt.Printf("Deleting host: %s.\n", curdomaininfo.Name)
        err = theforeman.DeleteHost(labConfig.ForemanURL, session, curHost.Name)
        if err != nil {
            return err
        }
    }

    return nil

}



