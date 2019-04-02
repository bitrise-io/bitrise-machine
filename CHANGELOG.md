# Changelog

-----------------

## 0.9.14 (2017 Feb 20)

### Release Notes

* Fixing a session handling bug for the `recreate` cleanup mode.
  The bug does not affect you if you don't use the session time ID environment variable.


### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.14/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.13 -> 0.9.14

* [4496830] Viktor Benei - session handling fix for recreate cleanup mode (#11) (2017 Feb 20)


## 0.9.13 (2017 Feb 19)

### Release Notes

Session handling for `destroy` and `recreate` cleanup modes.

TL;DR;

Bitrise Machine exposes a "session time id" as an environment variable,
which can be used in the Vagrantfile.

This session time id can be used to e.g. include it in the Virtual Machine ID,
to help with unique ID generation, as the ID will be kept from "vagrant up"
to "vagrant destroy", and the next "vagrant up" will generate a new session (time id).

_For more info see the README._

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.13/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.12 -> 0.9.13

* [ee305bd] Viktor Benei - Feature/session handling (#10) (2017 Feb 19)


## 0.9.12 (2017 Feb 05)

### Release Notes

* Max log buffer size doubled
* A couple of logging changes / revs

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.12/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.11 -> 0.9.12

* [e8b4cff] Viktor Benei - 0.9.12 (2017 Feb 05)
* [5054e37] Viktor Benei - bit of additional logging with time stamps (#9) (2017 Feb 05)
* [c6b4d55] Viktor Benei - Feature/log buffer max size bump and lot rev (#8) (2017 Feb 05)


## 0.9.11 (2017 Feb 04)

### Release Notes

* Log buffer max size, and proper "overflow" handling (#7)

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.11/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.10 -> 0.9.11

* [d4653e6] Viktor Benei - 0.9.11 (2017 Feb 04)
* [4ed94ea] Viktor Benei - Log buffer max size, and proper "overflow" handling (#7) (2017 Feb 04)
* [ba335c6] Viktor Benei - bitrise.yml normalized by workflow editor (2017 Feb 03)
* [081d8c1] Viktor Benei - go deps update (#6) (2017 Feb 03)


## 0.9.10 (2016 Nov 16)

### Release Notes

* dependency updates & recompile with Go 1.7.3
* a new `version` command, with `--full` flag: `bitrise-machine version --full`
    * prints a verbose version info, including the Go version, OS and ARCH
      where the binary was created

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.10/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.9 -> 0.9.10

* [043bec4] Viktor Benei - Feature/dep updates (#5) (2016 Nov 16)


## 0.9.9 (2016 Oct 08)

### Release Notes

* __BREAKING__ : the `rollback` cleanup mode now uses `vagrant`'s built in
  `vagrant snapshot` command (`vagrant snapshot pop --no-delete`), instead
  of the `vagrant-sahara` plugin. This also means that it requires `vagrant` v1.8.0
  or newer

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.9/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.8 -> 0.9.9

* [f68a96b] Viktor Benei - v0.9.9 (2016 Oct 08)
* [c1463cf] Viktor Benei - godeps update + bitrise.yml deps update (#4) (2016 Oct 08)
* [05536f0] Viktor Benei - Merge pull request #3 from bitrise-tools/feature/vagrant-snapshot-instead-of-sandbox (2016 Oct 08)
* [3a2292a] Viktor Benei - replaced vagrant sandbox with snapshot (2016 Oct 08)
* [2f99c59] Viktor Benei - minor bitrise.yml revision (2016 Aug 01)
* [1ae18cb] Viktor Benei - comment fix (2016 May 31)


## 0.9.8 (2016 May 07)

### Release Notes

* __NEW__ : Config Type Envs can now be defined, in addition
  to "generic" envs in bitrise-machine config.
  By default these are not loaded / used, but you can define which
  Config Type Envs you want to "activate" (add to the base Envs)
  with the new `-config-type-id` flag.

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.8/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.7 -> 0.9.8

* [a2cf1a0] Viktor Benei - Docker : Go - use 1.6.2 ; bitrise CLI 1.3.3 (2016 May 07)
* [1fa37e6] Viktor Benei - v0.9.8 (2016 May 07)
* [1d28643] Viktor Benei - LOG: MachineConfigTypeID - only debug print (2016 May 07)
* [cabd38b] Viktor Benei - Config Type Envs can now be defined, in addition to "generic" envs in bitrise-machine config. (2016 May 07)
* [0e080b1] Viktor Benei - bitrise.yml : script step minor (logging) updates (2016 May 06)
* [75fe394] Viktor Benei - Godeps update (2016 May 06)
* [5178de4] Viktor Benei - format_version upgrade (2016 May 06)
* [5991a9d] Viktor Benei - stingified (2016 May 06)
* [045c75a] Viktor Benei - go vet doesn't have to be installed anymore, it's part of the Go toolkit now (2016 May 06)
* [cb8ea64] Viktor Benei - no need to install Go 1.6.0 anymore - it's preinstalled everywhere now (2016 May 06)


## 0.9.7 (2016 Mar 05)

### Release Notes

* Go CPU Profiling can now be enabled by setting the `BITRISE_MACHINE_CPU_PROFILE_FILEPATH`
  Environment Variable to a file path, e.g. `export BITRISE_MACHINE_CPU_PROFILE_FILEPATH=./cpu.profile`.
  This profile can be used directly with Go's [`pprof` command line tool](http://blog.golang.org/profiling-go-programs).
* **Optimized Log handling**: less frequent log chunk processing (instead of doing it every 100ms the tick is now 500ms),
  and optimized log chunk buffer handling. **These changes should significantly reduce the CPU usage**,
  in a typical, sustained use case the difference can be 5-10x less CPU usage (~3% CPU usage
  where the previous version was around 20-25% on the same machine).

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.7/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.6 -> 0.9.7

* [bd2fc91] Viktor Benei - Initialized new Changelog, using `releaseman` - full migration is work-in-progress (2016 Mar 01)
* [68f3f80] Viktor Benei - Changelog: 0.9.7 (2016 Mar 01)
* [9d4c647] Viktor Benei - on CI make sure Go 1.6 is installed (2016 Mar 01)
* [ccdfbbd] Viktor Benei - Upgrades for Go 1.6 vendor folder - for testing (2016 Mar 01)
* [5bc5fca] Viktor Benei - Godeps.json (2016 Mar 01)
* [16c1e85] Viktor Benei - Docker image: use Go 1.6 (2016 Mar 01)
* [0a78829] Viktor Benei - deps update -> moved into vendor folder (2016 Mar 01)
* [7345dc3] Viktor Benei - v0.9.7 (2016 Mar 01)
* [d98dd36] Viktor Benei - Optimization: "tick" functions frequency reduced from 100 milliseconds to 500 milliseconds (2016 Mar 01)
* [3264036] Viktor Benei - optimization in ReadRunes function, as it was the main CPU bottleneck (2016 Mar 01)
* [dee129f] Viktor Benei - Go CPU profiling can now be turned-on with the `BITRISE_MACHINE_CPU_PROFILE_FILEPATH` Environment Variable (2016 Mar 01)
* [645472e] Viktor Benei - Release URL change (2016 Feb 29)


## 0.9.6 (2016 Jan 12)

### Release Notes

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

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.6/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.5 -> 0.9.6

* [f50178c] Viktor Benei - changelog for 0.9.6 (2016 Jan 12)
* [f4f2f9a] Viktor Benei - err check fix (2016 Jan 12)
* [47f2019] Viktor Benei - implemented CleanupMode destroy ; and global CLI params are now stored in Freezable objects, to prevent modification (2016 Jan 12)
* [440d333] Viktor Benei - godeps-update (2016 Jan 12)
* [f6eff19] Viktor Benei - finishing the command line Additional Env Vars handling (2016 Jan 11)
* [33a3602] Viktor Benei - removed the install_bitrise script (2016 Jan 11)
* [895c1ab] Viktor Benei - FIX : machine_config_test (2016 Jan 11)
* [6663d0a] Viktor Benei - godeps update (2016 Jan 11)
* [41ada6b] Viktor Benei - new flag: `--environment` / `-e` (2016 Jan 11)
* [ec25654] Viktor Benei - upgraded bitrise CLI version (2016 Jan 06)


## 0.9.5 (2015 Dec 31)

### Release Notes

* Added support for `vagrant` v1.8.x 's new `vagrant ssh-config` output format,
  which now wraps the IdentityPath in quotation marks.
* Can now also work with Identity Paths which include space in the path.

__IMPORTANT__ : the CLI package changed how it handles if a command has arguments with flags and specifies `SkipFlagParsing: true` - so we'll switch back to `SkipFlagParsing: false` - the command to run on the remote should be prefixed with `--` in case of `run`.

This means that `bitrise-machine --workdir=... run -timeout=0 ls -alh` won't work anymore,
you have to add the `--` separator: `bitrise-machine --workdir=... run -timeout=0 -- ls -alh`.

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.5/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.4 -> 0.9.5

* [5ecc6ea] Viktor Benei - added CLI package change note to 0.9.5 changelog (2015 Dec 31)
* [6b3917c] Viktor Benei - FIX : the CLI package changed how it handles if a command has arguments and specifies `SkipFlagParsing: true`- so we'll switch back to `SkipFlagParsing: false` and the command should be prefixed with `--` in case of run, if called with run flags (2015 Dec 31)
* [14e6ef4] Viktor Benei - final release notes for v0.9.5 (2015 Dec 31)
* [5e66ac0] Viktor Benei - Dockerfile : golang 1.5.2 (2015 Dec 31)
* [add86f9] Viktor Benei - Dockerfile: updated Bitrise CLI (1.2.4), new github path (bitrise-tools instead of bitrise-io) (2015 Dec 31)
* [6ebe72e] Viktor Benei - debug log the SSH commands (2015 Dec 31)
* [9de843f] Viktor Benei - v0.9.5 & changelog (2015 Dec 31)
* [b6de8a0] Viktor Benei - * Added support for `vagrant` v1.8.x 's new `vagrant ssh-config` output format, which now wraps the IdentityPath in quotation marks. * Can now also work with Identity Paths which include space in the path. (2015 Dec 31)
* [8059ecc] Viktor Benei - added more Debug info to setup:doSetupSSH (2015 Dec 31)
* [6bd987d] Viktor Benei - bitrise.yml : new workflow to create test binaries, for both OS X and Linux, with Go 1.5+ cross compile (2015 Dec 31)
* [db59c12] Viktor Benei - github URL change : moved from bitrise-io org to bitrise-tools (2015 Dec 31)
* [2b83a90] Viktor Benei - Godeps dependencies update (2015 Dec 31)


## 0.9.4 (2015 Oct 27)

### Release Notes

* __NEW__ flag : if you define `--abort-check-url` for `bitrise-machine run` it'll periodically check the given URL, and will abort the `run` if it receives a JSON response with `"status": "ok"` and `"is_aborted": true`

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.4/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.3 -> 0.9.4

* [1952623] Viktor Benei - changelog (2015 Oct 27)
* [e78fea4] Viktor Benei - v0.9.4 (2015 Oct 27)
* [651d4bb] Viktor Benei - abort-check-url handling : call the provided URL (if any) periodically, and abort the build if the URL returns the expected "is_aborted = true" response (2015 Oct 27)


## 0.9.3 (2015 Oct 19)

### Release Notes

* First version with official __LINUX__ binary release!
* New cleanup modes: `recreate` implemented, `custom` also available (you can specify a custom `vagrant` action for cleanup. For example for the DigitalOcean provider a built-in `recreate` action is available, which is more efficient than a full recreate).
* Environment Variables can also be specified in the `bitrise.machine.config.json` file.
* __NEW COMMAND__ : `destroy` - to help with completely destroying the "machine".

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.3/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.2 -> 0.9.3

* [a3bf9aa] Viktor Benei - v0.9.3 (2015 Oct 19)
* [d31e8af] Viktor Benei - godeps-update (2015 Oct 19)
* [eff00a6] Viktor Benei - new command: destroy - to easily destroy the host completely + files to build in Docker, for Linux (2015 Oct 19)
* [649005e] Viktor Benei - implemented: custom environments handling; cleanup mode 'recreate' and 'custom-action' ; vagrant '--machine-readable' processing (2015 Oct 01)
* [fe701a3] Viktor Benei - Godeps-update (2015 Oct 01)
* [8388a1e] Viktor Benei - goddess update (2015 Sep 30)
* [60ef9c1] Viktor Benei - godeps-update : added bitrise go-utils/testutil (2015 Sep 30)
* [0eab7ea] Viktor Benei - deps.go to force-include "`_test`" specific packages (2015 Sep 30)
* [67635a7] Viktor Benei - godeps-update (2015 Sep 30)
* [84e9afd] Viktor Benei - v0.9.3 (2015 Sep 30)


## 0.9.2 (2015 Sep 03)

### Release Notes

* log flush tuning
* Log Summary MetaInfo printing

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.2/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.1 -> 0.9.2

* [c35aa7b] Viktor Benei - 0.9.2 (2015 Sep 03)
* [cce440c] Viktor Benei - better `golint` CI (2015 Sep 02)
* [b58ece4] Viktor Benei - Control MetaInfo handling/printing implemented, the first (and so far only) one is `/logs/summary` (2015 Sep 02)
* [2f9f5cd] Viktor Benei - do log flush more frequently (3 sec instead of 5) (2015 Sep 02)


## 0.9.1 (2015 Sep 01)

### Release Notes

* __NEW__/__BREAKING__ action : `bitrise-machine setup` will skip the cleanup & SSH setup if it detects that the host is already prepared. You can force the full setup by adding the `--force` flag: `bitrise-machine setup --force`

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-io/bitrise-machine/releases/download/0.9.1/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.0 -> 0.9.1

* [dc80e95] Viktor Benei - v0.9.1 (2015 Sep 01)


-----------------

Generated at: 2017 Feb 20
