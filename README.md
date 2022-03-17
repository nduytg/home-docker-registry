# home-docker-registry
Build your own Docker Registry on K8s from scratch

Author: nduytg@gmail.com
Github: github.com/nduytg

## Preparation - Setup local K8s env by Minikube

To simplify the installation steps, we will use mini

```bash
# Download minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
sudo install minikube-darwin-amd64 /usr/local/bin/minikube

# Start minikube, wait for 5mins
minikube start

# Check if the cluster is up and running
❯ minikube kubectl -- get pods -A
NAMESPACE     NAME                               READY   STATUS    RESTARTS      AGE
kube-system   coredns-64897985d-jcw87            1/1     Running   0             90s
kube-system   etcd-minikube                      1/1     Running   0             100s
kube-system   kube-apiserver-minikube            1/1     Running   0             100s
kube-system   kube-controller-manager-minikube   1/1     Running   0             100s
kube-system   kube-proxy-h9z2r                   1/1     Running   0             90s
kube-system   kube-scheduler-minikube            1/1     Running   0             100s
kube-system   storage-provisioner                1/1     Running   1 (58s ago)   98s

# Set alias for minikube
# Add this in your bash profile
vim ~/.bash_profile
...
alias kubectl="minikube kubectl --"
alias k="kubectl"

source ~/.bash_profile
```

## Start docker registry deployment

```bash
k apply -f docker-registry/deployment.yml
```

Deployment file
```yaml
TODO: Put contents here after we finish
```

## Create service to access the app from localhost

```bash
❯ minikube tunnel
✅  Tunnel successfully started

📌  NOTE: Please do not close this terminal as this process must stay alive for the tunnel to be accessible ...

🏃  Starting tunnel for service docker-registry.


❯ kubectl get svc
NAME              TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
docker-registry   LoadBalancer   10.101.208.41   127.0.0.1     8080:30260/TCP   13m
kubernetes        ClusterIP      10.96.0.1       <none>        443/TCP          45m

```

```
/ # cat /etc/docker/registry/config.yml
version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
```


## Redis configuration + Replica set



# Further improvements

Use [S3](https://docs.docker.com/registry/storage-drivers/) as backend driver for Docker registry


# Reference

Reference documents for the assignment

https://docs.docker.com/registry/

https://docs.docker.com/registry/storage-drivers/

https://hub.docker.com/_/registry


