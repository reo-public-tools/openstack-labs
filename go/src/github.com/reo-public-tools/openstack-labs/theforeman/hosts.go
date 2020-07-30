package theforeman

import (
    "fmt"
    "bytes"
    "strings"
    "encoding/json"
    "net/http"
    "crypto/tls"
    "io/ioutil"
)


// Structures used when pulling back a host listing
type HostQueryResults struct {
        Total    int         `json:"total"`
        Subtotal int         `json:"subtotal"`
        Page     int         `json:"page"`
        PerPage  int         `json:"per_page"`
        Search   string      `json:"search"`
        Sort     Sort        `json:"sort"`
        Results  []Host      `json:"results"`
}

// Structures used to hold the returned host value
type Host struct {
        Name                     string                     `json:"name"`
        IP                       string                     `json:"ip,omitempty"`
        IP6                      string                     `json:"ip6,omitempty"`
        EnvironmentID            int                        `json:"environment_id,omitempty"`
        EnvironmentName          string                     `json:"environment_name,omitempty"`
        Mac                      string                     `json:"mac,omitempty"`
        RealmID                  int                        `json:"realm_id,omitempty"`
        RealmName                string                     `json:"realm_name,omitempty"`
        SpMac                    string                     `json:"sp_mac,omitempty"`
        SpIP                     string                     `json:"sp_ip,omitempty"`
        SpName                   string                     `json:"sp_name,omitempty"`
        DomainID                 int                        `json:"domain_id,omitempty"`
        DomainName               string                     `json:"domain_name,omitempty"`
        ArchitectureID           int                        `json:"architecture_id,omitempty"`
        ArchitectureName         string                     `json:"architecture_name,omitempty"`
        OperatingsystemID        int                        `json:"operatingsystem_id,omitempty"`
        OperatingsystemName      string                     `json:"operatingsystem_name,omitempty"`
        SubnetID                 int                        `json:"subnet_id,omitempty"`
        SubnetName               string                     `json:"subnet_name,omitempty"`
        Subnet6ID                int                        `json:"subnet6_id,omitempty"`
        Subnet6Name              string                     `json:"subnet6_name,omitempty"`
        SpSubnetID               int                        `json:"sp_subnet_id,omitempty"`
        PtableID                 int                        `json:"ptable_id,omitempty"`
        PtableName               string                     `json:"ptable_name,omitempty"`
        MediumID                 int                        `json:"medium_id,omitempty"`
        MediumName               string                     `json:"medium_name,omitempty"`
        PXELoader                string                     `json:"pxe_loader,omitempty"`
        Build                    bool                       `json:"build,omitempty"`
        Comment                  string                     `json:"comment,omitempty"`
        Disk                     interface{}                `json:"disk,omitempty"`
        InstalledAt              interface{}                `json:"instanlled_at,omitempty"`
        ModelID                  int                        `json:"model_id,omitempty"`
        HostgroupID              int                        `json:"hostgroup_id,omitempty"`
        OwnerID                  int                        `json:"owner_id,omitempty"`
        OwnerName                string                     `json:"owner_name,omitempty"`
        OwnerType                string                     `json:"owner_type,omitempty"`
        Enabled                  bool                       `json:"enabled,omitempty"`
        Managed                  bool                       `json:"managed,omitempty"`
        UseImage                 interface{}                `json:"use_image,omitempty"`
        ImageFile                string                     `json:"image_file,omitempty"`
        UUID                     string                     `json:"uuid,omitempty"`
        ComputeResourceID        int                        `json:"compute_resource_id,omitempty"`
        ComputeResourceName      string                     `json:"compute_resource_name,omitempty"`
        ComputeProfileID         int                        `json:"compute_profile_id,omitempty"`
        ComputeProfileName       string                     `json:"compute_profile_name,omitempty"`
        Capabilities             []string                   `json:"capabilities,omitempty"`
        CertName                 string                     `json:"certname,omitempty"`
        ImageID                  int                        `json:"image_id,omitempty"`
        ImageName                string                     `json:"image_name,omitempty"`
        CreatedAt                string                     `json:"created_at,omitempty"`
        UpdatedAt                string                     `json:"updated_at,omitempty"`
        LastCompile              interface{}                `json:"last_compile,omitempty"`
        GlobalStatus             int                        `json:"global_status,omitempty"`
        GlobalStatusLabel        string                     `json:"global_status_label,omitempty"`
        UptimeSeconds            interface{}                `json:"uptime_seconds,omitempty"`
        OrganizationID           int                        `json:"organization_id,omitempty"`
        OrganizationName         string                     `json:"organization_name,omitempty"`
        LocationID               int                        `json:"location_id,omitempty"`
        LocationName             string                     `json:"location_name,omitempty"`
        PuppetStatus             int                        `json:"puppet_proxy_status,omitempty"`
        ModelName                string                     `json:"model_name,omitempty"`
        ConfigurationStatus      int                        `json:"configuration_status,omitempty"`
        ConfigurationStatusLabel string                     `json:"configuration_status_label,omitempty"`
        BuildStatus              int                        `json:"build_status,omitempty"`
        BuildStatusLabel         string                     `json:"build_status_label,omitempty"`
        ID                       int                        `json:"id,omitempty"`
        PuppetProxyID            int                        `json:"puppet_proxy_id,omitempty"`
        PuppetProxyName          string                     `json:"puppet_proxy_name,omitempty"`
        PuppetCAProxyID          int                        `json:"puppet_ca_proxy_id,omitempty"`
        PuppetCAProxyName        string                     `json:"puppet_ca_proxy_name,omitempty"`
        PuppetProxy              PuppetProxy                `json:"puppet_proxy,omitempty"`
        PuppetCaProxy            PuppetCaProxy              `json:"puppet_ca_proxy,omitempty"`
        Parameters               []interface{}              `json:"parameters,omitempty"`
        HostGroupName            string                     `json:"host_group_name,omitempty"`
        HostGroupTitle           string                     `json:"host_group_title,omitempty"`
        AllParameters            []HostParameters           `json:"all_parameters,omitempty"`
        Interfaces               []NetInterfaceRet          `json:"interfaces,omitempty"`
        PuppeClasses             []interface{}              `json:"puppetclasses,omitempty"`
        ConfigGroups             []interface{}              `json:"config_groups,omitempty"`
        AllPuppeClasses          []interface{}              `json:"all_puppetclasses,omitempty"`
}
type PuppetProxy struct {
        Name string `json:"name,omitempty"`
        ID   int    `json:"id,omitempty"`
        URL  string `json:"url,omitempty"`
}
type PuppetCaProxy struct {
        Name string `json:"name,omitempty"`
        ID   int    `json:"id,omitempty"`
        URL  string `json:"url,omitempty"`
}
type HostParameters struct {
        Name          string      `json:"name,omitempty"`
        Priority      interface{} `json:"priority,omitempty"`
        ID            int         `json:"id,omitempty"`
        Value         interface{} `json:"value,omitempty"`
        ParameterType string      `json:"parameter_type,omitempty"`
        CreatedAt     string      `json:"created_at,omitempty"`
        UpdatedAt     string      `json:"updated_at,omitempty"`
        Permissions   Permissions `json:"permissions,omitempty"`
}
type NetInterfaceRet struct {
        SubnetID        int         `json:"subnet_id,omitempty"`
        SubnetName      string      `json:"subnet_name,omitempty"`
        Subnet6ID       int         `json:"subnet6_id,omitempty"`
        Subnet6Name     string      `json:"subnet6_name,omitempty"`
        DomainID        int         `json:"domain_id,omitempty"`
        DomainName      string      `json:"domain_name,omitempty"`
        CreatedAt       string      `json:"created_at,omitempty"`
        UpdatedAt       string      `json:"updated_at,omitempty"`
        Managed         bool        `json:"managed,omitempty"`
        Identifier      string      `json:"identifier,omitempty"`
        ID              int         `json:"id,omitempty"`
        Name            string      `json:"name,omitempty"`
        IP              string      `json:"ip,omitempty"`
        IP6             string      `json:"ip6,omitempty"`
        Mac             string      `json:"mac,omitempty"`
        Mtu             int         `json:"mtu,omitempty"`
        Fqdn            string      `json:"fqdn,omitempty"`
        Primary         bool        `json:"primary,omitempty"`
        Provision       bool        `json:"provision,omitempty"`
        Type            string      `json:"type,omitempty"`
        Virtual         bool        `json:"virtual,omitempty"`
        Username        string      `json:"username,omitempty"`
        Password        string      `json:"password,omitempty"`
        Provider        string      `json:"provider,omitempty"`
        Tag             string      `json:"tag,omitempty"`
        AttachedTo      string      `json:"attached_to,omitempty"`
        Mode            string      `json:"mode,omitempty"`
        AttachedDevices interface{} `json:"attached_devices,omitempty"`
        BondOptions     string      `json:"bond_options,omitempty"`
}
type Permissions struct {
        ViewHosts      bool `json:"view_hosts,omitempty"`
        CreateHosts    bool `json:"create_hosts,omitempty"`
        EditHosts      bool `json:"edit_hosts,omitempty"`
        DestroyHosts   bool `json:"destroy_hosts,omitempty"`
        BuildHosts     bool `json:"build_hosts,omitempty"`
        PowerHosts     bool `json:"power_hosts,omitempty"`
        ConsoleHosts   bool `json:"console_hosts,omitempty"`
        IpmiBootHosts  bool `json:"ipmi_boot_hosts,omitempty"`
        PuppetrunHosts bool `json:"puppetrun_hosts,omitempty"`
}

