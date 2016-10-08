# Bitrise Machine - manage bitrise CLI runner hosts

## Requirements

* [vagrant](https://www.vagrantup.com)
* [bitrise-bridge](https://github.com/bitrise-io/bitrise-bridge) for communication with the host


## Cleanup modes

* `rollback` : runs `vagrant snapshot pop` to clean up - **requires** at least `vagrant` v1.8.0
* `recreate` : runs `vagrant destroy -f` and then `vagrant up` to clean up
* `destroy` : runs `vagrant destroy -f` to cleanup, and allows `bitrise-machine setup` to create the Virtual Machine with `vagrant up`
* `custom-command` : runs `vagrant CUSTOM-COMMAND` to clean up, or `vagrant up` in case the Virtual Machine is not yet created
  * useful for provider plugins which add custom `vagrant` actions, which can be used for cleanup.
    For example the [vagrant-digitalocean](https://github.com/smdahlen/vagrant-digitalocean) plugin adds
    a `rebuild` command to `vagrant` and makes a cleanup / re-build faster than
    re-creating the Virtual Machine with `vagrant destroy` and `vagrant up`.
