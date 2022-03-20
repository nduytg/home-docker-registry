# home-docker-registry
Build your own Docker Registry on K8s from scratch

Author: nduytg@gmail.com

Github: [github.com/nduytg](github.com/nduytg)

## Quick start

```bash
# Start minikube with virtualbox driver, wait for 5mins
minikube start --driver=virtualbox

k apply -f redis/
k apply -f docker-registry/
```

## Requirements

Write the Kubernetes deployment manifest to run Docker Registry in Kubernetes with at least the following resources:
- [x] deployment
- [x] service
- [x] persistent volume claim
- [x] garbage collect cron job
- [x] ingress
- [x] secret (if needed). 
- [x] Redis
- [x] Increase replicas

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

^ Does not work because we dont have SSL certs yet! Let's add one!

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

## Create self-signed cert

```bash
# Create certs
openssl req -newkey rsa:4096 -nodes -keyout ./certs/registry.key -x509 -days 365 -out ./certs/registry.crt

# Create secrets
k create secret generic registry-certs-keys --from-file=./certs/registry.crt --from-file=./certs/registry.key 

# Create basic auth file
docker run \
  --entrypoint htpasswd \
  httpd:2 -Bbn duyuser password123456 > auth/htpasswd

# Create secret for auth file
k create secret generic registry-htpasswd --from-file=./auth/htpasswd

```

## Enable Ingress controller for minikube

In our case, we use local minikube env, so we need to enable ingress add-on for minikube

```bash
# Enable add-on
minikube addons enable ingress

# Check if ingress is running
kubectl get pods -n ingress-nginx
```

Ok, now we need to define the ruleset for our newly created Ingress

Create the Ingress object by running the following command:

```bash
k apply -f docker-registry/ingress-registry.yaml

#The output should be:
ingress.networking.k8s.io/ingress-registry created

```
Verify the IP address is set:

```bash
❯ kubectl get ingress
NAME              CLASS   HOSTS             ADDRESS        PORTS     AGE
docker-registry   nginx   registry.duy.io   192.168.64.3   80, 443   13h

# Set this in your local /etc/hosts file
192.168.64.3 registry.duy.io 

# Access the domain from your laptop browser
registry.duy.io
```

We will encounter invalid cert error if we use the Ingress cert out of the box

```bash
❯ docker tag ubuntu:16.04 registry.duy.io/ubuntu:16.04
❯ docker push registry.duy.io/ubuntu:16.04
The push refers to repository [registry.duy.io/ubuntu]
Get "https://registry.duy.io/v2/": x509: certificate is valid for ingress.local, not r
egistry.duy.io
❯ docker push registry.duy.io/ubuntu:16.04
The push refers to repository [registry.duy.io/ubuntu]
Get "https://registry.duy.io/v2/": x509: certificate is valid for ingress.local, not registry.duy.io
```

We need to create new certs (the same way we did with the cert in registry deployment), but this time for ingress

Let's go!!

```
# Create cert by openssl
openssl req -newkey rsa:4096 -nodes -keyout ./ingress-certs/registry-ingress.key -x509 -days 365 -out ./ingress-certs/registry-ingress.crt -subj "/CN=registry.duy.io/O=registry.duy.io"

# Create new secret
kubectl create secret tls registry-ingress --key ./ingress-certs/registry-ingress.key --cert ./ingress-certs/registry-ingress.crt
```

Replace default secret in ingress config with new secret we just created
```yaml
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - registry.duy.io
      secretName: registry-ingress
```

Try to run docker push command again, we encounter this error
```
❯ docker push registry.duy.io/ubuntu:16.04
The push refers to repository [registry.duy.io/ubuntu]
Get "https://registry.duy.io/v2/": x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0
```

Need to import the cert into Mac Keychain Access

Also need to add the following settings into your docker daemon config => Reload docker daemon

```json
  "insecure-registries" : ["registry.duy.io"],
```


Retry => Profit!!

