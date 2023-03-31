# shieldoo cli

Shieldoo cli can be used for Infrastructure as a Code automation.

# build

```bash
go build -o out/shieldoo
```

# man

Befor you will use cli tool `shieldoo` you must set environment variables:

- `SHIELDOO_URI` - use shieldoo Uri which you can find in shieldoo admin portal
- `SHIELDOO_APIKEY` - shieldoo ApiKey which you can find in shieldoo adimn portal

## shieldoo

```
A simple CLI tool to manage shieldoo servers and firewalls.

Usage:
  shieldoo [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  firewall    Manage firewall settings
  group       Manage groups
  help        Help about any command
  server      Manage servers

Flags:
  -h, --help   help for shieldoo

Use "shieldoo [command] --help" for more information about a command.
```

### shieldoo group

```
Manage groups

Usage:
  shieldoo group [command]

Available Commands:
  list        List all groups
  show        Show a group

Flags:
  -h, --help   help for group

Use "shieldoo group [command] --help" for more information about a command.
```

### shieldoo firewall

```
Manage firewall settings

Usage:
  shieldoo firewall [command]

Available Commands:
  delete      Delete a firewall
  ensure      Ensure a firewall (create or update)
  list        List firewalls
  show        Show a firewall

Flags:
  -h, --help   help for firewall

Use "shieldoo firewall [command] --help" for more information about a command.
```

### shieldoo server

```
Manage servers

Usage:
  shieldoo server [command]

Available Commands:
  delete      Delete a server
  ensure      Ensure a server (create or update)
  list        List all servers
  show        Show a server

Flags:
  -h, --help   help for server

Use "shieldoo server [command] --help" for more information about a command.
```
