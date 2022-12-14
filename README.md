# Cynomys (cyn)

A network diagnostic testing tool to ensure communication behaves as expected.

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

## How to use it

For simple use cases, just use command line args.  By convention, lowercase
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

A configuration file can be provided.  This is useful when trying to template
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
