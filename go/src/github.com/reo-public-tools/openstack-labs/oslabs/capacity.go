package main

import (
    "fmt"
    "github.com/reo-public-tools/openstack-labs/osutils"
)

func Capacity() (error) {

    // Populate struct from config file
    _ = sysLog.Info(fmt.Sprintf("Pulling in global config used for listing.\n"))
    labConfig, err := PopulateConfigFile("")
    if err != nil {
        return err
    }

    // Connect to openstack
    provider, err := osutils.OpenstackLogin(labConfig.OpenstackCloud)
    if err != nil {
        return err
    }

    // Get a detailed list of ironic nodes
    ironicNodes, err := osutils.GetIronicCapacity(&provider)
    if err != nil {
        return err
    }

    // Print header
    fmt.Printf("\n\n%-30s %-30s %-10s %-10s\n", "Flavor", "Type", "Used", "Free")
    fmt.Printf("=============================================================================================\n")

    // Print node info
    for _, node := range ironicNodes.NodesByType {
        fmt.Printf("%-30s %-30s %-10d %-10d\n", node.CapacityType, node.Flavor, node.Used, node.Free)
    }
    fmt.Printf("\n\n")

    return nil

}

