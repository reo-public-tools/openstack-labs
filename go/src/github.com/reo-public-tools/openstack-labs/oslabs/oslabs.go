package main

import (
  "fmt"
  "log"
  "github.com/reo-public-tools/openstack-labs/theforeman"
)

func main() {

    // Set the url to your local forman server
    var url string = "https://172.20.41.28"

    // Get the session id string
    session, err := theforeman.TheForemanLogin(url)
    if err != nil {
        log.Fatal(err)
    }

    // Find a static vlan backed domain available for use
    /*
    curdomaininfo, err := theforeman.FindAvailableVLANDomain(url, session)
    fmt.Printf("Found available domain %s to check out\n", curdomaininfo.Name)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Checking out lab %s\n", curdomaininfo.Name)
    err = theforeman.CheckOutVLANDomain(url, session, curdomaininfo.ID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Releasing lab %s\n", curdomaininfo.Name)
    err = theforeman.ReleaseVLANDomain(url, session, curdomaininfo.ID)
    if err != nil {
        log.Fatal(err)
    }


    */

    err = theforeman.CreateDynamicLab(url, session)
    if err != nil {
        log.Fatal(err)
    }

    err = theforeman.DeleteDynamicLab(url, session, "lab2.phobos.rpc.rackspace.com")
    if err != nil {
        log.Fatal(err)
    }
/*

    err = theforeman.DeleteDynamicLab(url, session, "lab3.phobos.rpc.rackspace.com")
    if err != nil {
        log.Fatal(err)
    }

*/

    fmt.Println("finished")


}
