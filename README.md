# Cynomys

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
# On Machine A - 192.168.0.1
cyn --listen 192.168.0.1 --call-tcp 192.168.0.2 --call-tcp 192.168.0.3

# On Machine B - 192.168.0.2 (shorthand flags)
cyn -l 192.168.0.2 -t 192.168.0.1 -t 192.168.0.3

# On Machine C - 192.168.0.3 (broadcasting)
cyn -l 192.168.0.3 --broadcast 192.168.0.255
```

## Why the name

[Prarie dogs talk to each other](https://en.wikipedia.org/wiki/Prairie_dog)
