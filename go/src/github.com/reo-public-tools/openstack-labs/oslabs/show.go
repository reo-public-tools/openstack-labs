package main

import (
    "fmt"
    "github.com/reo-public-tools/openstack-labs/theforeman"
    "github.com/reo-public-tools/openstack-labs/osutils"
)

func Show(labName string) (error) {

    // Populate struct from config file
    _ = sysLog.Info(fmt.Sprintf("Pulling in global config used for listing.\n"))
    labConfig, err := PopulateConfigFile("")
    if err != nil {
        return err
    }

    // Connect to theforeman
    _ = sysLog.Info(fmt.Sprintf("Logging into the foreman at: %s\n", labConfig.ForemanURL))
    session, err := theforeman.TheForemanLogin(labConfig.ForemanURL)
    if err != nil {
        return err
    }

    // Get detailed domain info
    _ = sysLog.Info(fmt.Sprintf("Pulling detailed domain info for lab: %s\n", labName))
    domainInfo, err := theforeman.GetDomainDetails(labConfig.ForemanURL, session, labName)
    if err != nil {
        return err
    }

    // Display general info
    fmt.Printf("Domain Info:\n")
    fmt.Printf("  Name: %s\n", domainInfo.Name)
    fmt.Printf("  Description: %s\n", domainInfo.Fullname)
    fmt.Printf("  Domain Parameters:\n")

    // Display settings from domain parameters
    var curType string
    hasIOIData := false
    for _, param := range domainInfo.Parameters {
        switch pname := param.Name; pname {
        case "type":
            fmt.Printf("    %-35s %s\n", "Type", param.Value)
            curType = param.Value
        case "in-use":
            fmt.Printf("    %-35s %s\n", "In Use:", param.Value)
        case "ioi_data":
            fmt.Printf("    %-35s %s\n", "Has Ironic-On-Ironic nodes:", "yes")
            hasIOIData = true
        case "multicast_gruop":
            fmt.Printf("    %-35s %s\n", "Multicast Group:", param.Value)
        case "internal_vrid":
            fmt.Printf("    %-35s %s\n", "Haproxy/Keepalived Internal VRID:", param.Value)
        case "external_vrid":
            fmt.Printf("    %-35s %s\n", "Haproxy/Keepalived External VRID:", param.Value)
        case "internal_floating_ip":
            fmt.Printf("    %-35s %s\n", "Haproxy/Keepalived Internal IP:", param.Value)
        case "external_floating_ip":
            fmt.Printf("    %-35s %s\n", "Haproxy/Keepalived External IP:", param.Value)
        default:
            fmt.Printf("    %-35s %s\n", pname + ":", param.Value)
        }
    }

    fmt.Printf("\n  Subnet Info:\n")
    for _, subnet := range domainInfo.Subnets {

        // Pull a more detailed subnet info for the current subnet
        _ = sysLog.Info(fmt.Sprintf("Pulling detailed subnet info for subnet: %s\n", subnet.Name))
        subnetInfo, err := theforeman.GetSubnetDetails(labConfig.ForemanURL, session, subnet.Name)
        if err != nil {
            return err
        }

        fmt.Printf("    Subnet: %s\n", subnetInfo.Name)
        fmt.Printf("      %-20s %s\n", "Description:", subnet.Description)
        fmt.Printf("      %-20s %s\n", "Network Address:", subnet.NetworkAddress)
        fmt.Printf("      %-20s %s\n", "Gateway:", subnetInfo.Gateway)
        fmt.Printf("      %-20s %s-%s\n", "DHCP Range:", subnetInfo.From, subnetInfo.To)
        fmt.Printf("      %-20s %s\n", "Ipam:", subnetInfo.Ipam)
        if curType == "vlan" {
            fmt.Printf("      %-20s %d\n", "VLAN:", subnetInfo.Vlanid)
        } else {

            // Get the vxlan id from parameters
            for _, curParam := range subnetInfo.Parameters {
                  if curParam.Name == "vxlan-id" {
                      fmt.Printf("      %-20s %s\n", "VXLAN ID:", curParam.Value)
                  }
            }

        } // End "if curType == "vlan" {"
    } // End "for _, subnet := range domainInfo.Subnets {"


    if hasIOIData {
        // Pull the ioi_data parameter
        _ = sysLog.Info(fmt.Sprintf("Pulling IOI data  for domain: %s\n", domainInfo.Name))
        IOIData, err := theforeman.GetDomainParameter(labConfig.ForemanURL, session, domainInfo.Name, "ioi_data")
        if err != nil {
            return err
        }

        // Take the base64 encoded json string and convert back to a noda data struct
        _ = sysLog.Info(fmt.Sprintf("Decoding basd64 data and clearing domain param for domain: %s\n", domainInfo.Name))
        IOINodeList, err := osutils.JSONStringToNodeData(IOIData, true)
        if err != nil {
            return err
        }

        fmt.Printf("\n  Ironic On Ironic Nodes:\n")
        for _, IOINode := range IOINodeList.IronicNodeDetails {
            fmt.Printf("    ID: %s\n", IOINode.ID)
            fmt.Printf("      %-20s %s,%s,%s\n", "IPMI IP,User,Pass:", IOINode.IpmiAddress,
                                                                       IOINode.IpmiUsername,
                                                                       IOINode.IpmiPassword)
            fmt.Printf("      %-20s %.0f,%.0f,%.0f\n", "Cpus,Mem,Disk", IOINode.Cpus,
                                                                        IOINode.MemoryMb,
                                                                        IOINode.Size)
            fmt.Printf("      %-20s %s,%s\n", "Flavor,Capability:", IOINode.Flavor,
                                                                    IOINode.SystemType)
            //fmt.Printf("      %-20s %s\n", "Macs:", IOINode.Macs)
        }
    }


    // Display host listing
    _ = sysLog.Info(fmt.Sprintf("Pulling host listing for domain: %s\n", domainInfo.Name))
    detailedHostList, err := theforeman.GetHostsDetailsByDomainID(labConfig.ForemanURL, session, domainInfo.ID)
    if err != nil {
        return err
    }
    if len(detailedHostList) > 0 {
        fmt.Printf("\n  Host Information\n")
        for _, curHost := range detailedHostList {
            fmt.Printf("    Name: %s\n", curHost.Name)
            fmt.Printf("      %-35s %s\n", "Primary/Host IP:", curHost.IP)
            fmt.Printf("      %-35s %s\n", "Arch:", curHost.ArchitectureName)
            fmt.Printf("      %-35s %s\n", "OS:", curHost.OperatingsystemName)
            fmt.Printf("      %-35s %s\n", "Puppet Env:", curHost.EnvironmentName)
            fmt.Printf("      Allocated IPs:\n")
            for _, curInterface := range curHost.Interfaces {
                if ! curInterface.Primary {
                    fmt.Printf("        %-20s %s\n", curInterface.SubnetName + ":", curInterface.IP)
                }
            }
        }
    }

    return nil

}

