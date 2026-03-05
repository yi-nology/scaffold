# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add --no-cache ca-certificates tzdata git gcc musl-dev

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build both binaries
RUN CGO_ENABLED=1 GOOS=linux go build -a -o scaffold ./cmd/scaffold
RUN CGO_ENABLED=1 GOOS=linux go build -a -o scaffold-server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata sqlite git

ENV TZ=Asia/Shanghai

RUN adduser -D -s /bin/sh app

WORKDIR /app

RUN mkdir -p /app/data /app/cache /app/configs && chown -R app:app /app

COPY --from=builder /app/scaffold .
COPY --from=builder /app/scaffold-server .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/templates ./templates

RUN chown -R app:app ./

USER app

EXPOSE 9090

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:9090/health || exit 1

ENV CONFIG_PATH=/app/configs/config.yaml

ENTRYPOINT ["./scaffold-server"]
