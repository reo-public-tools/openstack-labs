

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
    pushd forklift > /dev/null
    stdbuf -i0 -e0 -o0 ansible-playbook \
          -i ./inventories/localhost \
          -e @../playbooks/vars/forklift_mnaio.yml \
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
    popd > /dev/null

    # dhcp and dns need enabled.  the installer sets it up http for some reason.
    FPRESTART=0
    for PRSV in dhcp dns; do
      if [ -e "/etc/foreman-proxy/settings.d/${PRSV}.yml" ]; then
        grep '^:enabled: true' /etc/foreman-proxy/settings.d/${PRSV}.yml > /dev/null 2>&1
        if [ $? != 0 ]; then
          sed -i -e 's/^:enabled:.*/:enabled: true/g' /etc/foreman-proxy/settings.d/${PRSV}.yml
          FPRESTART=1
        fi
      fi
    done
    if [ $FPRESTART == 1 ]; then
       systemctl restart foreman-proxy
    fi
 

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
  
    # Get the primary ip for ssh testing
    CURIP=$(ip addr show dev $(ip r | awk '/default/{print $5}') | awk '/ inet /{print $2}' | awk -F '/' '{print $1}')

    # test it out 
    su - foreman -s /bin/bash -c "ssh -o StrictHostKeyChecking=no root@${CURIP} id"
    if [ $? != 0 ]; then
      echo "The foreman user is unable to ssh in as root with the ssh key"
      exit 255
    fi

    # Set up foreman-proxy ssh key to use for libvirt connectivity.
    if [ ! -e /usr/share/foreman-proxy/.ssh/id_rsa.pub ]; then
      chown foreman-proxy: /usr/share/foreman-proxy
      su foreman-proxy -s /bin/bash -c "ssh-keygen -f /usr/share/foreman-proxy/.ssh/id_rsa -q -N ''"
      restorecon -RvF /usr/share/foreman-proxy/.ssh
    fi

    # Make sure foreman-proxy can ssh to root for libvirt connectivity
    grep 'foreman-proxy@' /root/.ssh/authorized_keys
    if [ $? != 0 ]; then
      cat /usr/share/foreman-proxy/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys
    fi
  
    # test it out 
    su - foreman-proxy -s /bin/bash -c "ssh -o StrictHostKeyChecking=no root@${CURIP} id"
    if [ $? != 0 ]; then
      echo "The foreman-proxy user is unable to ssh in as root with the ssh key"
      exit 255
    fi

  fi

  # Make sure the foreman proxy has the needed sudo setup
  if [ ! -e /etc/sudoers.d/foreman-proxy ]; then
    echo 'Defaults !requiretty' > /etc/sudoers.d/foreman-proxy
    echo 'foreman-proxy ALL = NOPASSWD : /usr/bin/virsh' >> /etc/sudoers.d/foreman-proxy
    chmod 0440 /etc/sudoers.d/foreman-proxy
  fi
}


function install_ansible_module_foreman {

  # Run the install if not already completed
  pushd playbooks > /dev/null
  if [ ! -e ../.install_ansible_module_foreman ]; then
    ansible-galaxy collection install theforeman.foreman
    touch ../.install_ansible_module_foreman
  fi
  popd > /dev/null
}


function go_compile {

  # Make sure golang is installed
  rpm -qi golang > /dev/null 2>&1
  if [ $? != 0 ]; then
    yum -y install golang
  fi

  pushd go > /dev/null
  . sourceme
  make all
  if [ $? == 0 ]; then
    touch ../.go_compile
  else
    echo "Go compile failed"
    exit 255
  fi
  popd go > /dev/null
}
