

# Function to deploy katello/theforeman with forklift
function deploy_katello_mnaio {

  # Add system entry into /etc/hosts if it doesn't exist
  grep "$(ip route get 1 | awk '{print $NF;exit}') $(hostname)" /etc/hosts >/dev/null 2>&1
  if [ $? != 0 ]; then
    echo "$(ip route get 1 | awk '{print $NF;exit}') $(hostname)" >> /etc/hosts
  fi

  # Run the install if not already completed
  if [ ! -e .setup_forklift_complete ]; then
    if [ ! -e "./forklift" ]; then
      git clone https://github.com/theforeman/forklift.git
    fi
    pushd forklift
    stdbuf -i0 -e0 -o0 ansible-playbook \
          -i ./inventories/localhost \
          -e @../playbooks/vars/forklift.yml \
          playbooks/katello.yml | tee forklift_install.log
    if [ $? == 0 ]; then
      # Touch file to skip on next run
      touch ../.setup_forklift_complete

      # Add the local device as a host
      /opt/puppetlabs/puppet/bin/puppet agent --test
    else
      echo "The forklift playbooks failed."
      exit 255
    fi
    popd

    # Set up foreman ssh key to use for libvirt connectivity.
    if [ ! -e /usr/share/foreman/.ssh/id_rsa.pub ]; then
      su foreman -s /bin/bash -c "ssh-keygen -f /usr/share/foreman/.ssh/id_rsa -q -N ''"
      restorecon -RvF /usr/share/foreman/.ssh
    fi

    # Make sure foreman can ssh to root for libvirt connectivity
    grep 'foreman@' /root/.ssh/authorized_keys
    if [ $? != 0 ]; then
      cat /usr/share/foreman/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys
    fi
  
    # test it out 
    su - foreman -s /bin/bash -c 'ssh -o StrictHostKeyChecking=no root@localhost id'
    if [ $? != 0 ]; then
      echo "The foreman user is unable to ssh in as root with the ssh key"
      exit 255
    fi

  fi
}


function install_ansible_module_foreman {

  # Run the install if not already completed
  pushd playbooks
  if [ ! -e ../.install_ansible_module_foreman ]; then
    ansible-galaxy collection install theforeman.foreman
    touch ../.install_ansible_module_foreman
  fi
  popd
}
