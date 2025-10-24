.PHONY: all
all: build
FORCE: ;

.PHONY: build

build: build-auth build-feedback build-vote

build-deposit:
	cd microservices/deposit; go build -o bin/deposit main.go

build-deposit-docker: build-deposit-linux
	cd microservices/deposit; docker build -t eminetto/deposit -f Dockerfile .

build-deposit-linux:
	cd microservices/deposit; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/deposit main.go

build-withdraw:
	cd microservices/withdraw; go build -o bin/withdraw main.go

build-withdraw-docker: build-withdraw-linux
	cd microservices/withdraw; docker build -t withdraw -f Dockerfile .

build-withdraw-linux:
	cd microservices/withdraw; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/withdraw main.go

build-refund:
	cd microservices/refund; go build -o bin/refund main.go

build-refund-docker: build-refund-linux
	cd microservices/refund; docker build -t refund -f Dockerfile .

build-refund-linux:
	cd microservices/refund; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/refund main.go

run-docker: build-deposit-docker build-withdraw-docker build-refund-docker
    docker run -d -p 8080:8080 eminetto/deposit
    docker run -d -p 8081:8081 withdraw
    docker run -d -p 8082:8082 refund

deploy-k8s: #build-deposit-docker build-withdraw-docker build-refund-docker
	docker push eminetto/deposit:latest
	kubectl create namespace deposit
	kubectl apply --namespace deposit -f microservices/deposit/deposit.yaml
	kubectl port-forward --namespace deposit deployment/deposit 8080:8080

