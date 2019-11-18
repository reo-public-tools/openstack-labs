#!/bin/bash

# Source the virtual env
. ./katello_venv/bin/activate


# Run the vm setup playboks
if [ ! -e .setup_katello_complete ]; then
  pushd playbooks
  stdbuf -i0 -e0 -o0 ansible-playbook setup_katello_vm.yml
  if [ $? == 0 ]; then
    touch ../.setup_katello_complete
  else
    echo "setup_katello_vm.yml playbook failed"
    exit 255
  fi
  popd
fi


# Run the forklift playbook
if [ ! -e .setup_forklift_complete ]; then
  if [ ! -e "./forklift" ]; then
    git clone https://github.com/theforeman/forklift.git
  fi
  pushd forklift
  stdbuf -i0 -e0 -o0 ansible-playbook \
  	-l katello_vms \
  	-i ../playbooks/inventory/katello_vms.ini \
        -e @../playbooks/vars/forklift.yml \
  	playbooks/katello.yml
  if [ $? == 0 ]; then
    touch ../.setup_forklift_complete
  else
    echo "The forklift playbooks failed."
    exit 255
  fi
  popd
fi

# Drop out of the venv
deactivate
