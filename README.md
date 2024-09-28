# ✨ gfping

## Install

```bash
go build
```

## 🚀 Usage

```bash
$ gping --help
Batch network probe tool.

Usage:
  gfping [command]

Available Commands:
  help        Help about any command
  icmp        icmp gfping network

Flags:
  -c, --concurrent int    number of goroutines to use (concurrent) (Default 300) (default 300)
  -f, --file string       read list of targets from a file
  -h, --help              help for gfping
  -i, --singleip string   single ip，EP: 192.168.1.1
  -g, --subnet string     generate target list (only if no -f -i specified), EP: 192.168.1.1/16
  -t, --timeout int       individual target initial timeout, unit ms (default 3000)

Use "gfping [command] --help" for more information about a command.
```

eg:

```bash
gfping icmp -i 192.168.0.100
gfping icmp -g 192.168.0.0/24
gfping icmp -f ./ipList.txt
```
