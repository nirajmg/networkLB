kconfig :=  $(shell ls ~/.kube/config)

all: clean build docker helm

build:
	rm -rf nlb/bin
	mkdir nlb/bin
	cd ./nlb;pwd; GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o bin/nlb main.go

docker:
	cd nlb; docker build -t nlb:latest .

helm:
	helm install postgresql infra/postgresql
	helm install nlb infra/nlb  --set kubeconfig=$(shell ls ~/.kube/config)

clean:
	rm -rf nlb/bin
	helm delete postgresql
	helm delete nlb

