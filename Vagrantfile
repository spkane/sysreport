# -*- mode: ruby -*-
# vi: set ft=ruby :

FileUtils.chdir(File.dirname(__FILE__))

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

$script = <<SCRIPT
sudo apt-get -y update
sudo apt-get -y install git curl
sudo cp -f /home/vagrant/go/tools/init.ubuntu /etc/init.d/sysreport
sudo chmod +x /etc/init.d/sysreport
sudo cp -f /home/vagrant/go/tools/fake-rpm.sh /usr/bin/rpm
sudo chmod +x /usr/bin/rpm
sudo update-rc.d sysreport defaults
if [ ! -f /home/vagrant/.goenv ]; then
  git clone https://github.com/wfarr/goenv.git ~/.goenv
  export PATH="$HOME/.goenv/bin:$PATH"
  echo 'export PATH="$HOME/.goenv/bin:$PATH"' >> ~/.bashrc
  eval "$(goenv init -)"
  echo 'eval "$(goenv init -)"' >> ~/.bashrc
  goenv install 1.2
  goenv global 1.2
  goenv rehash
  go get github.com/spkane/go-utils/debugtools
  go get github.com/spkane/go-utils/jsonutils
  go get github.com/spkane/go-utils/strutils
  mkdir -p /home/vagrant/gopath
  export GOPATH="/home/vagrant/gopath"
  go get github.com/spkane/go-utils/debugtools
  go get github.com/spkane/go-utils/jsonutils
  go get github.com/spkane/go-utils/strutils
  echo 'export GOPATH="/home/vagrant/gopath"' >> ~/.bashrc
  echo 'export GOROOT=`go env GOROOT`' >> ~/.bashrc
  echo 'export GOBIN="$GOPATH/bin"' >> ~/.bashrc
  echo '#export GOARCH=amd64' >> ~/.bashrc
  echo '#export GOOS=darwin' >> ~/.bashrc
  echo '#export CGO_ENABLED=1' >> ~/.bashrc
  echo 'export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"' >> ~/.bashrc
  sudo service sysreport start
fi
echo "ALL DONE"
echo ""
echo ""
SCRIPT

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  # All Vagrant configuration is done here. The most common configuration
  # options are documented and commented below. For a complete reference,
  # please see the online documentation at vagrantup.com.

  config.vm.box = "precise64_vmware_fusion"

  config.vm.synced_folder ".", "/home/vagrant/go"

  config.ssh.forward_agent = true

  config.vm.provider :virtualbox do |virtualbox, override|
    override.vm.box_url = "http://files.vagrantup.com/precise64.box"
    override.vm.host_name = 'test-target-1'
    override.vm.network :private_network, type: :static, ip: "192.168.50.34"
  end

  config.vm.provider :vmware_fusion do |vmware_fusion, override|
    override.vm.box_url = "http://files.vagrantup.com/precise64_vmware_fusion.box"
    override.vm.host_name = 'test-target-1'
    override.vm.network :private_network, type: :static, ip: "192.168.50.34"
  end

  config.vm.provision "shell", inline: $script, :privileged => false

end
