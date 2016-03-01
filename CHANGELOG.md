## Changelog (Current version: 0.9.7)

-----------------

## 0.9.7 (2016 Mar 01)

### Release Notes

* __BREAKING__ : change 1
* change 2

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/bitrise-machine/releases/download/0.9.7/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.6 -> 0.9.7

* [ccdfbbd] Viktor Benei - Upgrades for Go 1.6 vendor folder - for testing (2016 Mar 01)
* [5bc5fca] Viktor Benei - Godeps.json (2016 Mar 01)
* [16c1e85] Viktor Benei - Docker image: use Go 1.6 (2016 Mar 01)
* [0a78829] Viktor Benei - deps update -> moved into vendor folder (2016 Mar 01)
* [7345dc3] Viktor Benei - v0.9.7 (2016 Mar 01)
* [d98dd36] Viktor Benei - Optimization: "tick" functions frequency reduced from 100 milliseconds to 500 milliseconds (2016 Mar 01)
* [3264036] Viktor Benei - optimization in ReadRunes function, as it was the main CPU bottleneck (2016 Mar 01)
* [dee129f] Viktor Benei - Go CPU profiling can now be turned-on with the `BITRISE_MACHINE_CPU_PROFILE_FILEPATH` Environment Variable (2016 Mar 01)
* [645472e] Viktor Benei - Release URL change (2016 Feb 29)
* [f50178c] Viktor Benei - changelog for 0.9.6 (2016 Jan 12)


## 0.9.6 (2016 Jan 12)

### Release Notes

* __BREAKING__ : change 1
* change 2

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/bitrise-machine/releases/download/0.9.6/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.5 -> 0.9.6

* [f4f2f9a] Viktor Benei - err check fix (2016 Jan 12)
* [47f2019] Viktor Benei - implemented CleanupMode destroy ; and global CLI params are now stored in Freezable objects, to prevent modification (2016 Jan 12)
* [440d333] Viktor Benei - godeps-update (2016 Jan 12)
* [f6eff19] Viktor Benei - finishing the command line Additional Env Vars handling (2016 Jan 11)
* [33a3602] Viktor Benei - removed the install_bitrise script (2016 Jan 11)
* [895c1ab] Viktor Benei - FIX : machine_config_test (2016 Jan 11)
* [6663d0a] Viktor Benei - godeps update (2016 Jan 11)
* [41ada6b] Viktor Benei - new flag: `--environment` / `-e` (2016 Jan 11)
* [ec25654] Viktor Benei - upgraded bitrise CLI version (2016 Jan 06)
* [5ecc6ea] Viktor Benei - added CLI package change note to 0.9.5 changelog (2015 Dec 31)


## 0.9.5 (2015 Dec 31)

### Release Notes

* __BREAKING__ : change 1
* change 2

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/bitrise-machine/releases/download/0.9.5/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.4 -> 0.9.5

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
* [1952623] Viktor Benei - changelog (2015 Oct 27)


## 0.9.4 (2015 Oct 27)

### Release Notes

* __BREAKING__ : change 1
* change 2

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/bitrise-machine/releases/download/0.9.4/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.3 -> 0.9.4

* [e78fea4] Viktor Benei - v0.9.4 (2015 Oct 27)
* [651d4bb] Viktor Benei - abort-check-url handling : call the provided URL (if any) periodically, and abort the build if the URL returns the expected "is_aborted = true" response (2015 Oct 27)
* [a3bf9aa] Viktor Benei - v0.9.3 (2015 Oct 19)


## 0.9.3 (2015 Oct 19)

### Release Notes

* __BREAKING__ : change 1
* change 2

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/bitrise-machine/releases/download/0.9.3/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.2 -> 0.9.3

* [d31e8af] Viktor Benei - godeps-update (2015 Oct 19)
* [eff00a6] Viktor Benei - new command: destroy - to easily destroy the host completely + files to build in Docker, for Linux (2015 Oct 19)
* [649005e] Viktor Benei - implemented: custom environments handling; cleanup mode 'recreate' and 'custom-action' ; vagrant '--machine-readable' processing (2015 Oct 01)
* [fe701a3] Viktor Benei - Godeps-update (2015 Oct 01)
* [8388a1e] Viktor Benei - goddess update (2015 Sep 30)
* [60ef9c1] Viktor Benei - godeps-update : added bitrise go-utils/testutil (2015 Sep 30)
* [0eab7ea] Viktor Benei - deps.go to force-include "`_test`" specific packages (2015 Sep 30)
* [67635a7] Viktor Benei - godeps-update (2015 Sep 30)
* [84e9afd] Viktor Benei - v0.9.3 (2015 Sep 30)
* [c35aa7b] Viktor Benei - 0.9.2 (2015 Sep 03)


## 0.9.2 (2015 Sep 03)

### Release Notes

* __BREAKING__ : change 1
* change 2

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/bitrise-machine/releases/download/0.9.2/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.1 -> 0.9.2

* [cce440c] Viktor Benei - better `golint` CI (2015 Sep 02)
* [b58ece4] Viktor Benei - Control MetaInfo handling/printing implemented, the first (and so far only) one is `/logs/summary` (2015 Sep 02)
* [2f9f5cd] Viktor Benei - do log flush more frequently (3 sec instead of 5) (2015 Sep 02)
* [dc80e95] Viktor Benei - v0.9.1 (2015 Sep 01)


## 0.9.1 (2015 Sep 01)

### Release Notes

* __BREAKING__ : change 1
* change 2

### Install or upgrade

To install this version, run the following commands (in a bash shell):

```
curl -fL https://github.com/bitrise-tools/bitrise-machine/releases/download/0.9.1/bitrise-machine-$(uname -s)-$(uname -m) > /usr/local/bin/bitrise-machine
```

Then:

```
chmod +x /usr/local/bin/bitrise-machine
```

That's all, you're ready to call `bitrise-machine`!

### Release Commits - 0.9.0 -> 0.9.1

* [6bc1f2f] Viktor Benei - readme : note about required `bitrise-bridge` (2015 Aug 31)


-----------------

Updated: 2016 Mar 01
