# Group

The `methodokta group` family of commands provide information about an Okta instance's groups

## Enumerate

The enumerate command will gather information about the groups in a given organization.

### Usage

```bash
methodokta group enumerate
```

### Help Text

```bash
% methodokta group -h
Audit and command Users

Usage:
  methodokta group [command]

Available Commands:
  enumerate   Enumerate Groups

Flags:
  -h, --help   help for group

Global Flags:
  -t, --apitoken string     Okta API Token
  -d, --domain string        Okta Domain
  -o, --output string        Output format (signal, json, yaml). Default value is signal (default "signal")
  -f, --output-file string   Path to output file. If blank, will output to STDOUT
  -q, --quiet                Suppress output
  -v, --verbose              Verbose output
  ```