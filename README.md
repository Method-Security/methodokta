<div align="center">
<h1>methodokta</h1>

[![GitHub Release][release-img]][release]
[![Verify][verify-img]][verify]
[![Go Report Card][go-report-img]][go-report]
[![License: Apache-2.0][license-img]][license]

[![GitHub Downloads][github-downloads-img]][release]
[![Docker Pulls][docker-pulls-img]][docker-pull]

</div>

methodokta provides security operators with a number of data-rich Okta enumeration capabilities to help them gain visibility into their Okta Instance. Designed with data-modeling and data-integration needs in mind, methodokta can be used on its own as an interactive CLI, orchestrated as part of a broader data pipeline, or leveraged from within the Method Platform.

The number of security-relevant Okta resources that methodokta can enumerate are constantly growing. For the most up to date listing, please see the documentation [here](docs-capabilities)

To learn more about methodokta, please see the [Documentation site](https://method-security.github.io/methodokta/) for the most detailed information.

## Quick Start

### Get methodokta

For the full list of available installation options, please see the [Installation](./docs/getting-started/index.md) page. For convenience, here are some of the most commonly used options:

- `docker run methodsecurity/methodokta`
- `docker run ghcr.io/method-security/methodokta:0.0.1`
- Download the latest binary from the [Github Releases](releases) page
- [Installation documentation](./docs/getting-started/index.md)

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

Interested in contributing to methodokta? Please see our [Contribution](#) page.

## Want More?

If you're looking for an easy way to tie methodokta into your broader cybersecurity workflows, or want to leverage some autonomy to improve your overall security posture, you'll love the broader Method Platform.

For more information, see [https://method.security]

## Community

methodokta is a Method Security open source project.

Learn more about Method's open source source work by checking out our other projects [here](github-org).

Have an idea for a Tool to contribute? Open a Discussion [here](discussion).

[verify]: https://github.com/Method-Security/methodokta/actions/workflows/verify.yml
[verify-img]: https://github.com/Method-Security/methodokta/actions/workflows/verify.yml/badge.svg
[go-report]: https://goreportcard.com/report/github.com/Method-Security/methodokta
[go-report-img]: https://goreportcard.com/badge/github.com/Method-Security/methodokta
[release]: https://github.com/Method-Security/methodokta/releases
[releases]: https://github.com/Method-Security/methodokta/releases/latest
[release-img]: https://img.shields.io/github/release/Method-Security/methodokta.svg?logo=github
[github-downloads-img]: https://img.shields.io/github/downloads/Method-Security/methodokta/total?logo=github
[docker-pulls-img]: https://img.shields.io/docker/pulls/methodsecurity/methodokta?logo=docker&label=docker%20pulls%20%2F%20methodokta
[docker-pull]: https://hub.docker.com/r/methodsecurity/methodokta
[license]: https://github.com/Method-Security/methodokta/blob/main/LICENSE
[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[homepage]: https://method.security
[docs-home]: https://method-security.github.io/methodokta
[docs-capabilities]: https://method-security.github.io/methodokta/docs/index.html
[discussion]: https://github.com/Method-Security/methodokta/discussions
[github-org]: https://github.com/Method-Security