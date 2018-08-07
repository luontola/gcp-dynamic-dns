
# sync-k8s-nodes-to-gcp-dns

Syncs the external node IPs of a Kubernetes cluster to GCP Cloud DNS records. 
Designed to be used with a [Port Proxy](https://git.k8s.io/contrib/for-demos/proxy-to-service)
to get high enough availability without having to pay for a load balancer.  
