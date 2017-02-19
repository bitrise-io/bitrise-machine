# Bitrise Machine - manage bitrise CLI runner hosts

## Requirements

- [vagrant](https://www.vagrantup.com)
- [bitrise-bridge](https://github.com/bitrise-io/bitrise-bridge) for communication with the host

## Cleanup modes

__IMPORTANT:__ if you want to change a config's Cleanup Mode, first make sure that
the machine/VM is not running, e.g. by running `bitrise-machine destroy`. Only after that
you should change the Cleanup Mode in the config!

- `rollback` : runs `vagrant snapshot pop` to clean up
    - __requires__ at least `vagrant` v1.8.0
    - SESSION: partial session support, initializes the session when the VM is created,
      it does not end or start a new session when a simple cleanup/rollback happens,
      only when the VM is actually re-created (e.g. after a destroy)
- `recreate` : runs `vagrant destroy -f` and then `vagrant up` to clean up
    - `bitrise-machine cleanup` is the same as `bitrise-machine destroy && bitrise-machine setup`
    - SESSION: full session handling support, init session when VM created
- `destroy` : runs `vagrant destroy -f` to cleanup, and allows/requires `bitrise-machine setup` (to create the Virtual Machine with `vagrant up`)
    - useful for on-demand conigurations, where you might want to have periods when the virtual machine/host is not created
    - `bitrise-machine cleanup` only cleans up / destroys the virtual machine, it does not recreate it.
      You have to run `bitrise-machine setup` to create a new one after a cleanup.
    - SESSION: full session handling support, init session when VM created
- `custom-command` : runs `vagrant CUSTOM-COMMAND` to clean up, or `vagrant up` in case the Virtual Machine is not yet created
    - useful for provider plugins which add custom `vagrant` actions, which can be used for cleanup.
      For example the [vagrant-digitalocean](https://github.com/smdahlen/vagrant-digitalocean) plugin adds
      a `rebuild` command to `vagrant` and makes a cleanup / re-build faster than
      re-creating the Virtual Machine with `vagrant destroy` and `vagrant up`.
    - SESSION: partial session support, initializes the session when the VM is created,
      it does not end or start a new session when a simple cleanup/rollback happens,
      only when the VM is actually re-created (e.g. after a destroy)

## Session

Bitrise Machine exposes a "session time id" as an environment variable,
which can be used in the Vagrantfile.

Where session is fully supported, this "session time id" persists between
a setup and a destroy, so the session time id will be the same
during `bitrise-machine setup` and the following `bitrise-machine destroy/cleanup`,
and will be re-generated at the next "vagrant up".

The environment variable is: `BITRISE_MACHINE_SESSION_TIME_ID`

Format: `YYYYMMDDHHMMSS`
Example: `20170215093215`


## TODO

- complete restructure: make everything session based
    - e.g. session.Start should create the SSH keys too
    - and session.End should delete them
- use Interface for "lifecycle handlers", one for each cleanup mode,
  so the common code just determines which handler to use (based on the cleanup mode),
  and then calls the handler's Setup, Destroy, Cleanup, ... methods,
  instead of doing this logic inline in every function (e.g. `doCleanup`)
