#kconfig :=  $(shell ls ~/.kube/config)
kconfig := "/mnt/c/Users/hanna/Desktop/Coursework/Grad/NetworkSystems/networkLB/config"
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
	helm install server-1 infra/server --set resources.limits.memory=128Mi --set resources.requests.memory=128Mi
	helm install server-2 infra/server --set resources.limits.memory=256Mi --set resources.requests.memory=256Mi
	helm install server-3 infra/server --set resources.limits.memory=450Mi --set resources.requests.memory=450Mi
	helm install server-4 infra/server --set resources.limits.memory=600Mi --set resources.requests.memory=600Mi 

clean:
	rm -rf nlb/bin
	helm delete postgresql
	helm delete nlb
	helm delete server-1
	helm delete server-2
	helm delete server-3
	helm delete server-4


