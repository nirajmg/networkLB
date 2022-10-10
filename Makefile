all: database

database:
	helm install postgresql infra/postgresql

clean:
	helm delete postgresql