# easycmd

A terminal tool template, you can create your own tools base this app

## Usage

### Version

```bash
# equal --version or version
$ easycmd -v 
0.0.3 (9a23740 2023-01-11)
```

### Ip

Query your host ip address

```bash
# random service
$ easy ip
https://icanhazip.com                     43.77.88.118
```

```bash
# all services
$ easy ip -a
https://ifconfig.me                       43.77.88.118
http://ip.sb                              43.77.88.118
http://ip.gs                              43.77.88.118
https://icanhazip.com                     43.77.88.118
https://ipecho.net/plain                  43.77.88.118
https://ifconfig.minidump.info/ip         43.77.88.118
https://ip.3322.org                       43.77.88.118
```

```bash
# query with given service
$ ./easycmd ip -s http://ip.sb
http://ip.sb                              43.77.88.118
```


