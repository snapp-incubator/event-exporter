fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

tidy:
	go mod tidy

test: fmt vet tidy # envtest ## Run tests.
	golangci-lint run
	go test ./... -covermode=atomic -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

##@ Build

build: test ## Build manager binary.
	go build -race -o bin/main main.go

run: manifests generate fmt vet ## Run a controller from your host.
	go run ./main.go

docker-build: test ## Build docker image with the manager.
	sudo podman build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	sudo podman push ${IMG}

docker-login:
	sudo podman login ${REG} -u ${REG_USER} -p ${REG_PASSWORD}

redeploy: docker-build docker-login docker-push
