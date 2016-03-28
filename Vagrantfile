Vagrant.configure(2) do |config|
  config.vm.box = "bento/debian-8.2"
  config.vm.network "public_network", ip: "192.168.1.40"

  config.vm.provision "shell", path: "vagrant_provision.sh"
end
