#!/bin/bash

yum -y groupinstall "Development Tools"
yum -y install libvirt-devel python3-devel

VENV_LOC="./katello_venv"

# Create the venv
if [ ! -d $VENV_LOC ]; then
  python3 -m venv $VENV_LOC
fi

# Source and make sure any dependancies are installed
. ${VENV_LOC}/bin/activate
pip install --upgrade pip
pip install -r pip-requirements.txt
deactivate
