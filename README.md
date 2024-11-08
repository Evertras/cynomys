# Cynomys (cyn)

A simple network diagnostic testing tool to ensure communication behaves as
expected.

Available as a Docker image: `evertras/cynomys` ([Dockerhub link](https://hub.docker.com/r/evertras/cynomys))

Available as native binaries: [on releases page](https://github.com/Evertras/cynomys/releases).

## Why does this exist

I got tired of setting up `nc -stuff` on multiple machines and manually trying
to send data across to ensure a non-trivial network setup was working properly.
There was also a need to test broadcast/multicast, which got surprisingly weird
with multiple versions of `nc`.

Cynomys is intended to allow for simple network communication testing on a
variety of platforms with simple, consistent behavior.

## Features

### Current

Test connectivity between different machines on with UDP to make sure
the machines can talk to each other as expected.

Test broadcast/multicast (UDP).

Test connectivity and communication with TCP.

Customizable interval to send data.

Use in a Docker container for Docker-related networking, or just use the raw
binary for native level testing.

Optional HTTP server that provides a simple landing page to show current
status and configuration of the given `cyn` instance.

Measure RTT latency from senders to receivers.

### Future

These will be converted to issues in the future, but some other ideas...

Test that connectivity is NOT made between different machines that should not
talk to each other, for firewall/security reasons.

Allow metric collection (Prometheus, etc).

Realtime streaming updates for the HTML view.

## How to install it

Binaries are self-contained and available for most major platforms. Grab
[a native binary from the releases page](https://github.com/Evertras/cynomys/releases).

Or run with docker.

```bash
docker run --rm -it evertras/cynomys --help
```

## How to use it

For simple use cases, just use command line args. By convention, lowercase
means listen while uppercase means send.

```bash
# In this example:
# Machine A is listening on UDP on :1234 and TCP on :5555. It sends UDP to B and C.
# Machine B sends both UDP and TCP to Machine A, and UDP to C.
# Machine C only sends UDP to A and B.

# On Machine A - 192.168.58.2
cyn --listen-udp 192.168.58.2:1234 \
    --listen-tcp 192.168.58.2:5555 \
    --send-udp 192.168.58.3:3456 \
    --send-udp 192.168.58.4:3456

# On Machine B - 192.168.58.3 (shorthand flags)
cyn -u 192.168.58.3:3456 \
    -T 192.168.58.2:5555 \
    -U 192.168.58.2:1234 \
    -U 192.168.58.4:3456

# On Machine C - 192.168.58.4 (mixed, 1 minute send interval)
cyn -u 192.168.58.4:3456 \
    --send-udp 192.168.58.2:1234 \
    --send-interval 1m
```

```bash
# Listen for broadcast messages on Machine A and echo the data received
cyn -u :1234 -e

# Broadcast messages from Machine B using the regular UDP sender
cyn -U 192.168.58.255:1234
```

Instances that are listening will produce output when they receive messages.

```
2022/12/27 03:00:39 Read 2 bytes from 192.168.58.4:50372
2022/12/27 03:00:39 Received: hi
2022/12/27 03:00:39 Read 2 bytes from 192.168.58.3:54115
2022/12/27 03:00:39 Received: hi
```

### Configuration

Configuration is available through command line flags, a configuration file,
and/or environment variables.

#### Command line flags

All available flags can be seen by running `cyn --help`. This README snippet should
have the latest, but when in doubt, just run the command to check.

```text
$ cyn --help
Usage:
   [flags]
   [command]

Available Commands:
  help        Help about any command
  version     Displays the version

Flags:
  -c, --config string            A file path to load as additional configuration.
  -h, --help                     help for this command
      --http.address string      An address:port to host an HTTP server on for realtime data, such as '127.0.0.1:8080'
  -e, --listen.echo              If enabled, echo the data that's received. Otherwise just print the length received (default).
  -t, --listen.tcp strings       An IP:port address to listen on for TCP.  Can be specified multiple times.
  -u, --listen.udp strings       An IP:port address to listen on for UDP.  Can be specified multiple times.
  -d, --send.data string         The string data to send. (default "hi")
  -i, --send.interval duration   How long to wait between attempting to send data (default 1s)
  -T, --send.tcp strings         An IP:port address to send to (TCP).  Can be specified multiple times.
  -U, --send.udp strings         An IP:port address to send to (UDP).  Can be specified multiple times.
      --sinks.stdout.enabled     Whether to enable the stdout metrics sink

Use " [command] --help" for more information about a command.
```

#### Environment variables

Configuration can be supplied throuh environment variables that follow the
`CYNOMYS_VAR_NAME` pattern, replacing `.` with `_`. For example:

```bash
export CYNOMYS_SEND_INTERVAL=500ms
cyn -T 192.168.58.3:1234
```

#### Configuration File

A configuration file can be provided. This is useful when trying to template
configuration such as with Consul or similar tools.

A full configuration file with all options is given below.

```yaml
# my-cyn-config.yaml
listen:
  echo: true
  udp:
    - 192.168.58.4:2345
  tcp:
    - 192.168.58.4:2346
send:
  udp:
    - 192.168.58.3:1234
  tcp:
    - 192.168.58.3:1235
  interval: 30s
  data: "my custom data"
http:
  address: 127.0.0.1:8080
sinks:
  stdout:
    enabled: true
```

The configuration file must be supplied with the `--config/-c` command line
flag.

```bash
cyn --config ./my-cyn-config.yaml

# Can also use -c for shorthand
cyn -c ./my-cyn-config.yaml
```

## Why the name

[Prarie dogs talk to each other](https://en.wikipedia.org/wiki/Prairie_dog)

### How do I pronounce it?

idk
