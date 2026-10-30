FROM golang:1.23-alpine AS builder
WORKDIR /src/user-service
COPY contracts/ /src/contracts/
COPY user-service/go.mod ./
RUN go mod download
COPY user-service/ .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/service ./cmd/user

FROM alpine:3.21
RUN apk add --no-cache ca-certificates su-exec && addgroup -S app -g 10001 && adduser -S app -u 10001 -G app
WORKDIR /app
COPY --from=builder /out/service /app/service
COPY user-service/scripts/entrypoint.sh /app/entrypoint.sh
RUN chmod 755 /app/entrypoint.sh
EXPOSE 8081 50051
ENTRYPOINT ["/app/entrypoint.sh"]
