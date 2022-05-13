# BUILD STAGE
FROM golang:1.18.1-alpine3.15 AS builder
WORKDIR /app
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOFLAGS=-buildvcs=false \
    && go mod tidy \
    && go build -ldflags '-w -s' -o trade

FROM alpine:3.15
RUN apk update && \
  apk add ca-certificates && \
  update-ca-certificates && \
  apk add git && \
  rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder /app/trade .
EXPOSE 8090
ENTRYPOINT ["./trade"]