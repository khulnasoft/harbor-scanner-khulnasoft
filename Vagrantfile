# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/focal64"

  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.memory = "2096"
  end

  config.vm.provision "shell", env: {
    "KHULNASOFT_REGISTRY_USERNAME" => ENV['KHULNASOFT_REGISTRY_USERNAME'],
    "KHULNASOFT_REGISTRY_PASSWORD" => ENV['KHULNASOFT_REGISTRY_PASSWORD']
    }, inline: <<-SHELL
    if [ -z "$KHULNASOFT_REGISTRY_USERNAME" ]; then echo "KHULNASOFT_REGISTRY_USERNAME env is unset" && exit 1; fi
    if [ -z "$KHULNASOFT_REGISTRY_PASSWORD" ]; then echo "KHULNASOFT_REGISTRY_PASSWORD env is unset" && exit 1; fi
  SHELL

  config.vm.provision "install-go", type: "shell", path: "vagrant/install-go.sh"
  config.vm.provision "install-docker-ce", type: "shell", path: "vagrant/install-docker.sh"
  config.vm.provision "install-harbor", type: "shell", path: "vagrant/install-harbor.sh", env: {
    "HARBOR_VERSION" => ENV["HARBOR_VERSION"] || "v2.4.0"
  }
  config.vm.provision "install-khulnasoft", type: "shell", path: "vagrant/install-khulnasoft.sh", env: {
    "KHULNASOFT_REGISTRY_USERNAME" => ENV['KHULNASOFT_REGISTRY_USERNAME'],
    "KHULNASOFT_REGISTRY_PASSWORD" => ENV['KHULNASOFT_REGISTRY_PASSWORD'],
    "KHULNASOFT_VERSION" => ENV['KHULNASOFT_VERSION'] || "6.5",
    "HARBOR_VERSION" => ENV['HARBOR_VERSION'] || "v2.4.0"
  }

  # Access Harbor Portal at http://localhost:8181 (admin/@Harbor12345)
  config.vm.network :forwarded_port, guest: 8080, host: 8181

  # Access Khulnasoft Management Console at http://localhost:9181 (administrator/@Khulnasoft12345)
  config.vm.network :forwarded_port, guest: 9080, host: 9181
end