// Structures used to create a new host
type HostPostData struct {
        OrganizationID int           `json:"organization_id,omitempty"`
        LocationID     int           `json:"location_id,omitempty"`
        Host           NewHostData   `json:"host"`
}

type NewHostData struct {
        Name                     string                     `json:"name"`
        LocationID               int                        `json:"location_id,omitempty"`
        OrganizationID           int                        `json:"organization_id,omitempty"`
        EnvironmentID            string                     `json:"environment_id,omitempty"`
        IP                       string                     `json:"ip,omitempty"`
        Mac                      string                     `json:"mac,omitempty"`
        ArchitectureID           int                        `json:"architecture_id,omitempty"`
        DomainID                 int                        `json:"domain_id,omitempty"`
        RealmID                  int                        `json:"realm_id,omitempty"`
        PuppetProxyID            int                        `json:"puppet_proxy_id,omitempty"`
        PuppetCAProxyID          int                        `json:"puppet_ca_proxy_id,omitempty"`
        PuppeClaassIDs           []interface{}              `json:"puppet_class_ids,omitempty"`
        ConfigGroupIDs           []interface{}              `json:"config_group_ids,omitempty"`
        OperatingsystemID        int                        `json:"operatingsystem_id,omitempty"`
        MediumID                 int                        `json:"medium_id,omitempty"`
        PXELoader                string                     `json:"pxe_loader,omitempty"`
        PtableID                 int                        `json:"ptable_id,omitempty"`
        SubnetID                 int                        `json:"subnet_id,omitempty"`
        ComputeResourceID        int                        `json:"compute_resource_id,omitempty"`
        RootPass                 string                     `json:"root_pass,omitempty"`
        ModelID                  int                        `json:"model_id,omitempty"`
        HostgroupID              int                        `json:"hostgroup_id,omitempty"`
        OwnerID                  int                        `json:"owner_id,omitempty"`
        OwnerType                string                     `json:"owner_type,omitempty"`
        ImageID                  int                        `json:"image_id,omitempty"`
        HostParametersAttributes []HostParametersAttributes `json:"host_parameters_attributes,omitempty"`
        Build                    bool                       `json:"build,omitempty"`
        Enabled                  bool                       `json:"enabled,omitempty"`
        ProvisionMethod          string                     `json:"provision_method,omitempty"`
        Managed                  bool                       `json:"managed,omitempty"`
        ProgressReportID         string                     `json:"progress_report_id,omitempty"`
        Comment                  string                     `json:"comment,omitempty"`
        Capabilities             string                     `json:"capabilities,omitempty"`
        ComputeProfileID         int                        `json:"compute_profile_id,omitempty"`
        InterfacesAttributes     InterfacesAttributes       `json:"interfaces_attributes,omitempty"`
        ComputeAttributes        ComputeAttributes          `json:"compute_attributes,omitempty"`
}

