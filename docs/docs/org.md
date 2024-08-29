# Org

The `methodokta org` family of commands provide information about an Okta organization

## Enumerate

The enumerate command will gather information about the organization

### Usage

```bash
methodokta org enumerate
```

### Help Text

```bash
% methodokta org -h
Audit and command Users

Usage:
  methodokta org [command]

Available Commands:
  enumerate   Enumerate Org

Flags:
  -h, --help   help for org

Global Flags:
  -t, --apitoken string     Okta API Token
  -d, --domain string        Okta Domain
  -o, --output string        Output format (signal, json, yaml). Default value is signal (default "signal")
  -f, --output-file string   Path to output file. If blank, will output to STDOUT
  -q, --quiet                Suppress output
  -v, --verbose              Verbose output
  ```