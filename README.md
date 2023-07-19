# Google Cloud Dynamic DNS Client

Syncs the current IP address to Google Cloud DNS records. Can discover your public IP based on
(1) a 3rd party web service, or (2) directly from the local network interface.

[![Docker Build Status](https://img.shields.io/docker/cloud/build/luontola/gcp-dynamic-dns.svg)](https://hub.docker.com/r/luontola/gcp-dynamic-dns)

## Using

Run this application's [container](https://hub.docker.com/r/luontola/gcp-dynamic-dns) using the command `sync` (
entrypoint `/app`), `host` network and restart policy `always`. It will then sync the IP address automatically whenever
it changes.

For a list of other commands, run the `help` command.

### Environment variables

#### `MODE` (optional)

The method for determining your public IP address. Possible values:

- `service` (default) - Asks a 3rd party web service for your external IP address.
- `interface` - Asks your operating system for the IP address assigned to a network interface.
- TODO: `upnp` - Asks your network router for its external IP address using Universal Plug and Play.

Default: `service`

#### `SERVICE_URLS` (optional, MODE=service)

Web addresses of services which report your public IP address. Multiple services may be separated by space, in which
case they will be used in a round-robin fashion. The continuous check interval is 5 minutes, so use more than once
service to call each individual service less often.

Default: `https://ipv4.icanhazip.com/ https://checkip.amazonaws.com/ https://ifconfig.me/ip https://ipinfo.io/ip`

#### `INTERFACE_NAME` (optional, MODE=interface)

Name of the network interface whose IP to use. If not defined, the program will detect the primary network interface
automatically.

Example: `eth0`

#### `DNS_NAMES`

List of domain names to update. Separate the domain names with one space. Each name must end with a period. The DNS
records must already exist on Cloud DNS and they must be type `A` records.

Example: `example.com. subdomain.example.com. example.org.`

#### `GOOGLE_PROJECT`

The name of your Google Cloud project. The above mentioned DNS names must be hosted under this project's Cloud DNS.

Example: `your-project-123456`

#### `GOOGLE_APPLICATION_CREDENTIALS`

Path to service account credentials with permissions to update your Cloud DNS records.

Example: `/path/to/dns-updater-gcp-keys.json`

> In the GCP console, under **IAM & admin > Service accounts**, create a service account for the application. Name
> it `dns-updater`, grant it the **DNS > DNS Administrator** role, create a key for it in JSON format and save it
> as `dns-updater-gcp-keys.json`.

## Developing

Run tests and build the project

    docker-compose build --force-rm

Run the application

    docker-compose run --rm app help

The application container doesn't have `sh` or other fancy stuff,
so to inspect its contents use the `docker export` command:

    docker export <container> | tar tv
