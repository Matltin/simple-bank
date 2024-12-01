# simple-bank

Install tools

Docker desktop

TablePlus

Golang

Homebrew

Migrate

brew install golang-migrate
DB Docs

npm install -g dbdocs
dbdocs login
DBML CLI

npm install -g @dbml/cli
dbml2sql --version
Sqlc

brew install sqlc
Gomock

go install github.com/golang/mock/mockgen@v1.6.0
Setup infrastructure
Create the bank-network

make network
Start postgres container:

make postgres
Create simple_bank database:

make createdb
Run db migration up all versions:

make migrateup
Run db migration up 1 version:

make migrateup1
Run db migration down all versions:

make migratedown
Run db migration down 1 version:

make migratedown1
Documentation
Generate DB documentation:

make db_docs
Access the DB documentation at this address. Password: secret

How to generate code
Generate schema SQL file with DBML:

make db_schema
Generate SQL CRUD with sqlc:

make sqlc
Generate DB mock with gomock:

make mock
Create a new db migration:

make new_migration name=<migration_name>
How to run
Run server:

make server
Run test:

make test
Deploy to kubernetes cluster
Install nginx ingress controller:

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/aws/deploy.yaml
Install cert-manager:

kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
