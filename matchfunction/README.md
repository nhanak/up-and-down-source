This folder provides a sample Match Function for Open Match Matchmaker 101 Tutorial.

Run the below steps in this folder to set up the Match Function.

Step 1: Specify your Registry URL.
```
REGISTRY=[YOUR_REGISTRY_URL]
```

Step 2: Build the Match Function image.
```
docker build -t $REGISTRY/mm101-tutorial-matchfunction .
```

Step 3: Push the Match Function image to the configured Registry.
```
docker push $REGISTRY/mm101-tutorial-matchfunction
```

Step 4: Update the install yaml for your setup.
```
sed "s|REGISTRY_PLACEHOLDER|$REGISTRY|g" matchfunction.yaml | kubectl apply -f -
```
