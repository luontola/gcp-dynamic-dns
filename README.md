# Google Cloud Dynamic DNS Client

Syncs the current IP address, based on the local network interface, to Google Cloud DNS records.

[![Docker Build Status](https://img.shields.io/docker/cloud/build/luontola/gcp-dynamic-dns.svg)](https://hub.docker.com/r/luontola/gcp-dynamic-dns)
[![Docker Image Size](https://images.microbadger.com/badges/image/luontola/gcp-dynamic-dns.svg)](https://microbadger.com/images/luontola/gcp-dynamic-dns)

## Using

Run this application's [container](https://hub.docker.com/r/luontola/gcp-dynamic-dns) using the command `sync` (
entrypoint `/app`), `host` network and restart policy `always`. It will then sync the IP address automatically whenever
it changes.

For a list of other commands, run the `help` command.

### Environment variables

#### `INTERFACE_NAME` (optional)

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
