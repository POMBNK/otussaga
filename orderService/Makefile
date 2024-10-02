generate:
	go generate ./...
	go run docs/merger/main.go
	statik -src=docs/dist -include=*.html,*.css,*.js,*.png,*.json

.PHONY: generate

start:
	cd deployment/k8s/ && helm install auth auth && minikube tunnel

.PHONY: start

test:
	cd http_tests && newman run auth-collection.postman_collection.json --insecure

stop:
	cd deployment/k8s/ && helm delete auth && kubectl delete jobs auth-job

.PHONY: stop
