# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  # Use the bento version because it's considered friendlier for multiple
  # providers compared to the "official" ubuntu version
  config.vm.box = "bento/ubuntu-20.04"

  config.vm.define "cyna", primary: true do |cyna|
    cyna.vm.hostname = "cyna"
    cyna.vm.network "private_network", ip: "192.168.58.2"
  end

  config.vm.define "cynb" do |cynb|
    cynb.vm.hostname = "cynb"
    cynb.vm.network "private_network", ip: "192.168.58.3"
  end

  # If we really need it, uncomment this...
  #config.vm.define "cync" do |cync|
    #cync.vm.hostname = "cync"
    #cync.vm.network "private_network", ip: "192.168.58.4"
  #end
end