```bash
# Before
❯ docker push registry.duy.io/ubuntu:16.04
The push refers to repository [registry.duy.io/ubuntu]
Get "https://registry.duy.io/v2/": x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0

# After ;) (remember to relogin to new domain before you run this though)
❯ docker login -u duyuser -p password123456 registry.duy.io
❯ docker push registry.duy.io/ubuntu:16.04
The push refers to repository [registry.duy.io/ubuntu]
1251204ef8fc: Pushed 
47ef83afae74: Pushed 
df54c846128d: Pushed 
be96a3f634de: Pushed 
16.04: digest: sha256:a3785f78ab8547ae2710c89e627783cfa7ee7824d3468cae6835c9f4eae23ff7 size: 1150
```

## Redis configuration + Replica set

k apply -f redis/

```bash
❯ k exec -it redis -- redis-cli
127.0.0.1:6379> CONFIG GET maxmemory
1) "maxmemory"
2) "209715200"
127.0.0.1:6379> 
```

Important notes for Redis configuration
* Need to set a HTTP secret
* Set storage.cache.blobdescriptor: redis
* Set redis addr in registry configuration

```bash
k apply -f docker-registry/registry-config.yml
k apply -f docker-registry/deployment.yml
```

Ideally, we will see something like this

```bash
# In redis
127.0.0.1:6379> keys *
 1) "repository::redis::blobs::sha256:04ab1bfc453f19989c401c2f0622df3363b5182bdea1af9e59ee2ceea3a9931c"
 2) "blobs::sha256:961b8e95c0f4561047ea48e53e564dae6a4e4b5d3334579a0344af2b04ecb3f9"
 3) "repository::grafana::blobs::sha256:dae972374a52168fce7ad44c1f169b6ae65ffa5e1c43d93b78abffa43207925e"
 4) "blobs::sha256:29b14abb751a5802c3b7174e21ed6aa8486627788969e4214a2c90147fef056f"
 5) "blobs::sha256:cdd789ccb9ea8c941d008916c02350057379875d56187e95a0d9ee823d3e2f6f"
 6) "blobs::sha256:54fec2fa59d0a0de9cd2dec9850b36c43de451f1fd1c0a5bf8f1cf26a61a5da4"
 7) "repository::alpine::blobs::sha256:3d243047344378e9b7136d552d48feb7ea8b6fe14ce0990e0cc011d5e369626a"
....
```

