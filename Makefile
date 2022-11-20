kconfig :=  $(shell ls ~/.kube/config)
DOCKERUSER=us-west3-docker.pkg.dev/lab5-364722/lab7
PROJECT=lab5-364722
CLUSTER=cluster-1
REGION=us-west3-a

all:  build docker helm

build:	
	mkdir -p nlb/bin;cd ./nlb;pwd; GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o bin/nlb main.go
	mkdir -p server/bin; cd ./server;pwd; GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o bin/server main.go

docker:
	cd nlb; docker build -t nlb:latest .
	cd server; docker build -t server:latest .

helm:
	helm install postgresql infra/postgresql
	helm install nlb infra/nlb  --set kubeconfig=$(kconfig)
	helm install server-1 infra/server --set resources.limits.memory=256Mi --set resources.requests.memory=256Mi
	helm install server-2 infra/server --set resources.limits.memory=256Mi --set resources.requests.memory=256Mi
	helm install server-3 infra/server --set resources.limits.memory=512Mi --set resources.requests.memory=512Mi
	# helm install server-4 infra/server --set resources.limits.memory=1024Mi --set resources.requests.memory=1024Mi 

clean:
	helm delete postgresql
	helm delete nlb
	helm delete server-1
	helm delete server-2
	helm delete server-3