type HostParametersAttributes struct {
        Name  string `json:"name,omitempty"`
        Value string `json:"value,omitempty"`
}

type InterfacesAttributes struct {
                Primary            NetInterface `json:"1,omitempty"`
                Management         NetInterface `json:"2,omitempty"`
                Storage            NetInterface `json:"3,omitempty"`
                StorageManagement  NetInterface `json:"4,omitempty"`
                Tenant             NetInterface `json:"5,omitempty"`
                LBAAS              NetInterface `json:"6,omitempty"`
                InsideNet          NetInterface `json:"7,omitempty"`
                GatewayNet         NetInterface `json:"8,omitempty"`
}

type NetInterface struct {
        Name            string   `json:"name,omitempty"`
        Primary         bool     `json:"primary,omitempty"`
        IP              string   `json:"ip,omitempty"`
        IP6             string   `json:"ip6,omitempty"`
        Mac             string   `json:"mac,omitempty"`
        Type            string   `json:"type,omitempty"`
        SubnetID        int      `json:"subnet_id,omitempty"`
        Subnet6ID       int      `json:"subnet6_id,omitempty"`
        DomainID        int      `json:"domain_id,omitempty"`
        Identifier      string   `json:"identifier,omitempty"`
        Managed         bool     `json:"managed,omitempty"`
        Provision       bool     `json:"provision,omitempty"`
        Username        string   `json:"username,omitempty"`
        Password        string   `json:"password,omitempty"`
        Provider        string   `json:"provider,omitempty"`
        Virtual         bool     `json:"virtual,omitempty"`
        Tag             string   `json:"tag,omitempty"`
        Mtu             int      `json:"mtu,omitempty"`
        AttachedTo      string   `json:"attached_to,omitempty"`
        Mode            string   `json:"mode,omitempty"`
        AttachedDevices []string `json:"attached_devices,omitempty"`
        BondOptions     string   `json:"bond_options,omitempty"`
}

