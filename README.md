# Easycmd 
<a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/license-MIT-_red.svg"></a>
<a href="https://github.com/sincerefly/easycmd/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
<a href='https://coveralls.io/github/sincerefly/easycmd?branch=main'><img src='https://coveralls.io/repos/github/sincerefly/easycmd/badge.svg?branch=main' alt='Coverage Status' /></a>

A terminal tool sample, you can create your own tools base this app

## Usage

```bash
$ ./easycmd -h   
Long Terminal Usage desc

Usage:
  easycmd [flags]
  easycmd [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  ip          Query your local ip address
  version     Print the version number

Flags:
  -h, --help      help for easycmd
  -v, --version   output version

Use "easycmd [command] --help" for more information about a command.
```

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


