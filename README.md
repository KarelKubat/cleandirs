# cleandirs

<!-- toc -->
- [Installation](#installation)
- [Trying it out](#trying-it-out)
- [Typical usage](#typical-usage)
<!-- /toc -->

`cleandirs` is a small utility, written in Go, which cleans out files beyond a certain age in your temporary directories, typically `/tmp`. Empty subdirectories are removed too.

- Prerequisites: Go V1.17 (or better)
- Dependencies: none, self-contained

## Installation

1. Run `go get github.com/KarelKubat/cleandirs`

1. Change-dir to where the files are, typically `~/go/src/github.com/KarelKubat/cleandirs/`

1. Install using `go install cleandirs.go`. This will build and deploy a binary to wherever your Go programs go, maybe `~/go/bin/`.

## Trying it out

Just run `cleandirs` to see it in action in dry-run mode. The defaults are:

- Files and empty subdirectories would be removed under `/tmp`. Use the flag `-dirs-to-clean=$DIR1,$DIR2,$DIR3` to extend the list or to specify a different temporary directory. The argument to `-dirs-to-clean` is a comma-separated list.

- Without `-dry-run=false` the utility will only spit out the files that are candidate for removal, that's why when you leave out this flag, nothing happens.

- If you don't want `cleandirs` to remove empty subdirs under your temporary directories, add `-prune-dirs=false`.

- If you want `cleandirs` to also remove special files (pipes, FIFO's), add `--all-files`. Default is to consider only regular files.

- Stale files are considered those entries that are more than a day old. Use `-ttl=...` to specify a different time. E.g., if you want to keep up to 1week old files, use `-ttl=168h`.

- When `-version` is given, then `cleandirs` only reports its version and then stops. When `-help` is given, then an overview of all flags is shown and nothing else happens.

## Typical usage

To clean files out of `/tmp` you will probably want user `root` to run this, e.g. via `root`'s crontab. (Have a good look at the sources if you want to check whether you can trust this code. It's not long and easy to read, pretty straight-forward.)

Here is an example:

```shell
PATH=/bin:/usr/bin:/sbin:/usr/sbin:/WHEREVER/YOU/HAVE/THE/CLEANDIRS/BINARY

# Remove cleandir logging at reboot
@reboot rm -f /tmp/cleandirs.log

# Clean under /tmp/ what's older than a day
5 * * * * cleandirs -dry-run=false >>/tmp/cleandirs.log 2>&1

# Clean under /var/log what's older than a week
10 * * * * cleandirs -dry-run=false -dirs-to-clean=/var/log -ttl=168h >>/tmp/cleandirs.log 2>&1
```