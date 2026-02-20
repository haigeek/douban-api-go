FROM golang:1.21-alpine AS builder
WORKDIR /src
ARG TARGETOS
ARG TARGETARCH
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags "-s -w" -o /out/douban-api-go ./cmd/server

FROM alpine:3.20
RUN apk --no-cache add ca-certificates tini tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
WORKDIR /data
COPY --from=builder /out/douban-api-go /usr/bin/douban-api-go
EXPOSE 80
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["/usr/bin/douban-api-go", "--port", "80"]
