postgres:
	docker run --name postgres15 -p 5432:5432 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

restart:
	docker restart postgres15

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1


destroy:
	docker stop postgres15 
	docker rm postgres15

sqlc:
	sqlc generate 

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination ./db/mock/store.go github.com/ngyale-pro/simplebank/db/sqlc Store

kubedeploy:
	kubectl apply -f eks/deployment.yaml
	kubectl apply -f eks/service.yaml
	kubectl apply -f eks/ingress.yaml
	kubectl apply -f eks/issuer.yaml

kubedelete:
	kubectl delete clusterissuer letsencrypt-staging
	kubectl delete ingressclass nginx 
	kubectl delete ingress simple-bank-ingress 
	kubectl delete service simple-bank-api-service 
	kubectl delete deployment simple-bank-api-deployment 

.PHONY: postgres createdb dropdb migrateup migratedown destroy sqlc test server mock migrateup1 migratedown1 kubedeploy kubedelete