REPO = github.com/imega/stock-miner
CWD = /go/src/$(REPO)
GO_IMG = golang:1.15.8-alpine3.13
BUILDER=builder

test: lint unit acceptance

builder:
	@docker build --build-arg GO_IMG=$(GO_IMG) \
		-t builder -f $(CURDIR)/tests/Dockerfile .
	@touch builder

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

down:
	GO_IMG=$(BUILDER) CWD=$(CWD) docker-compose down -v --remove-orphans