type ComputeAttributes struct {
        Start           string   `json:"start,omitempty"`
}

// Structures used to create a new interface
type InterfacePostData struct {
        OrganizationID int     `json:"organization_id,omitempty"`
        LocationID     int     `json:"location_id,omitempty"`
        HostID         string  `json:"host_id,omitempty"`
        Interface NetInterface `json:"interface"`
}

// Create a new host
func CreateHost(url string,
                session string,
                labName string,
                hostName string,
                hostGroup string,
                role string) (Host, error) {

    sysLogPrefix := "theforeman(package).hosts(file).CreateHost(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating a new host by the name of %s.", sysLogPrefix, hostName))

    // Init the return value
    var queryResults Host

    // Get a struct populated with common parameter data
    globalParams, err := GetGlobalParameters(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Convert the location name to id for host creation
    locationID, err := ConvLocNameToID(url, session, globalParams.LabLocationName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Convert the organization name to id for host creation
    organizationID, err := ConvOrgNameToID(url, session, globalParams.LabOrgName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Convert the host group name to an id for host creation
    hostGroupID, err := ConvHostGroupNameToID(url, session, hostGroup)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Get domain details to use in the post data
    domainInfo, err := GetDomainDetails(url, session, labName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Set up some subnet name -> id mappings 
    subnetmap := make(map[string]int)
    for _, curSubnet := range domainInfo.Subnets {
        subnetmap[curSubnet.Name] = curSubnet.ID
    }

    // Pull the lab name from the fqdn of the lab name
    labShortName := strings.ToUpper(strings.Split(labName, ".")[0])


    // Fill out the data structure to be used for creating the domain
    newHostStruct := HostPostData{
        Host: NewHostData {
            Name: hostName,
            OrganizationID: organizationID,
            LocationID: locationID,
            HostgroupID: hostGroupID,
            DomainID: domainInfo.ID,
            InterfacesAttributes: InterfacesAttributes{
                Primary: NetInterface{
                    Name: hostName,
                    DomainID: domainInfo.ID,
                    Primary: true,
                    Managed: false,
                    Provision: false,
                },
                Management: NetInterface{
                    Identifier: "br-mgmt",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-MGMT",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                Storage: NetInterface{
                    Identifier: "br-storage",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-STORAGE",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                StorageManagement: NetInterface{
                    Identifier: "br-stor-mgmt",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-STOR-MGMT",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                Tenant: NetInterface{
                    Identifier: "br-vxlan",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-TENANT",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                LBAAS: NetInterface{
                    Identifier: "br-lbaas",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-LBAAS",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                InsideNet: NetInterface{
                    Identifier: "neutron-internal",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-INSIDE-NET",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                GatewayNet: NetInterface{
                    Identifier: "neutron-gateway",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-GW-NET",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
            },
            HostParametersAttributes: []HostParametersAttributes {
                { Name: "role", Value: role },
            },
        },
    }

    // Convert data to json
    postData, err := json.Marshal(newHostStruct)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/hosts", url)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("POST", requesturl, bytes.NewBufferString(string(postData)))
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }

    // Make sure we are using the proper content type for the configs api
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json,version=2")

    // Set the session Cookie header
    req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", session))

    // Disable tls verify
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    // Set up the http client and do the request
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }
    defer resp.Body.Close()


    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 201 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return Host{}, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }

    return queryResults, nil

}

// Create a new host
func CreateSimulatedHost(url string,
                         session string,
                         labName string,
                         hostName string,
                         hostGroup string,
                         role string) (Host, error) {

    sysLogPrefix := "theforeman(package).hosts(file).CreateHost(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating a new host by the name of %s.", sysLogPrefix, hostName))

    // Init the return value
    var queryResults Host

    // Get a struct populated with common parameter data
    globalParams, err := GetGlobalParameters(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Convert the location name to id for host creation
    locationID, err := ConvLocNameToID(url, session, globalParams.LabLocationName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Convert the organization name to id for host creation
    organizationID, err := ConvOrgNameToID(url, session, globalParams.LabOrgName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Convert the host group name to an id for host creation
    hostGroupID, err := ConvHostGroupNameToID(url, session, hostGroup)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Get domain details to use in the post data
    domainInfo, err := GetDomainDetails(url, session, labName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return Host{}, err
    }

    // Set up some subnet name -> id mappings
    subnetmap := make(map[string]int)
    for _, curSubnet := range domainInfo.Subnets {
        subnetmap[curSubnet.Name] = curSubnet.ID
    }

    // Pull the lab name from the fqdn of the lab name
    labShortName := strings.ToUpper(strings.Split(labName, ".")[0])

    // Fill out the data structure to be used for creating the domain
    newHostStruct := HostPostData{
        Host: NewHostData {
            Name: hostName,
            OrganizationID: organizationID,
            LocationID: locationID,
            HostgroupID: hostGroupID,
            DomainID: domainInfo.ID,
            ProvisionMethod: "image",
            ComputeAttributes: ComputeAttributes{
                Start: "1",
            },
            InterfacesAttributes: InterfacesAttributes{
                Primary: NetInterface{
                    Name: hostName,
                    DomainID: domainInfo.ID,
                    Primary: true,
                    Managed: true,
                    Provision: true,
                    SubnetID: subnetmap["MNAIO-PROV"],
                },
                Management: NetInterface{
                    Identifier: "br-mgmt",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-MGMT",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                Storage: NetInterface{
                    Identifier: "br-storage",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-STORAGE",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                StorageManagement: NetInterface{
                    Identifier: "br-stor-mgmt",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-STOR-MGMT",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                Tenant: NetInterface{
                    Identifier: "br-vxlan",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-TENANT",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                LBAAS: NetInterface{
                    Identifier: "br-lbaas",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-LBAAS",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                InsideNet: NetInterface{
                    Identifier: "neutron-internal",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-INSIDE-NET",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
                GatewayNet: NetInterface{
                    Identifier: "neutron-gateway",
                    DomainID: domainInfo.ID,
                    SubnetID: subnetmap[fmt.Sprintf("%s-GW-NET",labShortName)],
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
                },
            },
            HostParametersAttributes: []HostParametersAttributes {
                { Name: "role", Value: role },
            },
        },
    }

    // Convert data to json
    postData, err := json.Marshal(newHostStruct)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/hosts", url)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("POST", requesturl, bytes.NewBufferString(string(postData)))
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }

    // Make sure we are using the proper content type for the configs api
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json,version=2")

    // Set the session Cookie header
    req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", session))

    // Disable tls verify
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    // Set up the http client and do the request
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }
    defer resp.Body.Close()


    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 201 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return Host{}, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return Host{}, err
    }

    // Return the results
    return queryResults, nil
}


