#!/bin/bash

# Source the virtual env
. ./katello_venv/bin/activate

# Install the ansible modules for foreman
ansible-galaxy collection install theforeman.foreman


# Run the vm setup playboks
if [ ! -e .prep_lab_environment ]; then
  pushd playbooks
  stdbuf -i0 -e0 -o0 ansible-playbook -i inventory/katello_vms.ini prep_lab_environment.yml
  if [ $? == 0 ]; then
    touch ../.prep_lab_environment
  else
    echo "prep_lab_environment.yml playbook failed"
    exit 255
  fi
  popd
fi

# Drop out of the venv
deactivate

