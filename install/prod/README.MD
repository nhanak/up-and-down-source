# INSTALLING PRODUCTION
## Make kubectl use the correct context
```
kubectl config use-context gke_expanded-goal-281218_us-central1-c_up-and-down-cluster-1
```

## Make sure the firwall is open
```
gcloud compute firewall-rules create game-server-firewall-tcp --allow tcp:7000-8000 --target-tags game-server --description "Firewall to allow game server udp traffic"
gcloud compute firewall-rules create game-server-firewall-tcp --allow tcp:7000-8000 --description "Firewall to allow game server tcp traffic"
```

### Ping
```
Test-NetConnection -ComputerName 34.68.194.98 -Port 7659
```

### Certificates 
To enable HTTPS, check certificates folder