REPO = github.com/imega/stock-miner
CWD = /go/src/$(REPO)
GO_IMG = golang:1.15.8-alpine3.13
NODE_IMG = node:16.1.0-alpine3.13
BUILDER=builder

test: build unit acceptance

builder:
	@docker build --build-arg GO_IMG=$(GO_IMG) \
		-t builder -f $(CURDIR)/tests/Dockerfile .
	@touch builder

build: node_modules
	@docker run --rm -v $(CURDIR):/data -w /data \
		-e TAG=$(TAG) \
		-e VERSION=$(GITHUB_REF) \
		-e GRAPHQL_HOST=$(GRAPHQL_HOST) \
		-e GRAPHQL_SCHEMA=$(GRAPHQL_SCHEMA) \
		$(NODE_IMG) \
		sh -c "npm run build"

node_modules:
	@docker run --rm -v $(CURDIR):/data -w /data $(NODE_IMG) npm install

lint:
	@-docker run --rm -t -v $(CURDIR):$(CWD) -w $(CWD) \
		golangci/golangci-lint golangci-lint run

unit: builder
	@docker run --rm -w $(CWD) -v $(CURDIR):$(CWD) \
		$(BUILDER) sh -c "\
			go list ./... | grep -v 'tests' | xargs go test -vet=off -coverprofile cover.out \
		"

acceptance: builder down
	GO_IMG=$(BUILDER) CWD=$(CWD) docker-compose up -d --build --scale acceptance=0
	GO_IMG=$(BUILDER) CWD=$(CWD) docker-compose up --abort-on-container-exit acceptance

b:
	GO_IMG=$(BUILDER) CWD=$(CWD) docker-compose up -d --build --scale acceptance=0

a:
	GO_IMG=$(BUILDER) CWD=$(CWD) docker-compose up --abort-on-container-exit acceptance

down:
	GO_IMG=$(BUILDER) CWD=$(CWD) docker-compose down -v --remove-orphans

release: build
	go run -tags=dev assets/generate.go
