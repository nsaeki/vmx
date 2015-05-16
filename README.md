# vmx

Thin CLI Wrapper for VMware `vmrun`.

Some VMware products have a command line application [`vmrun`](https://www.vmware.com/support/ws55/doc/ws_learning_cli_vmrun.html), it is useful for command line operation or launching VMs in no GUI mode.
`vmrun` requires vmx file path for its argument, and typing vmx path is annoying because default VMs are created in deep in home directory and those names are included spaces.

This utility runs `vmrun` with a simple VM name like:

```bash
$ vmx start Debian
```

Your virtual machines are found by:

```bash
$ vmx list
```

This tools appends `nogui` option for `start` command.
If you want to launch a VM in GUI mode, specify `gui` option explicitly:

```bash
$ vmx start Debian gui
```

## Install

Just go get this repository

```bash
go get github.com/nsaeki/vmx
```

## Caveats

Supports only VMware Fusion for Mac.
