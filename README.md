
# Sync Kubernetes nodes to Google Cloud DNS

Syncs the external node IPs of a Kubernetes cluster to Google Cloud DNS records. 
Designed to be used with a [Port Proxy](https://git.k8s.io/contrib/for-demos/proxy-to-service)
to get high enough availability (using client-side DNS failover) without having to pay for a load balancer.

[![Docker Build Status](https://img.shields.io/docker/build/luontola/sync-k8s-nodes-to-gcp-dns.svg)](https://hub.docker.com/r/luontola/sync-k8s-nodes-to-gcp-dns/)


## Using

### Kubernetes permissions

[As a prerequisite](https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control#prerequisites_for_using_role-based_access_control) you will need to grant yourself some permissions. Because of the way Kubernetes Engine checks permissions when you create a Role or ClusterRole, you must first create a RoleBinding that grants you all of the permissions included in the role you want to create.

Here is an example role binding which gives your Google identity the `cluster-admin` role, after which you can freely create additional Role and ClusterRole permissions.

```yaml
# cluster-admins.yaml
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cluster-admins
subjects:
- kind: User
  # update this
  name: your.name@gmail.com
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: ""
```

Apply this with:

    kubectl apply -f cluster-admins.yaml

Next create a service account which permits the application to find out the node IPs of the Kubernetes cluster.

```yaml
# dns-updater-account.yaml
kind: ServiceAccount
apiVersion: v1
metadata:
  name: dns-updater-account
  namespace: default
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: dns-updater-role
rules:
- apiGroups: [""]
  resources: ["nodes"]
  verbs: ["list"]
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: dns-updater-role-binding
subjects:
- kind: ServiceAccount
  name: dns-updater-account
  namespace: default
roleRef:
  kind: ClusterRole
  name: dns-updater-role
  apiGroup: rbac.authorization.k8s.io
```

Apply this with:

    kubectl apply -f dns-updater-account.yaml

### Cloud DNS permissions

You will also need to grant the application permissions to update your Cloud DNS records.

In the GCP console, under **IAM & admin > Service accounts**, create a service account for the application. Name it `dns-updater`, grant it the **DNS > DNS Administrator** role, create a key for it in JSON format and save it as `dns-updater-gcp-keys.json`.

Store the key in Kubernetes as a secret:

    kubectl create secret generic dns-updater-gcp-keys --from-file=gcp-keys.json=dns-updater-gcp-keys.json

### Application deployment

Finally you can deploy the application to Kubernetes. In `DNS_NAMES` list all the domain names you wish to update. They all must be type `A` DNS records. Separate the domain names with one space. Each name must end with a period. The domains must be under the Cloud DNS of the project specified in `GOOGLE_PROJECT`.

```yaml
# dns-updater-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dns-updater
spec:
  selector:
    matchLabels:
      app: dns-updater
  replicas: 1
  template:
    metadata:
      labels:
        app: dns-updater
    spec:
      containers:
      - name: dns-updater
        image: luontola/sync-k8s-nodes-to-gcp-dns
        imagePullPolicy: Always
        command: ["/app", "sync"]
        env:
        - name: DNS_NAMES
          # update this
          value: example.com. subdomain.example.com. example.org.
        - name: GOOGLE_PROJECT
          # update this
          value: your-project-123456
        - name: GOOGLE_APPLICATION_CREDENTIALS
          value: /var/secrets/gcp-keys/gcp-keys.json
        volumeMounts:
        - name: gcp-keys-vol
          mountPath: /var/secrets/gcp-keys
          readOnly: true
        resources:
          requests:
            memory: 10Mi
            cpu: 0m
          limits:
            memory: 20Mi
      volumes:
      - name: gcp-keys-vol
        secret:
          secretName: dns-updater-gcp-keys
      serviceAccountName: dns-updater-account
```

Apply this with:

    kubectl apply -f dns-updater-deployment.yaml

## Developing

Run tests and build the project

    docker-compose build

Run the application

    docker-compose run --rm app help

The application container doesn't have `sh` or other fancy stuff,
so to inspect its contents use the `docker export` command:

    docker export <container> | tar tv
