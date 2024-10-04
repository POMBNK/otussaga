start:
	cd k8s/ \
	&& helm install database database \
	&& helm install rabbitmq rabbitmq \
	&& helm install orderservice orderservice \
	&& helm install paymentservice paymentservice \
	&& helm install warehouseservice warehouseservice \
	&& minikube tunnel
.PHONY: start