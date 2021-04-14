This folder provides a sample Director for Open Match Matchmaker 101 Tutorial.

Run the below steps in this folder to set up the Director.

Step 1: Specify your Registry URL.
```
REGISTRY=[YOUR_REGISTRY_URL]
```

Step 2: Build the Director image.
```
docker build -t $REGISTRY/mm101-tutorial-director .
```

Step 3: Push the Director image to the configured Registry.
```
docker push $REGISTRY/mm101-tutorial-director
```

Step 4: Update the install yaml for your setup.
```
sed "s|REGISTRY_PLACEHOLDER|$REGISTRY|g" director.yaml | kubectl apply -f -
```
