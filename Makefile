kconfig :=  $(shell ls ~/.kube/config)

all:  build docker helm

build:	
	mkdir -p nlb/bin;cd ./nlb;pwd; GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o bin/nlb main.go
	mkdir -p server/bin; cd ./server;pwd; GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o bin/server main.go

docker:
	cd nlb; docker build -t nlb:latest .
	cd server; docker build -t server:latest .

helm:
	helm install postgresql infra/postgresql
	helm install nlb infra/nlb  --set kubeconfig=$(shell ls ~/.kube/config)
	helm install server infra/server 

clean:
	rm -rf nlb/bin
	helm delete postgresql
	helm delete nlb
	helm delete server

