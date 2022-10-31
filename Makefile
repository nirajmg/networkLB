all: database

build:
	rm -rf nlb/bin
	mkdir nlb/bin
	cd ./nlb;pwd; GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -o bin/nlb main.go

database:
	helm install postgresql infra/postgresql

clean:
	rm -rf nlb/bin
	helm delete postgresql

