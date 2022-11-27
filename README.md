# Cynomys (cyn)

A network diagnostic testing tool to ensure communication behaves as expected.

## Features

### Current

### In progress

Test connectivity between different machines on with UDP/TCP to make sure
the machines can talk to each other as expected.

Test broadcast/multicast.

Use in a Docker container for Docker-related networking, or just use the raw
binary for native level testing.

### Future

Test that connectivity is NOT made between different machines that should not
talk to each other, for firewall/security reasons.

Config file for more advanced setups.

## How to use it

For simple use cases, just use command line args.

```bash
# On Machine A - 192.168.58.2
cyn --listen-udp 192.168.58.1:1234 --call-tcp 192.168.58.3:3456 --call-tcp 192.168.58.4:3456

# On Machine B - 192.168.58.3 (shorthand flags)
cyn -l 192.168.58.2:2345 -t 192.168.58.2:1234 -t 192.168.58.4:3456

# On Machine C - 192.168.58.4 (broadcasting)
cyn -l 192.168.58.3:3456 --broadcast 192.168.58.255
```

## Why the name

[Prarie dogs talk to each other](https://en.wikipedia.org/wiki/Prairie_dog)

### How do I pronounce it?

idk
