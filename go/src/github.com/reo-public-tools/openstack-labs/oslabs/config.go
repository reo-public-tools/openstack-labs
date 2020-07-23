package main

import (
    "os"
    "io/ioutil"
    "gopkg.in/yaml.v2"
)

type LabConfig struct {
        ForemanURL          string `yaml:"foreman_url"`
        ForemanArchitecture string `yaml:"foreman_architecture"`
	LabType             string `yaml:"lab_type"`
	OpenstackCloud      string `yaml:"openstack_cloud"`
	PubSSHKey           string `yaml:"pub_ssh_key"`
        OpenstackNetwork    string `yaml:"openstack_network"`
	Environment         struct {
                  ForemanHosts        []ForemanHosts        `yaml:"foreman_hosts"`
                  IronicOnIronicHosts []IronicOnIronicHosts `yaml:"ironic_on_ironic_hosts"`
                  SimulatedHosts      []SimulatedHosts      `yaml:"simulated_hosts"`
	} `yaml:"environment"`
}

type ForemanHosts struct {
        Role                  string `yaml:"role"`
        ForemanHostgroup      string `yaml:"foreman-hostgroup"`
        Count                 int    `yaml:"count"`
        Name                  string `yaml:"name,omitempty"`
        SimulatorType         string `yaml:"simulator-type,omitempty"`
}

type IronicOnIronicHosts struct {
        Flavor string `yaml:"flavor"`
        Count  int    `yaml:"count"`
}

type SimulatedHosts struct {
        Role                  string `yaml:"role"`
        ForemanHostgroup      string `yaml:"foreman-hostgroup"`
        Count                 int    `yaml:"count"`
        SimulatorName         string `yaml:"simulator-name"`
}


func PopulateConfigFile(configFile string) (LabConfig, error) {

    // Initialize the config
    returnConfig := LabConfig{}

    // Read in global config first if it exists
    if _, err := os.Stat(globalConfigFile); ! os.IsNotExist(err) {

        // Read in the config
        globalData, err := ioutil.ReadFile(globalConfigFile)
        if err != nil {
            return LabConfig{}, err
        }

        // Pull yaml into LabConfig struct
        err = yaml.Unmarshal(globalData, &returnConfig)
        if err != nil {
            return LabConfig{}, err
        }

    }

    // Read in the config(only if not empty)
    if configFile != "" {
        data, err := ioutil.ReadFile(configFile)
        if err != nil {
            return LabConfig{}, err
        }

        // Pull yaml into LabConfig struct
        err = yaml.Unmarshal(data, &returnConfig)
        if err != nil {
            return LabConfig{}, err
        }
    }

    return returnConfig, nil

}



