KEY_PATH ?= ~/.ssh/serhii-aws-key.pem
TF_DIR = infra/tofu
AN_DIR = infra/ansible

export ANSIBLE_HOST_KEY_CHECKING=False

start-local-compose:
	docker compose up --build
clean-compose:
	docker compose down
start-local-k3d-helmfile:
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
start-local-k3d-argo:
	@echo "Start a k3d cluster"
	k3d cluster create local-ha --servers 1 --agents 2 --port "8080:8080@loadbalancer"

	@echo "Make a local image"
	docker build -t ghcr.io/zapi-web/url-shortener:local .
	k3d image import ghcr.io/zapi-web/url-shortener:local -c local-ha

	@echo "Install argo"
	helm upgrade --install argocd argo-cd --repo https://argoproj.github.io/argo-helm --namespace argo-cd --create-namespace -f deploy/platform/values/argo-cd.yaml

	@echo "Install prometheus CRD"
	kubectl apply --server-side -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_servicemonitors.yaml && kubectl apply --server-side -f https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_podmonitors.yaml

	@echo "Deploy platform"
	cd deploy/platform && kubectl apply -f argo/applicationSet.yaml && kubectl apply -f argo/monitoringDashboard.yaml

	@echo "Deploy App"
	cd deploy/platform && kubectl apply -f argo/application.yaml
clean-k3d:
	k3d cluster delete local-ha
upgrade-local-k3d:
	@echo "Deploy platform"
	cd deploy/platform && helmfile sync
	kubectl apply -f deploy/platform/manifests/otel-collector.yaml

	@echo "Deploy App"
	helm upgrade --install url-shortener ./deploy/app -f ./deploy/app/values.yaml --create-namespace --namespace app
start-aws-k3s-cluster:
	@echo "Make an infrastructure with OpenTofu"
	tofu -chdir=$(TF_DIR) apply -auto-approve

	@echo "Ping all servers"
	cd $(AN_DIR) && ansible all -m ping --private-key=$(KEY_PATH)

	@echo "Configure servers, and deploy"
	cd $(AN_DIR) && ansible-playbook site.yaml --private-key=$(KEY_PATH)
aws-k3s-cluster-destroy:
	tofu -chdir=$(TF_DIR) destroy -auto-approve