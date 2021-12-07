# Cloudflare Spectrum DDNS

Cloudflare Spectrum's functionality is limited to specifying IP address for SSH and Minecraft applications (at least for the Pro subscription). This is a problem for homelabbers, self-hosters, or anyone else that has a dynamic IP address.

This project serves to fill a need left by this limited functionality by providing tooling to:
1. Determine your current public IP address
2. Update your Cloudflare Spectrum application's IP address upon change.

## Deployment
This project has been designed to allow the aforementioned tooling to be consumed in a variety of ways.
- Run manually as a one-off command
- Run as a native or Kubernetes cronjob
- Dockerized in vanilla Docker, Docker Compose, or Kubernetes (see docker installation below)
- As a library see the `pkg` directory for the public `IPChecker`, `IPPoller` and `SpectrumClient` utilities.

## Requirements

- Existing DDNS solution (Dyn.com, NoIP.com, etc)
- Docker or Linux (should work on macOS and Windows but has not been tested)

## Installation

### Native (golang binary)

```
go get github.com/connormckelvey/cloudflare-spectrum-ddns/cmd/update-spectrum-ip
```

### Docker

```
docker pull connormckelvey/cloudflare-spectrum-ddns
```

## Configuration

Configuration can be done through environment variables or command line flags. The only exceptions are:
- `CLOUDFLARE_API_KEY` and `CLOUDFLARE_API_EMAIL` which can only be set via environment variables
- `-poll` and `-debug` which can only be set with command line flags.

**NOTE** Flags have a higher priority than environment variables and take precedence.

### Dotenv (.env) Files
The update-spectrum-dns binary will automatically load a `.env` file from the current working directory.


### Definitions

|Environment Variable|Flag|Description|Example
|---|---|---|---|
|`CLOUDFLARE_API_KEY`|x|User API Key for Cloudflare Account|`98eb470b2b60482e259d28648895d9e1`|
|`CLOUDFLARE_API_EMAIL`|x|Email address associated with API Key|`user@example.com`|
|`CLOUDFLARE_ZONE_NAME`|`-zone`|Cloudflare Zone or "Site" name hosting the Spectrum application|`example.com`|
|`SPECTRUM_APP_DOMAIN`|`-app-domain`|Domain name to be used for Spectrum App|`minecraft.example.com`|
|`SPECTRUM_APP_PROTOCOL`|`-app-protocol`|Protocol for the Spectrum application (`ssh` or `minecraft`)| `minecraft`
|`DDNS_HOSTNAME`|`-ddns-hostname`|Hostname (from DDNS service) used to get latest IP address|`example.noip.com`|
|`DNS_SERVER`|`-dns-server`|DNS Server used to query DDNSHostname|8.8.8.8|
|x|`-poll`|Check and apply updates continuously, omit this to only run once (cron)|`-poll`|
|x|`-poll`|Enable verbose logging|`-debug`|


## Usage

### Native (golang binary and .env)

```
echo CLOUDFLARE_API_KEY=xxx >> .env
echo CLOUDFLARE_API_EMAIL=xxx >> .env
```
(optionally define more environment variables in the .env file instead of using the command line flags below.)


```
update-spectrum-ip \
    -zone example.com \
    -app-domain minecraft.example.com \
    -ddns-hostname example.noip.com \
    -dns-server 1.1.1.1 \
    -poll \
    -debug
```

### Docker

```
echo CLOUDFLARE_API_KEY=xxx >> .env
echo CLOUDFLARE_API_EMAIL=xxx >> .env
```
(optionally define more environment variables in the .env file instead of using the command line flags below.)

```
docker run --env-file .env connormckelvey/cloudflare-spectrum-ddns
    -zone example.com \
    -app-domain minecraft.example.com \
    -ddns-hostname example.noip.com \
    -dns-server 1.1.1.1 \
    -poll \
    -debug
```

*OR* 

```
 docker run -v "/path/to/.env:/app/.env" cloudflare-spectrum-ddns:dev \
    -zone example.com \
    -app-domain minecraft.example.com \
    -ddns-hostname example.noip.com \
    -dns-server 1.1.1.1 \
    -poll \
    -debug
 ```


### Docker Compose (.env)

