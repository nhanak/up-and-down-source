# Installation

## How to Install Development Build (local, minikube, windows)

### 1. Create a local minikube cluster 
```
minikube start --cpus=3 --memory=2500mb --kubernetes-version v1.17.0
```
#### Wait... why 1.17.0? 
Agones doesn't work for Kubernetes 1.18.0 yet

### 2. Build Local Images
```
./build.ps1
```

### 3. Install
```
./install.ps1
```

# Notes
## You can't pull from Google Cloud Registry (GCR) without the proper keys.
If you intend to run the production images hosted on GCR for your local build, you're going to have to set up your proper keys

Technically, this doesn't seem like a big deal until you realize that Powershell and kubectl do not like eachother. This means it is MUCH EASIER to use git bash if you are on Windows. Better yet, use Linux as your dev environment.

This was a huge pain to figure out, thankfully there is a tutorial.

https://medium.com/hackernoon/today-i-learned-pull-docker-image-from-gcr-google-container-registry-in-any-non-gcp-kubernetes-5f8298f28969
### 1. Get Service Account JSON Key from cloud console
### 2. Run Bash Commands
First create the key
```
 kubectl create secret docker-registry gcr-json-key \
  --docker-server=https://gcr.io \
  --docker-username=_json_key \
  --docker-password="$(cat ~/json-key-file-from-gcp.json)" \
  --docker-email=any@valid.email
```

Then, patch the service account

```
 kubectl patch serviceaccount default \
-p '{"imagePullSecrets": [{"name": "gcr-json-key"}]}'
```

## We're using private repos
... And golang is sketched out by that. Basically, you have to warn it by setting an env variable
```
go env -w GOPRIVATE="github.com/nhanak"
```