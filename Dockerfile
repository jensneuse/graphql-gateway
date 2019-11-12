FROM golang:latest AS builder
WORKDIR /app
ADD go.mod /app
ADD go.sum /app
RUN go mod download
ADD . /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /graphql-gateway main.go

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /graphql-gateway ./
COPY --from=builder /app/logger.config.json ./
COPY --from=builder /app/schema.graphql ./
RUN chmod +x ./graphql-gateway
ENTRYPOINT ["./graphql-gateway"]