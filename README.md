# Simple XDP Firewall

This is a simple project that uses XDP and Dropbox's [gobpf](https://github.com/dropbox/goebpf) library for filtering received packets based on their IP address.
<br><br>


## How to run
----

The firewall contains 2 parts:

- The Kernel `eBPF` code written in `C` that needs to be compiled to `ELF` (Clang, llvm, etc)
- the Userspace Code written in `Go` that handles the logic of the firewall


## Prerequisites

We'll be using Clang to compile the `eBPF` code to `ELF`:

    $ apt-get install clang llvm make


You can either run `make` on the root directory to build the `ELF` file, or use `clang`:

    # Use tha make command 
    $ make 

    # Use Clang to compile to ELF
    $ cd bpf
    $ clang -I ../headers -O -target bpf -c xdp_drop.c -o xdp_drop.elf

The `ELF` file has to be named `xdp_drop.elf` because it's hardcoded into the source code.

<br>

## Run
Running an XDP/eBPF program, we're going to be needing root permission, because all processes that intend to load eBPF programs into the Linux kernel must be running in privileged mode. <br>
We also need to pass in the `-iface` flag with the value of the interface we want XDP to attach to. We'll attach `lo` in this example:

    $ go build           # or make build
    $ sudo ./xdp-blocker -iface lo


## Block IP Address

I've created a simple API that you can use to block or unblock an IP Address with a given subnet.

-   Block IP Address:  `POST /v1/add-ip`
-   Unblock IP Address: `POST /v1/remove-ip`

### Examples of using the API to block/unblock a given IP Address

    # Block IP Address(es)
    $ curl -XPOST -H "Content-Type: application/json" localhost:8080/v1/add-ip -d '
    {
        "ipAddress": "X.X.X.X", 
        "subnet": "32"
    }

    # Unblock IP Address(es)
    $ curl -XPOST -H "Content-Type: application/json" localhost:8080/v1/remove-ip -d '
    {
        "ipAddress": "X.X.X.X", 
        "subnet": "32"
    }


Hitting `CTRL-c` will detach the XDP Program and all blocked IP Addresses will be unblocked.

