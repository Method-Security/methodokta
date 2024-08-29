# methodokta Documentation

Hello and welcome to the methodokta documentation. While we always want to provide the most comprehensive documentation possible, we thought you may find the below sections a helpful place to get started.

- The [Getting Started](./getting-started/basic-usage.md) section provides onboarding material
- The [Development](./development/setup.md) header is the best place to get started on developing on top of and with methodokta
- See the [Docs](./index.md) section for a comprehensive rundown of methodokta capabilities

# About methodokta

methodokta provides security operators with powerful Okta enumeration capabilities, enabling them to gain comprehensive visibility into their Okta environments. Tailored for data-modeling and data-integration, methodokta can function independently as an interactive CLI, be integrated into broader security workflows, or be utilized within the Method Platform for enhanced identity and access management insights.

The number of security-relevant Okta resources that methodokta can enumerate are constantly growing. For the most up to date listing, please see the documentation [here](./index.md)

To learn more about methodokta, please see the [Documentation site](https://method-security.github.io/methodokta/) for the most detailed information.

## Quick Start

### Get methodokta

For the full list of available installation options, please see the [Installation](./getting-started/installation.md) page. For convenience, here are some of the most commonly used options:

- `docker run methodsecurity/methodokta`
- `docker run ghcr.io/method-security/methodokta`
- Download the latest binary from the [Github Releases](https://github.com/Method-Security/methodokta/releases/latest) page
- [Installation documentation](./getting-started/installation.md)

### Authentication
Authentication can be done in 2 ways:
1. By setting the `--api-token` and `--domain` flags
2. Setting the `$OKTA_API_TOKEN` and `$OKTA_API_DOMAIN` env variables

### General Usage

```bash
methodokta <resource> enumerate 
```

#### Examples

```bash
methodokta user enumerate
```

```bash
methodokta group enumerate
```

## Contributing

Interested in contributing to methodokta? Please see our organization wide [Contribution](https://method-security.github.io/community/contribute/discussions.html) page.

## Want More?

If you're looking for an easy way to tie methodokta into your broader cybersecurity workflows, or want to leverage some autonomy to improve your overall security posture, you'll love the broader Method Platform.

For more information, visit us [here](https://method.security)

## Community

methodoktais a Method Security open source project.

Learn more about Method's open source source work by checking out our other projects [here](https://github.com/Method-Security) or our organization wide documentation [here](https://method-security.github.io).

Have an idea for a Tool to contribute? Open a Discussion [here](https://github.com/Method-Security/Method-Security.github.io/discussions).
