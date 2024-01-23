cluster=prod
KUBECONFIG=$(HOME)/.kube/$(cluster)
image=paskalmaksim/db-table-metrics:dev

test:
	./scripts/validate-license.sh
	go fmt ./cmd ./pkg/...
	go vet ./cmd ./pkg/...
	go test ./pkg/...
	go mod tidy
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v
run:
	go run --race ./cmd \
	-grace.interval=0 \
	-web.listen-address=127.0.0.1:8080 \
	-log.level=debug
build:
	go run github.com/goreleaser/goreleaser@latest build --clean --skip=validate --snapshot
	mv ./dist/db-table-metrics_linux_amd64_v1/db-table-metrics .
	docker buildx build --pull --push . -t $(image)
deploy:
	helm uninstall db-table-metrics --namespace db-table-metrics || true
	helm upgrade db-table-metrics ./chart/db-table-metrics \
	--install \
	--create-namespace \
	--namespace db-table-metrics \
	--values values.$(cluster).yaml
deploy-all:
	make deploy cluster=prod
	make deploy cluster=azure-prod