# Overview

Just testing out some automation for flexible lab environments. 


We are setting up labs in a mix of ironic nodes and vms all on a flat 'ironic' network.  Depending
on the configuration the networks will be vxlan or vlan backed.  Since we have a limited number of 
vlan blocks, we will need to set those up as static lab environments that can be "checked out" when
in use.  We can spin up as many dynamic vxlan backed environments as hardware allows.  These will
also allow checking out ironic nodes for ironic-on-ironic environments if ironic testing is needed.


## Setting up the environment

```
# Set up your python-openstackclient ~/.config/openstack/clouds.yaml with a 'phobos' cloud settings.

# Set up a local katello_venv with all the requirements
./scripts/setup_venv.sh

# Create the foreman server vm on the 'ironic' network.
./scripts/deploy_katello.sh

# Do some prep work for foreman and openstack
./scripts/prep_lab_environment.sh

```


## Building out a lab
