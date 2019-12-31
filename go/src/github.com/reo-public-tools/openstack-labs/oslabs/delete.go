package main

import (
    "fmt"
    "github.com/reo-public-tools/openstack-labs/theforeman"
    "github.com/reo-public-tools/openstack-labs/osutils"
)

func Delete(labName string) (error) {

    // Populate struct from config file
    _ = sysLog.Info(fmt.Sprintf("Pulling in global config.\n"))
    fmt.Printf("Pulling in global config.\n")
    labConfig, err := PopulateConfigFile("")
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
    fmt.Printf("Pullinginfo for domain: %s.\n", labName)
    curdomaininfo, err := theforeman.GetDomainDetails(labConfig.ForemanURL, session, labName)
    if err != nil {
        return err
    }

    // Get the lab type
    labType := ""
    for _, param := range curdomaininfo.Parameters {
        if param.Name == "type" {
            if param.Value == "vlan" {
                labType = "static"
            }
        }
    }

    /* #### Delete the neutron port and clear the parameter #### */

    // Connect to openstack
    provider, err := osutils.OpenstackLogin(labConfig.OpenstackCloud)
    if err != nil {
        return err
    }

    // Delete the domain parameter
    _ = sysLog.Info(fmt.Sprintf("Clearing the external_floating_ip parameter\n"))
    fmt.Printf("Clearing the external_floating_ip parameter\n")
    err = theforeman.DeleteDomainParameter(labConfig.ForemanURL, session, labName, "external_floating_ip")
    if err != nil {
        return err
    }

    // Delete the port
    _ = sysLog.Info(fmt.Sprintf("Deleting the neutron port tied to the external_floating_ip parameter\n"))
    fmt.Printf("Deleting the neutron port tied to the external_floating_ip parameter\n")
    err = osutils.DeleteExtNetworkPortByLabName(&provider, labName)
    if err != nil {
        return err
    }


    /* #### Delete a dynamic(vxlan based) domain or check in a vlan domain #### */
    if labType == "static" {

        // Release the static lab
        _ = sysLog.Info(fmt.Sprintf("Releasing static lab %s\n", curdomaininfo.Name))
        fmt.Printf("Releasing lab %s\n", curdomaininfo.Name)
        err = theforeman.ReleaseVLANDomain(labConfig.ForemanURL, session, curdomaininfo.ID)
        if err != nil {
            return err
        }

    } else {

        // Delete the dynamic lab
        _ = sysLog.Info(fmt.Sprintf("Deleting dynamic lab %s.\n", curdomaininfo.Name))
        fmt.Printf("Deleting dynamic lab %s.\n", curdomaininfo.Name)
        err = theforeman.DeleteDynamicLab(labConfig.ForemanURL, session, curdomaininfo.Name)
        if err != nil {
            return err
        }

    }




    return nil

}



