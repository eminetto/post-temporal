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

build-client:
	cd cmd/client; go build -o bin/client main.go

build-client-docker: build-client-linux
	docker build -t client -f cmd/client/Dockerfile .

build-client-linux:
	cd cmd/client; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/client main.go

build-worker:
	cd cmd/worker; go build -o bin/worker main.go

build-worker-docker: build-worker-linux
	docker build -t worker -f cmd/worker/Dockerfile .

build-worker-linux:
	cd cmd/worker; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "netgo" -installsuffix netgo -o bin/worker main.go


run-docker: build-deposit-docker build-withdraw-docker build-refund-docker
    docker run -d -p 8080:8080 eminetto/deposit
    docker run -d -p 8081:8081 withdraw
    docker run -d -p 8082:8082 refund
    docker run -d -p 8083:8083 client
    docker run -d worker

deploy-k8s: #build-deposit-docker build-withdraw-docker build-refund-docker
	helm install \
        --repo https://go.temporal.io/helm-charts \
        --set server.replicaCount=1 \
        --set cassandra.config.cluster_size=1 \
        --set elasticsearch.replicas=1 \
        --set prometheus.enabled=false \
        --set grafana.enabled=false \
        temporaltest temporal \
        --timeout 15m

	docker push eminetto/deposit:latest
	kubectl create namespace deposit
	kubectl apply --namespace deposit -f microservices/deposit/deposit.yaml
	kubectl port-forward --namespace deposit deployment/deposit 8080:8080

	docker push eminetto/refund:latest
	kubectl create namespace refund
	kubectl apply --namespace refund -f microservices/refund/refund.yaml
	kubectl port-forward --namespace refund deployment/refund 8082:8082

	docker push eminetto/withdraw:latest
	kubectl create namespace withdraw
	kubectl apply --namespace withdraw -f microservices/withdraw/withdraw.yaml
	kubectl port-forward --namespace withdraw deployment/withdraw 8081:8081

	docker push eminetto/client:latest
	kubectl create namespace client
	kubectl apply --namespace client -f cmd/client/client.yaml
	kubectl port-forward --namespace client deployment/client 8083:8083

	docker push eminetto/worker:latest
	kubectl create namespace worker
	kubectl apply --namespace worker -f cmd/worker/worker.yaml

clean-k8s:
	kubectl delete namespace deposit
	kubectl delete namespace withdraw
	kubectl delete namespace refund
	helm uninstall temporaltest
