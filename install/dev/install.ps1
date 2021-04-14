#---------------------------------------------------
# install.ps1: Installs Up & Down on local minikube
#---------------------------------------------------
# Let shell know we want to use minikubes docker environment
[string]$minikubeCMD = 'minikube docker-env'
[string]$minikubeCMD2 = '& minikube -p minikube docker-env | Invoke-Expression'
iex $minikubeCMD
iex $minikubeCMD2

# Install Open Match core
# https://open-match.dev/site/docs/installation/yaml/#install-core-open-match
[string]$installOpenMatchCoreCMD = 'kubectl apply --namespace open-match -f https://open-match.dev/install/v1.0.0/yaml/01-open-match-core.yaml'
iex $installOpenMatchCoreCMD

# Delete Open Match

#kubectl delete psp,clusterrole,clusterrolebinding --selector=release=open-match
#kubectl delete namespace open-match

# Install Evaluator
[string]$installEvaluatorCMD = 'kubectl apply --namespace open-match -f ../yaml/06-open-match-override-configmap.yaml -f ../yaml/07-open-match-default-evaluator.yaml'
iex $installEvaluatorCMD

# Install Agones
[string]$setupAgonesNameSpaceCMD = 'kubectl create namespace agones-system'
[string]$installAgonesCMD = 'kubectl apply -f https://raw.githubusercontent.com/googleforgames/agones/release-1.6.0/install/yaml/install.yaml'
iex $setupAgonesNameSpaceCMD
iex $installAgonesCMD

# Install Up & Down
[string]$installUpAndDownCMD = 'kubectl apply -f ./up_and_down.yaml'
iex $installUpAndDownCMD

# Set default user on game-server-manager to be allow to read and create pods
[string]$clusterrolebindingCMD =  'kubectl create clusterrolebinding game-server-manager --clusterrole=gameservers-manager --serviceaccount=default:default'
iex $clusterrolebindingCMD

# Install Up & Down Fleet
[string]$installUpAndDownFleetCMD = 'kubectl apply -f ../yaml/fleet.yaml'
iex $installUpAndDownFleetCMD

# Install Up & Down Fleet Auto Scaler
[string]$installUpAndDownFleetCMD = 'kubectl apply -f ../yaml/fleet-autoscaler.yaml'
iex $installUpAndDownFleetCMD
