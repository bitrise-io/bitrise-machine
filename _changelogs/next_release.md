## Changes

* New cleanup mode: `CleanupModeDestroy` / `destroy`.
  Using this option the VM won't be re-created, only destroyed, you'll have to
  call `setup` before the next `run`. This can be useful if you want to
  set different parameters (environment variables) for every setup/cleanup,
  for example the template to use to create the Virtual Machine.
* You can now allow the creation of the Virtual Machine in `setup`,
  by setting the `is_allow_vagrant_create_in_setup` option to `true` in `bitrise.machine.config.json`.
  If this option is set to `false` (default), then, just like it was before,
  only `cleanup` is allowed to create the Virtual Machine, (in case of
  a `recreate` and `custom` cleanups), `setup` is not (unless it does a cleanup).
    * You don't need to set this option for the new `destroy` cleanup mode,
      it's explicitly allowed for `destroy` cleanup mode to create the
      Virtual Machine in `setup`.
    * It's important to note that this creation in `setup` will happen **after** the cleanup,
      in case the `is_cleanup_before_setup` option is set to `true`. This is
      required to prevent the Virtual Machine to be destroyed right away,
      for example if the cleanup mode is `destroy`.
* To support the dynamic creation use-case a new flag is now available: `-e` / `--environment`.
  With `-e MY_KEY=my-value` you can specify custom environment variable(s)
  for your commands. You can add multiple `-e` flags (e.g. `-e KEY1=val1 -e KEY2=val2`).
  The environment variables you define with `-e` will be appended to the
  ones defined in your `bitrise.machine.config.json` config file - this
  also means that you can overwrite the ones defined in `bitrise.machine.config.json`.
  These environment variables will be available for `setup`, `cleanup` and `destroy`
  commands, just like the `envs` defined in `bitrise.machine.config.json`.
