#!/bin/bash

# Source the virtual env
. ./katello_venv/bin/activate

# Run the destroy playbooks
pushd playbooks
ansible-playbook destroy_katello_vm.yml
popd

deactivate

# Remove the forklift repo
if [ -d "./forklift" ]; then
  rm -rf ./forklift
fi

# Remove any step markers
if [ -e ".setup_katello_complete" ]; then
  rm .setup_katello_complete 
fi

if [ -e ".setup_forklift_complete" ]; then
  rm .setup_forklift_complete 
fi

