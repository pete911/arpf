# arpf

ARP tools

```
Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  scan        list MAC and IPs on the network

Flags:
  -h, --help      help for arpf
  -v, --verbose   print debug messages
```

## examples

Get IPs and MAC addresses on local network

```
en0: flags=up|broadcast|multicast mtu 1500
    ether 6c:66:60:64:67:6c
    inet 192.168.86.22 netmask 0xffffff00 broadcast 192.168.86.255
sending ARPs: src 192.168.86.22 dst 192.168.86.1 ... 192.168.86.254
sending packets for 10s in 3s intervals
dude.home. (192.168.86.112) at b6:26:e6:e6:e6:46
dude-two.home. (192.168.86.118) at b6:26:e6:96:46:96
```

