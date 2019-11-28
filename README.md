# Telegram client as middleman

## Usage
```text
$ tgmid -h
Telegram client as middleman

Usage:
  tgmid [flags]
  tgmid [command]

Available Commands:
  gen-completion Generate shell completion file
  help           Help about any command
  request        Send a request and wait for a response
  start          Start the telegram client
  version        Print the version number

Flags:
  -c, --config-dir string   Set config dir (default "./configs/")
  -h, --help                help for tgmid
  -v, --verbose count       Verbose output

Use "tgmid [command] --help" for more information about a command.
```

## Build
```bash
git clone https://github.com/immid/tgmid.git
cd tgmid
make build
```
