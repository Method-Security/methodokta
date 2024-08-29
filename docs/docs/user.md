# User

The `methodokta user` family of commands provide information about an Okta instance's users

## Enumerate

The enumerate command will gather information about the users in a given organization.

### Usage

```bash
methodokta user enumerate
```

### Help Text

```bash
% methodokta user -h
Audit and command Users

Usage:
  methodokta user [command]

Available Commands:
  enumerate   Enumerate Users

Flags:
  -h, --help   help for user

Global Flags:
  -t, --apitoken string     Okta API Token
  -d, --domain string        Okta Domain
  -o, --output string        Output format (signal, json, yaml). Default value is signal (default "signal")
  -f, --output-file string   Path to output file. If blank, will output to STDOUT
  -q, --quiet                Suppress output
  -v, --verbose              Verbose output
  ```