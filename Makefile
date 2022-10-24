all: database

build:
	mkdir nlb/bin
	go build -o nlb/bin/nlb nlb/main.go


database:
	helm install postgresql infra/postgresql

clean:
	helm delete postgresql