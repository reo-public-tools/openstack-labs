package main

import (
    "fmt"
    "strings"
    "github.com/reo-public-tools/openstack-labs/theforeman"
)

func List() (error) {

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

    // Get detailed domain list
    _ = sysLog.Info(fmt.Sprintf("Pulling detailed domain list: %s\n"))
    domainDetails, err := theforeman.GetDomainsWithDetails(labConfig.ForemanURL, session)
    if err != nil {
        return err
    }

    // Print header
    fmt.Printf("+-%s-+-%s-+-%s-+-%s-+\n",
            strings.Repeat("-", 30),
            strings.Repeat("-", 7),
            strings.Repeat("-", 7),
            strings.Repeat("-", 40),
        )
    fmt.Printf("| %-30s | %-7s | %-7s | %-40s |\n", "Domain Name", "Type", "In Use", "Description")
    fmt.Printf("+-%s-+-%s-+-%s-+-%s-+\n",
            strings.Repeat("-", 30),
            strings.Repeat("-", 7),
            strings.Repeat("-", 7),
            strings.Repeat("-", 40),
        )

    // Loop and print details
    for _, domainInfo := range domainDetails {
        curType := ""
        inUse := ""
        for _, param := range domainInfo.Parameters {
            switch pname := param.Name; pname {
            case "type":
                curType = param.Value
            case "in-use":
                inUse = param.Value
            }
        }
        fmt.Printf("| %-30s | %-7s | %-7s | %-40s |\n", domainInfo.Name, curType, inUse, domainInfo.Fullname)
        fmt.Printf("+-%s-+-%s-+-%s-+-%s-+\n",
                strings.Repeat("-", 30),
                strings.Repeat("-", 7),
                strings.Repeat("-", 7),
            strings.Repeat("-", 40),
            )
    }

    return nil

}

