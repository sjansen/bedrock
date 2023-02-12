# -*- mode: ruby -*-
# vi: set ft=ruby :
require 'json'

boxes = JSON.parse(
  File.read("vagrant/boxes.json"), {symbolize_names: true}
)
is_arm64 = if Vagrant::Util::Platform.darwin?
  `sysctl -n hw.optional.arm64`.strip() == "1"
end
ENV["LC_ALL"] = "en_US.UTF-8"

unless Vagrant.has_plugin?("vagrant-hosts")
  puts("ERROR: vagrant-hosts plugin is missing")
  puts("$ vagrant plugin install vagrant-hosts")
  exit(1)
end

Vagrant.configure("2") do |config|
  config.ssh.forward_agent = true

  #config.vm.boot_timeout = 300
  config.vm.box = if is_arm64
    "bento/ubuntu-22.04-arm64"
  else
    "bento/ubuntu-22.04"
  end
  config.vm.box_check_update = false

  # vagrant plugin install vagrant-parallels
  config.vm.provider "parallels" do |vm, override|
    vm.customize ["set", :id, "--time-sync", "on"]
    vm.customize ["set", :id, "--disable-timezone-sync", "on"]
    vm.update_guest_tools = true
  end

  if File.exists?("Vagrantfile.local")
    eval File.read "Vagrantfile.local"
  end

  boxes.each do |name, opts|
    config.vm.define name do |config|
      config.vm.hostname = name
      config.vm.network "private_network", ip: opts[:eth1], :netmask => "255.255.255.0"
      config.vm.provision "hosts", sync_hosts: true
      config.vm.provision "prepare-vm", type: "shell", privileged: true, path: "vagrant/scripts/prepare-vm"

      if opts[:controller]
        config.vm.provision "install-ansible", type: "shell", privileged: true, path: "vagrant/scripts/install-ansible"
        config.vm.provision "update-ansible-hosts", type: "shell", privileged: true, path: "vagrant/scripts/update-ansible-hosts", run: "always"
      end

      if opts.has_key?(:forward)
        opts[:forward].each do |port|
          config.vm.network "forwarded_port", guest: port, host: port, host_ip: "127.0.0.1"
        end
      end

      config.vm.provider "parallels" do |vm, override|
        vm.name = "bedrock-#{config.vm.hostname}"
        vm.cpus = opts[:cpus]
        vm.memory = opts[:memory]
      end
    end
  end
end
