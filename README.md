
# Google Cloud Dynamic DNS Client

Syncs the current IP address, based on the local network interface, to Google Cloud DNS records. 

[![Docker Build Status](https://img.shields.io/docker/cloud/build/luontola/gcp-dynamic-dns.svg)](https://hub.docker.com/r/luontola/gcp-dynamic-dns/)
[![Docker Image Size](https://images.microbadger.com/badges/image/luontola/gcp-dynamic-dns.svg)](https://microbadger.com/images/luontola/gcp-dynamic-dns)


## Using

### Cloud DNS permissions

You will need to grant the application permissions to update your Cloud DNS records.

In the GCP console, under **IAM & admin > Service accounts**, create a service account for the application. Name it `dns-updater`, grant it the **DNS > DNS Administrator** role, create a key for it in JSON format and save it as `dns-updater-gcp-keys.json`.


### Application deployment

In `DNS_NAMES` list all the domain names you wish to update. They all must be type `A` DNS records. Separate the domain names with one space. Each name must end with a period. The domains must be under the Cloud DNS of the project specified in `GOOGLE_PROJECT`.

TODO


## Developing

Run tests and build the project

    docker-compose build

Run the application

    docker-compose run --rm app help

The application container doesn't have `sh` or other fancy stuff,
so to inspect its contents use the `docker export` command:

    docker export <container> | tar tv
