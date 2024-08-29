# Basic Usage

Before you get started, you will need to export Okta credentials that you want methodOkta to utilize as environment variables. 

## Binaries

Running as a binary means you don't need to do anything additional for methodokta to leverage the environment variables you have already exported. You can test that things are working properly by running:

```bash
methodokta user enumerate --api-token xxxx --domain myoktadomain
```

## Docker

Running an authenticated workflow with methodokta as a Docker container requires that you pass okta credentials to the container:

```bash
docker run -e OKTA_API_TOKEN="XXXX" -e OKTA_DOMAIN="myoktadomain" methodsecurity/methodokta
```