func GetHosts(url string, session string) ([]Host, error) {

    sysLogPrefix := "theforeman(package).hosts(file).GetHosts(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Get a list of hosts", sysLogPrefix))


    // Initialize the page and entries per page
    var entries_per_page int = 20
    var curpage int = 1
    var processed int = 0

    // Track current query and overall domain array
    hostList := []Host{}
    var queryResults HostQueryResults


    // Pager loop
    for {

        // Set the query url
        var requesturl string = fmt.Sprintf("%s/api/hosts?per_page=%d&page=%d", url, entries_per_page, curpage)

        // Set up the basic request from the url and body
        req, err := http.NewRequest("GET", requesturl, nil)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return hostList, err
        }

        // Make sure we are using the proper content type for the configs api
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json,version=2")

        // Set the session Cookie header
        req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", session))

        // Disable tls verify
        tr := &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }

        // Set up the http client and do the request
        client := &http.Client{Transport: tr}
        resp, err := client.Do(req)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return hostList, err
        }
        defer resp.Body.Close()

        // Read in the body and check status
        body, _ := ioutil.ReadAll(resp.Body)
        if resp.StatusCode != 200 {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
            return hostList, fmt.Errorf("%s %s", sysLogPrefix, string(body))
        }

        // Convert the body to a byte array
        bytes := []byte(body)

        // Unmarshall the json byte array into a struct
        err = json.Unmarshal(bytes, &queryResults)
        if err != nil {
            _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
            return hostList, err
        }

        // Append the results to the domain list
        hostList = append(hostList, queryResults.Results...)

        // Do the pager calculations. 
        processed += len(queryResults.Results)
        if (processed < queryResults.Subtotal) {
            curpage += 1
            continue
        } else {
            break
        }

    }

    return hostList, nil

}

