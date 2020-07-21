#!/bin/bash

# Run this from the parent of the scripts dir.

##############
# Variables
##############

# OSA or OSP
export MNAIO_TYPE=${MNAIO_TYPE:-OSA}
export ENVADD="-e @./vars/lab_prep_mnaio.yml"

# Set some default versions and pull in deploy took specific vars
if [ $MNAIO_TYPE == OSA ]; then
  export MNAIO_VERSION=${MNAIO_VERSION:-stable/train}
  export ENVADD="${ENVADD} -e @./vars/mnaio_osa.yml -e @./vars/forklift.yml -e version_match=${MNAIO_VERSION}"
else
  export MNAIO_VERSION=${MNAIO_VERSION:-OSP13}
  export ENVADD="${ENVADD} -e @./vars/mnaio_osp.yml -e @./vars/forklift.yml -e version_match=${MNAIO_VERSION}"
fi


#######
# Prep
#######

# Install apypie for python2 > /dev/null 2>&1
rpm -qi python2-pip
if [ $? != 0 ]; then
   echo "Installing python2-pip for apypie"
   yum -y install python2-pip
fi
pip install apypie

# Install rubygen-ruby-libvirt for libvirt dhcp smart proxy
rpm -qi rubygem-ruby-libvirt > /dev/null 2>&1
if [ $? != 0 ]; then
   echo "Installing rubygen-ruby-libvirt package for libvirt dhcp smart proxy"
   yum -y install rubygem-ruby-libvirt
fi

# Run the script to set up the venv if not done already
if [ ! -d ./katello_venv ]; then
  ./scripts/setup_venv.sh
fi

# Source the functions file
. scripts/functions.sh

# Deploy katello using forklift
deploy_katello_mnaio

# Prep libvirt for theforman compute resource
. katello_venv/bin/activate



######################################
# Install the foreman-ansible-modules
######################################
install_ansible_module_foreman

###############################
# Run the mnaio setup playbook
###############################

rm ./.setup_mnaio
# Run the setup_mnaio.yml playbook
if [ ! -e .setup_mnaio ]; then
  pushd playbooks

  #echo ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml | tee setup_mnaio.log
  #stdbuf -i0 -e0 -o0 ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml | tee setup_mnaio.log
  #echo ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml --tags create_organization
  #stdbuf -i0 -e0 -o0 ../katello_venv/bin/ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml --tags setup_libvirtd,create_organization,create_location,create_users,create_domains,create_subnets,create_operatingsystems,update_foreman_settings,update_foreman_global_parameters #,create_libvirt_resources,provisioning_templates
  #stdbuf -i0 -e0 -o0 ../katello_venv/bin/ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml --tags create_libvirt_resources
  ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml --tags create_operatingsystems,provisioning_templates,create_libvirt_resources,create_compute_profiles,create_hostgroups,provision_environments,foreman_hooks
  #../katello_venv/bin/ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml --tags setup_libvirtd,create_organization,create_location,create_users,create_domains,create_subnets,create_operatingsystems,update_foreman_settings,update_foreman_global_parameters #,create_libvirt_resources,provisioning_templates

  if [ $? == 0 ]; then
    touch ../.setup_mnaio
  else
    echo "setup_mnaio.yml playbook failed"
    exit 255
  fi
  popd
fi




# Deactivate the venv
deactivate

