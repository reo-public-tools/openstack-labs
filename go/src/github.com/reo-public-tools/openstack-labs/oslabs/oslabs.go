package main

import (
  "fmt"
  "log"
  "log/syslog"
  "github.com/reo-public-tools/openstack-labs/theforeman"
  "github.com/reo-public-tools/openstack-labs/osutils"
)

func main() {

    // Set up the syslog logger with some defaults
    sysLog, err := syslog.New(syslog.LOG_EMERG | syslog.LOG_USER, "oslabs")
    if err != nil {
        log.Fatal(err)
    }

    // Override the writers in the related packages
    // to change the identity id.  Later on we will
    // want to change where its writing for things
    // like remote logging.
    theforeman.OverrideLogWriter(sysLog)
    osutils.OverrideLogWriter(sysLog)

    // Set the url to your local forman server
    var url string = "https://172.20.41.28"
    _ = sysLog.Info(fmt.Sprintf("Setting foreman url to %s.", url))

    // Get the session id string
    _ = sysLog.Info("Logging into the foreman.")
    session, err := theforeman.TheForemanLogin(url)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(session)


/*
    // Find a static vlan backed domain available for use
    _ = sysLog.Info("Look for a free static domain.")
    curdomaininfo, err := theforeman.FindAvailableVLANDomain(url, session)
    fmt.Printf("Found available domain %s to check out\n", curdomaininfo.Name)
    if err != nil {
        log.Fatal(err)
    }

    _ = sysLog.Info(fmt.Sprintf("Checking out static lab %s\n", curdomaininfo.Name))
    fmt.Printf("Checking out lab %s\n", curdomaininfo.Name)
    err = theforeman.CheckOutVLANDomain(url, session, curdomaininfo.ID)
    if err != nil {
        log.Fatal(err)
    }

    _ = sysLog.Info(fmt.Sprintf("Releasing static lab %s\n", curdomaininfo.Name))
    fmt.Printf("Releasing lab %s\n", curdomaininfo.Name)
    err = theforeman.ReleaseVLANDomain(url, session, curdomaininfo.ID)
    if err != nil {
        log.Fatal(err)
    }



    _ = sysLog.Info("Creating a dynamic lab.")
    err = theforeman.CreateDynamicLab(url, session)
    if err != nil {
        log.Fatal(err)
    }

    _ = sysLog.Info("Releasing dynamic lab lab2.phobos.rpc.rackspace.com.")
    err = theforeman.DeleteDynamicLab(url, session, "lab2.phobos.rpc.rackspace.com")
    if err != nil {
        log.Fatal(err)
    }
*/
/*

    err = theforeman.DeleteDynamicLab(url, session, "lab3.phobos.rpc.rackspace.com")
    if err != nil {
        log.Fatal(err)
    }

*/

    // Auth to openstack test
    provider, err := osutils.OpenstackLogin("phobos")
    if err != nil {
        log.Fatal(err)
    }

    nodeList, err := osutils.GetAvailableIronicNodeListByCapability(&provider, "system_type", "standard")
    if err != nil {
        log.Fatal(err)
    }
    for _, curnode := range nodeList {
        fmt.Println(curnode.UUID)
    }
    //fmt.Println(nodeList)

    flavorCapability, err := osutils.GetFlavorCapability(&provider, "ironic-standard", "system_type")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(flavorCapability)


    ironicRequest := osutils.IronicOnIronicRequest{
            IronicOnIronicNodeRequests: []osutils.IronicOnIronicNodeRequest {
                { Flavor: "ironic-standard", Count: 5 },
                { Flavor: "ironic-storage-perf", Count: 3 },
            },
        }

    testIronic, err := osutils.CheckIronicCapacity(&provider, ironicRequest)
    if err != nil {
        log.Fatal(err)
    }
    if testIronic {
        fmt.Printf("Ironic on ironic capacity check good\n")
    }

    macList, err := osutils.GetIronicPXEMacs(&provider, "6fc8d998-848d-48d7-9e32-7594bc72e2e9")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(macList)


}