Registry container logs
```bash
# On pod 1
172.17.0.2 - - [20/Mar/2022:06:54:33 +0000] "HEAD /v2/grafana/blobs/sha256:cdd789ccb9ea8c941d008916c02350057379875d56187e95a0d9ee823d3e2f6f HTTP/1.1" 200 0 "" "docker/20.10.8 go/go1.16.6 git-commit/75249d8 kernel/5.10.47-linuxkit os/linux arch/amd64 UpstreamClient(Docker-Client/20.10.8 \\(darwin\\))"
172.17.0.2 - - [20/Mar/2022:06:54:33 +0000] "PUT /v2/grafana/manifests/latest HTTP/1.1" 201 0 "" "docker/20.10.8 go/go1.16.6 git-commit/75249d8 kernel/5.10.47-linuxkit os/linux arch/amd64 UpstreamClient(Docker-Client/20.10.8 \\(darwin\\))"
time="2022-03-20T06:54:33.490089529Z" level=info msg="response completed" go.version=go1.16.15 http.request.contenttype="application/vnd.docker.distribution.manifest.v2+json" http.request.host=registry.duy.io http.request.id=18977836-a222-4910-b455-f94516c89538 http.request.method=PUT http.request.remoteaddr=192.168.64.1 http.request.uri="/v2/grafana/manifests/latest" http.request.useragent="docker/20.10.8 go/go1.16.6 git-commit/75249d8 kernel/5.10.47-linuxkit os/linux arch/amd64 UpstreamClient(Docker-Client/20.10.8 \(darwin\))" http.response.duration=16.337273ms http.response.status=201 http.response.written=0


# On pod 2
172.17.0.2 - - [20/Mar/2022:06:54:33 +0000] "HEAD /v2/grafana/blobs/sha256:cdd789ccb9ea8c941d008916c02350057379875d56187e95a0d9ee823d3e2f6f HTTP/1.1" 404 157 "" "docker/20.10.8 go/go1.16.6 git-commit/75249d8 kernel/5.10.47-linuxkit os/linux arch/amd64 UpstreamClient(Docker-Client/20.10.8 \\(darwin\\))"
172.17.0.2 - - [20/Mar/2022:06:54:33 +0000] "PATCH /v2/grafana/blobs/uploads/4ce77739-d321-45be-b0b1-4949208c0743?_state=3ET9AKLzACngAL8hRjVGHMtQYKZA6eW_LmUroIL_yIV7Ik5hbWUiOiJncmFmYW5hIiwiVVVJRCI6IjRjZTc3NzM5LWQzMjEtNDViZS1iMGIxLTQ5NDkyMDhjMDc0MyIsIk9mZnNldCI6MCwiU3RhcnRlZEF0IjoiMjAyMi0wMy0yMFQwNjo1NDozMy4zNzA4NzM4MDlaIn0%3D HTTP/1.1" 202 0 "" "docker/20.10.8 go/go1.16.6 git-commit/75249d8 kernel/5.10.47-linuxkit os/linux arch/amd64 UpstreamClient(Docker-Client/20.10.8 \\(darwin\\))"
time="2022-03-20T06:54:33.404872485Z" level=info msg="response completed" go.version=go1.16.15 http.request.host=registry.duy.io http.request.id=667435cf-420d-4c04-80fc-f9e211673b31 http.request.method=PATCH http.request.remoteaddr=192.168.64.1 http.request.uri="/v2/grafana/blobs/uploads/4ce77739-d321-45be-b0b1-4949208c0743?_state=3ET9AKLzACngAL8hRjVGHMtQYKZA6eW_LmUroIL_yIV7Ik5hbWUiOiJncmFmYW5hIiwiVVVJRCI6IjRjZTc3NzM5LWQzMjEtNDViZS1iMGIxLTQ5NDkyMDhjMDc0MyIsIk9mZnNldCI6MCwiU3RhcnRlZEF0IjoiMjAyMi0wMy0yMFQwNjo1NDozMy4zNzA4NzM4MDlaIn0%3D" http.request.useragent="docker/20.10.8 go/go1.16.6 git-commit/75249d8 kernel/5.10.47-linuxkit os/linux arch/amd64 UpstreamClient(Docker-Client/20.10.8 \(darwin\))" http.response.duration=5.384046ms http.response.status=202 http.response.written=0
```

# Encountered issues
1. Minikube does not expose the service port correctly by default on Mac. It always return 127.0.0.1 as external IP. We can force it to use virtualbox driver to fix this. Source: https://github.com/kubernetes/minikube/issues/7344
2. If the registry is empty, the GC cronjob will fail!! Because there is no docker directory in that volume yet!

# Further improvements

0. Can change to use NFS + PVC ReadWriteMany when we have multiple nodes

1. Use [S3](https://docs.docker.com/registry/storage-drivers/) as backend driver for Docker registry

2. Use dragonfly as daemon set on each node to speed up image distribution time & offload traffic to docker-registry

3. Set ACL for ingress + docker registry to only accept internal traffic

4. Enable auto-scaling by prometheus + KEDA

5. Change the Auth method to Token, instead of htpasswd

6. Enable Docker proxy feature, cache dockerhub image on local to avoid the Dockerhub rate limit issue

# Reference

Reference documents for the assignment

https://docs.docker.com/registry/

https://docs.docker.com/registry/storage-drivers/

https://hub.docker.com/_/registry

https://docs.docker.com/registry/garbage-collection/

https://github.com/marketplace/actions/yamllint-github-action

https://kubernetes.github.io/ingress-nginx/examples/docker-registry/

https://kubernetes.io/docs/tasks/access-application-cluster/ingress-minikube/
