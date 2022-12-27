# Cynomys (cyn)

A network diagnostic testing tool to ensure communication behaves as expected.

Available as a Docker image: `evertras/cynomys` ([Dockerhub link](https://hub.docker.com/r/evertras/cynomys))

## Why does this exist

I got tired of setting up `nc -stuff` on multiple machines and manually trying
to send data across to ensure a non-trivial network setup was working properly.
With `cyn`, it can be run quickly and allow me to poke at configurations while
making sure that connectivity remains (or fails). There was also a need to test
broadcast/multicast, which got surprisingly weird with different versions of
`nc`.

## Features

### Current

Test connectivity between different machines on with UDP to make sure
the machines can talk to each other as expected.

Test broadcast/multicast (UDP).

Use in a Docker container for Docker-related networking, or just use the raw
binary for native level testing.

### In progress

Test connectivity with TCP.

### Future

Test that connectivity is NOT made between different machines that should not
talk to each other, for firewall/security reasons.

Customizable data to send.

Customizable intervals.

Allow metric collection (Prometheus, etc).

## How to install it

Binaries are self-contained and available for most major platforms. Grab
[a native binary from the releases page](https://github.com/Evertras/cynomys/releases).

Run with docker to see available flags.

```bash
docker run --rm -it evertras/cynomys
```

## How to use it

For simple use cases, just use command line args. By convention, lowercase
means listen while uppercase means send.

```bash
# On Machine A - 192.168.58.2
cyn --listen-udp 192.168.58.2:1234 --send-udp 192.168.58.3:3456 --send-udp 192.168.58.4:3456

# On Machine B - 192.168.58.3 (shorthand flags)
cyn -u 192.168.58.3:2345 -U 192.168.58.2:1234 -U 192.168.58.4:3456

# On Machine C - 192.168.58.4 (mixed)
cyn -u 192.168.58.4:3456 --send-udp 192.168.58.2:2345
```

```bash
# Listen for broadcast messages on Machine A
cyn -u :1234

# Broadcast messages from Machine B using the regular UDP sender
cyn -U 192.168.58.255:1234
```

### Configuration file

A configuration file can be provided. This is useful when trying to template
configuration such as with Consul or similar tools.

A full configuration file with all options is given below.

```yaml
# my-cyn-config.yaml
listen-udp:
  - 192.168.58.4:2345
send-udp:
  - 192.168.58.3:1234
```

The configuration is loaded via file.

```bash
cyn --config ./my-cyn-config.yaml

# Can also use -c for shorthand
cyn -c ./my-cyn-config.yaml
```

## Why the name

[Prarie dogs talk to each other](https://en.wikipedia.org/wiki/Prairie_dog)

### How do I pronounce it?

idk
