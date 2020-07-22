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
rpm -qi python2-pip > /dev/null 2>&1
if [ $? != 0 ]; then
  echo "Installing python2-pip for apypie"
  yum -y install python2-pip
fi
pip show apypie > /dev/null 2>&1
if [ $? != 0 ]; then
  pip install apypie
fi

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

# Run the setup_mnaio.yml playbook
if [ ! -e .setup_mnaio ]; then
  pushd playbooks > /dev/null

  echo ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml | tee setup_mnaio.log
  stdbuf -i0 -e0 -o0 ansible-playbook -i ',localhost' ${ENVADD} setup_mnaio.yml | tee setup_mnaio.log

  if [ $? == 0 ]; then
    touch ../.setup_mnaio
  else
    echo "setup_mnaio.yml playbook failed"
    exit 255
  fi
  popd > /dev/null
fi


# Deactivate the venv
deactivate


########################
# Compile any go code
########################

install_oslabs
