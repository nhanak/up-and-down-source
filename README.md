# up & down
1v1 web game where players control shapes and wage war against eachother

## How to Run Locally (Recommended)
### 1. Run the Client UI
Inside ./client/frontend/src/pages/Home.js, there is this.local inside the constructor. Set this.local to true.
If you have just cloned this repository, this has already be done for you.
```
cd ./client/frontend
yarn install
yarn start
```

### 2. Run the GameServer
Inside ./game-server/main.go  there is const local at the top of the file. Set const local = true.
If you have just cloned this repository, this has already be done for you.
```
cd ./game-server
go run .
```
You might also want to build the executable if you are going to be restarting the local game-server multiple times
```
go build -o ./game-server.exe
./game-server.exe
```

### 3. Play the game in Browser
1. Open up two instances of your favorite browser (Chrome, Firefox, etc.) and navigate to localhost:3000. 
2. Enter a Nickname and click play on each instance. 
3. Enjoy the game! 
#### Notes 
I do not recommend having just one browser with seperate tabs running the game, as this does not work 100% of the time. 

## How to Run on Google Cloud
Take a look inside the ./install and ./install/prod folders for instructions.
#### Notes
I have removed the custom certs that were used to get https working, so you will have to build your own or find a workaround that doesn't require certs.

## How to Build and Push to Google Cloud Registry (Powershell)
### 1. Setup REGISTRY Environment Variable
```
$Env:REGISTRY="gcr.io/[YOUR_REGISTRY_URL]"
```

### 2. Build & Push 
#### a) Client
```
docker build -t $env:REGISTRY/client ./client/
docker push $env:REGISTRY/client
```
#### b) Director
```
docker build -t $env:REGISTRY/director ./director/
docker push $env:REGISTRY/director
```

#### c) Matchfunction
```
docker build -t $env:REGISTRY/matchfunction ./matchfunction/
docker push $env:REGISTRY/matchfunction
```

#### d) GameServer
```
docker build -t $env:REGISTRY/game-server ./game-server/
docker push $env:REGISTRY/game-server
```

#### e) Cloud MySQL
```
docker build -t $env:REGISTRY/cloud-mysql ./cloud-mysql/
docker push $env:REGISTRY/cloud-mysql
```

## How to Delete Pods, Services, Deployments
### Delete all Pods
```
kubectl delete --all pods --namespace=default
```

### Delete all Deployments
```
 kubectl delete --all deployments --namespace=default
```

### Delete all Services
```
 kubectl delete --all services --namespace=default
```

### Nuke Docker
```
docker system prune -a
```

### Nuke all gameservers
```
kubectl delete gs --all
```

## Tech Stack Overview
1. Go backend
2. React frontend
3. MySQL leaderboards/mmr
4. Kubernetes & Docker for scalability

### Why Kubernetes 1.17.0? 
Agones doesn't work for Kubernetes 1.18.0 yet

## Resources 

### How was the matchmaker made?
https://open-match.dev/site/
