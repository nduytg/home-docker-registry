# home-docker-registry
Build your own Docker Registry on K8s from scratch

Author: nduytg@gmail.com

Github: github.com/nduytg

## Requirements

Write the Kubernetes deployment manifest to run Docker Registry in Kubernetes with at least the following resources:
- [x] deployment
- [x] service
- [x] persistent volume claim
- [] garbage collect cron job
- [] ingress
- [] secret (if needed). 
- [] configmap (added)
- [] self-signed ssl (added)

## Preparation - Setup local K8s env by Minikube

To simplify the installation steps, we will use mini

```bash
# Download minikube
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64
sudo install minikube-darwin-amd64 /usr/local/bin/minikube

# Start minikube with virtualbox driver, wait for 5mins
minikube start --driver=virtualbox

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

## Setup basic minikube monitoring

```bash
❯ minikube dashboard --url
🤔  Verifying dashboard health ...
🚀  Launching proxy ...
🤔  Verifying proxy health ...
http://127.0.0.1:52948/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/

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
# Enable minikube tunnel
minikube tunnel # Run in a seperate terminal

❯ k get svc -o wide
NAME              TYPE           CLUSTER-IP      EXTERNAL-IP     PORT(S)          AGE     SELECTOR
docker-registry   LoadBalancer   10.103.70.188   10.103.70.188   8080:31231/TCP   4m24s   app=docker-registry
kubernetes        ClusterIP      10.96.0.1       <none>          443/TCP          6m25s   <none>
❯ docker image tag ubuntu 10.103.70.188:8080/myfirstimage
❯ docker push 10.103.70.188:8080/myfirstimage
Using default tag: latest
The push refers to repository [10.103.70.188:8080/myfirstimage]
Get "https://10.103.70.188:8080/v2/": http: server gave HTTP response to HTTPS client

```

```bash
# Access to the docker registry container
k exec -it docker-registry-deployment-54d89b54b-g4tnl -- sh 
```

Check default config file of registry service
```yaml
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

## PV & PVC

Notes

```
In a production cluster, you would not use hostPath. Instead a cluster administrator would provision a network resource like a Google Compute Engine persistent disk, an NFS share, or an Amazon Elastic Block Store volume. Cluster administrators can also use StorageClasses to set up dynamic provisioning.
```

But in this case, we will use HostPath PV to simplify the process

Create PV, PVC
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: docker-registry-pv
  labels:
    type: local
    workload: docker-registry
spec:
  storageClassName: slow
  accessModes:
    - ReadWriteOnce
  capacity:
    storage: 10Gi
  hostPath:
    path: /data/docker-registry/

apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: docker-registry-pvc
spec:
  storageClassName: slow
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi # Request 5/10GB, can expand later
  selector:
    matchLabels:
      workload: "docker-registry" # Specifically set the PVC to this PV, instead of allowing k8s to auto-allocate
```

Apply the changes
```bash
k apply -f docker-registry/pv.yml   
k apply -f docker-registry/pvc.yml   
```

Mount PVC to the deployment 

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: docker-registry-deployment
  labels:
    app: docker-registry
spec:
  replicas: 1
  selector:
    matchLabels:
        app: docker-registry
  template:
    metadata:
      labels:
        app: docker-registry
    spec:
      containers:
      - name: registry
        image: registry:2
        ports:
        - containerPort: 5000
        volumeMounts:
          - name: docker-registry-storage
            mountPath: /var/lib/registry
      volumes:
        - name: docker-registry-storage
          persistentVolumeClaim:
            claimName: docker-registry-pvc
```

Apply the changes

```bash
k apply -f docker-registry/deployment.yml
```

## Redis configuration + Replica set


# Encountered issues
1. Minikube does not expose the service port correctly by default on Mac. It always return 127.0.0.1 as external IP. We can force it to use virtualbox driver to fix this. Source: https://github.com/kubernetes/minikube/issues/7344

# Further improvements

Use [S3](https://docs.docker.com/registry/storage-drivers/) as backend driver for Docker registry

Use dragonfly as daemon set on each node to speed up image distribution time & offload traffic to docker-registry

# TODO
1. Add yaml lint checker
2. Add Github actions


# Reference

Reference documents for the assignment

https://docs.docker.com/registry/

https://docs.docker.com/registry/storage-drivers/

https://hub.docker.com/_/registry


