ARG GO_IMG

FROM $GO_IMG
ARG CWD
ENV GOFLAGS=
ENV LOG_LEVEL=debug

RUN apk add --update curl alpine-sdk gcc

WORKDIR $CWD
COPY . .
RUN go build -v -o /app .

CMD ["/app"]

HEALTHCHECK --interval=16s --timeout=2s \
    CMD curl --fail http://127.0.0.1/healthcheck
