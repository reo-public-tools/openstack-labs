package main

import (
  "fmt"
  "github.com/reo-public-tools/openstack-labs/theforeman"
  
)

func main() {

  // Set the url to your local forman server
  var url string = "https://172.20.41.28"

  // Get the session id string
  session := theforeman.TheForemanLogin(url)
  fmt.Printf("Cookie: %s\n", session)
}
