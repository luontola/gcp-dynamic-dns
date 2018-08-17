
# Sync Kubernetes nodes to Google Cloud DNS

Syncs the external node IPs of a Kubernetes cluster to Google Cloud DNS records. 
Designed to be used with a [Port Proxy](https://git.k8s.io/contrib/for-demos/proxy-to-service)
to get high enough availability without having to pay for a load balancer.  

[![Docker Build Status](https://img.shields.io/docker/build/luontola/sync-k8s-nodes-to-gcp-dns.svg)](https://hub.docker.com/r/luontola/sync-k8s-nodes-to-gcp-dns/)


## Developing

Run tests and build the project

    docker-compose build

Run the application

    docker-compose run --rm app help

The application container doesn't have `sh` or other fancy stuff,
so to inspect its contents use the `docker export` command:

    docker export <container> | tar tv
