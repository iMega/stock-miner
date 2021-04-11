ARG GO_IMG

FROM $GO_IMG
ARG CWD
ENV GOFLAGS=
ENV LOG_LEVEL=debug

WORKDIR $CWD
ADD . .

RUN go run -tags=dev assets/generate.go && \
    go build -v -ldflags "-X main.isDevMode=true" -o /app .

CMD ["/app"]

HEALTHCHECK --interval=16s --timeout=2s \
    CMD curl --fail http://127.0.0.1/healthcheck
