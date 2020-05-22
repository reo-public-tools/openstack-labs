package main

import (
    "fmt"
    "github.com/reo-public-tools/openstack-labs/theforeman"
    "github.com/reo-public-tools/openstack-labs/osutils"
)

func GenerateInventory(labName string) (error) {

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

    fmt.Printf("---\nall:\n");
    fmt.Printf("  vars:\n")
    fmt.Printf("    domain_name: %s\n", domainInfo.Name)
    fmt.Printf("    domain_description: %s\n", domainInfo.Fullname)
    fmt.Printf("    domain_parameters:\n")

    // Display settings from domain parameters
    var curType string
    hasIOIData := false
    for _, param := range domainInfo.Parameters {
        fmt.Printf("      %s: %s\n", param.Name, param.Value)
        if param.Name == "type" {
            curType = param.Value.(string)
        }
    }

    fmt.Printf("    subnets:\n")
    for _, subnet := range domainInfo.Subnets {

        // Pull a more detailed subnet info for the current subnet
        _ = sysLog.Info(fmt.Sprintf("Pulling detailed subnet info for subnet: %s\n", subnet.Name))
        subnetInfo, err := theforeman.GetSubnetDetails(labConfig.ForemanURL, session, subnet.Name)
        if err != nil {
            return err
        }

        fmt.Printf("      %s:\n", subnetInfo.Name)
        fmt.Printf("        description: %s\n", subnet.Description)
        fmt.Printf("        network_address: %s\n", subnet.NetworkAddress)
        fmt.Printf("        gateway: %s\n", subnetInfo.Gateway)
        fmt.Printf("        netmask: %s\n", subnetInfo.Mask)
        fmt.Printf("        network: %s\n", subnetInfo.Network)
        fmt.Printf("        network_type: %s\n", subnetInfo.NetworkType)
        fmt.Printf("        dhcp_from: %s\n", subnetInfo.From)
        fmt.Printf("        dhcp_to: %s\n", subnetInfo.To)
        fmt.Printf("        ipam: %s\n", subnetInfo.Ipam)
        if curType == "vlan" {
            fmt.Printf("        vlan: %d\n", subnetInfo.Vlanid)
        } else {

            // Get the vxlan id from parameters
            for _, curParam := range subnetInfo.Parameters {
                  if curParam.Name == "vxlan-id" {
                      fmt.Printf("        vxlan: %s\n", curParam.Value)
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

        fmt.Printf("    ironic_on_ironic_nodes:\n")
        for _, IOINode := range IOINodeList.IronicNodeDetails {
            fmt.Printf("      - id: %s\n", IOINode.ID)
            fmt.Printf("        ipmi_address: %s\n", IOINode.IpmiAddress)
            fmt.Printf("        ipmi_username: %s\n", IOINode.IpmiUsername)
            fmt.Printf("        ipmi_password: %s\n", IOINode.IpmiPassword)
            fmt.Printf("        cpus: %s\n", IOINode.Cpus)
            fmt.Printf("        memory_mb: %s\n", IOINode.MemoryMb)
            fmt.Printf("        size: %s\n", IOINode.Size)
            fmt.Printf("        flavor: %s\n", IOINode.Flavor)
            fmt.Printf("        system_type: %s\n", IOINode.SystemType)
            fmt.Printf("        macs: %s\n", IOINode.Macs)
        }
    }


    // Display host listing
    _ = sysLog.Info(fmt.Sprintf("Pulling host listing for domain: %s\n", domainInfo.Name))
    detailedHostList, err := theforeman.GetHostsDetailsByDomainID(labConfig.ForemanURL, session, domainInfo.ID)
    if err != nil {
        return err
    }
    if len(detailedHostList) > 0 {
        fmt.Printf("  hosts:\n")
        for _, curHost := range detailedHostList {
            fmt.Printf("    %s:\n", curHost.Name)
            fmt.Printf("      ansible_host: %s\n", curHost.IP)
            fmt.Printf("      ansible_user: %s\n", "admin")
            fmt.Printf("      arch: %s\n", curHost.ArchitectureName)
            fmt.Printf("      os: %s\n", curHost.OperatingsystemName)
            fmt.Printf("      interfaces:\n")
            for _, curInterface := range curHost.Interfaces {
                if ! curInterface.Primary {
                    fmt.Printf("        - ip: %s\n", curInterface.IP)
                    fmt.Printf("          subnet: %s\n", curInterface.SubnetName)
                    fmt.Printf("          domain_name: %s\n", curInterface.DomainName)
                    fmt.Printf("          managed: %t\n", curInterface.Managed)
                    fmt.Printf("          primary: %t\n", curInterface.Primary)
                    fmt.Printf("          provision: %t\n", curInterface.Provision)
                    fmt.Printf("          virtual: %t\n", curInterface.Virtual)
                    fmt.Printf("          type: %s\n", curInterface.Type)
                    fmt.Printf("          mtu: %d\n", curInterface.Mtu)
                    fmt.Printf("          fqdn: %s\n", curInterface.Fqdn)
                    fmt.Printf("          mode: %s\n", curInterface.Mode)
                    fmt.Printf("          bond_options: %s\n", curInterface.BondOptions)
                }
            }
        }
    }

    return nil

}

