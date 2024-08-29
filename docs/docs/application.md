# Application

The `methodokta application` family of commands provide information about an Okta instance's applications

## Enumerate

The enumerate command will gather information about the applications in a given Organization.

### Usage

```bash
methodokta application enumerate
```

### Help Text

```bash
% methodokta application -h
Audit and command Users

Usage:
  methodokta user [command]

Available Commands:
  enumerate   Enumerate Applications

Flags:
  -h, --help   help for application

Global Flags:
  -t, --apitoken string     Okta API Token
  -d, --domain string        Okta Domain
  -o, --output string        Output format (signal, json, yaml). Default value is signal (default "signal")
  -f, --output-file string   Path to output file. If blank, will output to STDOUT
  -q, --quiet                Suppress output
  -v, --verbose              Verbose output
  ```