func GetHostDetails(url string, session string, hostID interface{}) (Host, error) {

    sysLogPrefix := "theforeman(package).hosts(file).GetHostDetails(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Getting details for host %v", sysLogPrefix, hostID))

    // Init some vars
    hostInfo := Host{}

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/hosts/%v", url, hostID)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("GET", requesturl, nil)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return hostInfo, err
    }

    // Make sure we are using the proper content type for the configs api
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json,version=2")

    // Set the session Cookie header
    req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", session))

    // Disable tls verify
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    // Set up the http client and do the request
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return hostInfo, err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return hostInfo, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &hostInfo)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return hostInfo, err
    }

    return hostInfo, nil

}

func GetHostsDetailsByDomainID(url string, session string, domainID int) ([]Host, error) {

    sysLogPrefix := "theforeman(package).hosts(file).GetHostsDetailsByDomainID(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Get a list of host details filtered by domain id", sysLogPrefix))

    var hostDetailList []Host

    // Get standard host list
    hostList, err := GetHosts(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return []Host{}, err
    }

    // Parse through and get details
    for _, curHost := range hostList {
        if curHost.DomainID == domainID {

            curHostInfo, err := GetHostDetails(url, session, curHost.Name)
            if err != nil {
                _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
                return []Host{}, err
            }

            hostDetailList = append(hostDetailList, curHostInfo)
        }
    }

    return hostDetailList, nil

}

