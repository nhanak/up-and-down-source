# Have minikube build everything on its own Docker daemon
[string]$minikubeCMD = 'minikube docker-env'
[string]$minikubeCMD2 = '& minikube -p minikube docker-env | Invoke-Expression'
[string]$buildClientCMD = 'docker build -t client ../../client'
[string]$buildDirectorCMD = 'docker build -t director ../../director'
[string]$buildMatchFunctionCMD = 'docker build -t matchfunction ../../matchfunction'
[string]$buildGameServerCMD = 'docker build -t game-server ../../game-server'

iex $minikubeCMD
iex $minikubeCMD2
iex $buildClientCMD
iex $buildGameServerCMD
iex $buildDirectorCMD
iex $buildMatchFunctionCMD