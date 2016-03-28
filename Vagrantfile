Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"
  config.vm.network "public_network"

  config.vm.provision "shell", path: "vagrant_provision.sh"
end
