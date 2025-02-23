ARG ALIPNE_VERSION=3.21
ARG GO_VERSION=1.24.0

FROM golang:${GO_VERSION}-alpine${ALIPNE_VERSION} AS builder
WORKDIR /app

RUN apk --no-cache add make git build-base

RUN --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

COPY . .
RUN make build

FROM alpine:${ALIPNE_VERSION}
WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/config/config.yaml config/config.yaml
COPY --from=builder /app/.env .

CMD ["/app/api"]
