start-local-compose:
	docker compose up --build
clean-compose:
	docker compose down
start-local-k3d:
	@echo "Start a k3d cluster"
	k3d cluster create local-ha --servers 1 --agents 2 --port "8080:8080@loadbalancer"

	@echo "Make a local image"
	docker build -t ghcr.io/zapi-web/url-shortener:local .
	k3d image import ghcr.io/zapi-web/url-shortener:local -c local-ha

	@echo "Install prometheus CRD"
	kubectl apply --server-side -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_servicemonitors.yaml && kubectl apply --server-side -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_podmonitors.yaml

	@echo "Deploy platform"
	cd deploy/platform && helmfile sync
	kubectl apply -f deploy/platform/manifests/otel-collector.yaml

	@echo "Deploy App"
	helm upgrade --install url-shortener ./deploy/app -f ./deploy/app/values.yaml --create-namespace --namespace app
clean-k3d:
	k3d cluster delete local-ha
upgrade-local-k3d:
	@echo "Deploy platform"
	cd deploy/platform && helmfile sync
	kubectl apply -f deploy/platform/manifests/otel-collector.yaml

	@echo "Deploy App"
	helm upgrade --install url-shortener ./deploy/app -f ./deploy/app/values.yaml --create-namespace --namespace app