func DeleteHost(url string, session string, hostID interface{}) (error) {

    sysLogPrefix := "theforeman(package).hosts(file).DeleteHost(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Deleting host %v.", sysLogPrefix, hostID))

    // Set the query url assuming the key doesn't exist
    var requesturl string = fmt.Sprintf("%s/api/hosts/%v", url, hostID)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("DELETE", requesturl, nil)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }

    // Make sure we are using the proper content type for the configs api
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json,version=2")

    // Set the session Cookie header
    req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", session))

    // Disable tls verify
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    // Set up the http client and do the request
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return err
    }
    defer resp.Body.Close()

    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    return nil
}


func AddBridgeInterface(url string, session string, identifier string, hostName interface{}, domainID int, subnetID int) (NetInterfaceRet, error) {


    sysLogPrefix := "theforeman(package).hosts(file).AddBridgeInterface(func):"
    _ = sysLog.Debug(fmt.Sprintf("%s Creating %s bridge interface for host %s", sysLogPrefix, identifier, hostName))

    // Init the return value
    var queryResults NetInterfaceRet

    // Get a struct populated with common parameter data
    globalParams, err := GetGlobalParameters(url, session)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return NetInterfaceRet{}, err
    }

    // Convert the location name to id for host creation
    locationID, err := ConvLocNameToID(url, session, globalParams.LabLocationName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return NetInterfaceRet{}, err
    }
    fmt.Println(locationID)

    // Convert the organization name to id for host creation
    organizationID, err := ConvOrgNameToID(url, session, globalParams.LabOrgName)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
       return NetInterfaceRet{}, err
    }
    fmt.Println(organizationID)

    // Fill out the data structure to be used for creating the domain
    newInterfaceStruct := InterfacePostData{
//        OrganizationID: organizationID,
//        LocationID: locationID,
        Interface: NetInterface{
                    Identifier: identifier,
                    DomainID: domainID,
                    SubnetID: subnetID,
                    Primary: false,
                    Managed: false,
                    Provision: false,
                    Type: "bridge",
        },
    }
//        HostID: hostName.(string),

    // Convert data to json
    postData, err := json.Marshal(newInterfaceStruct)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return NetInterfaceRet{}, err
    }

    // Set the query url
    var requesturl string = fmt.Sprintf("%s/api/hosts/%v/interfaces", url, hostName)

    // Set up the basic request from the url and body
    req, err := http.NewRequest("POST", requesturl, bytes.NewBufferString(string(postData)))
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return NetInterfaceRet{}, err
    }

    // Make sure we are using the proper content type for the configs api
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json,version=2")

    // Set the session Cookie header
    req.Header.Set("Cookie", fmt.Sprintf("_session_id=%s", session))

    // Disable tls verify
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }

    // Set up the http client and do the request
    client := &http.Client{Transport: tr}
    resp, err := client.Do(req)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return NetInterfaceRet{}, err
    }
    defer resp.Body.Close()


    // Read in the body and check status
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 201 {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, string(body)))
        return NetInterfaceRet{}, fmt.Errorf("%s %s", sysLogPrefix, string(body))
    }

    // Convert the body to a byte array
    bytes := []byte(body)

    // Unmarshall the json byte array into a struct
    err = json.Unmarshal(bytes, &queryResults)
    if err != nil {
        _ = sysLog.Err(fmt.Sprintf("%s %s", sysLogPrefix, err))
        return NetInterfaceRet{}, err
    }

    return queryResults, nil

}

