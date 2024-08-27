# Device

The `methodokta device` family of commands provide information about an Okta instance's devices

## Enumerate

The enumerate command will gather information about the devices in a given organization.

### Usage

```bash
methodokta device enumerate
```

### Help Text

```bash
% methodokta device -h
Audit and command Users

Usage:
  methodokta device [command]

Available Commands:
  enumerate   Enumerate Devices

Flags:
  -h, --help   help for device

Global Flags:
  -t, --apitoken string     Okta API Token
  -d, --domain string        Okta Domain
  -o, --output string        Output format (signal, json, yaml). Default value is signal (default "signal")
  -f, --output-file string   Path to output file. If blank, will output to STDOUT
  -q, --quiet                Suppress output
  -v, --verbose              Verbose output
  ```