FROM golang AS builder
WORKDIR /app
ADD go.mod /app
ADD go.sum /app
RUN go mod download
ADD . /app/
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o /graphql-gateway main.go
COPY /logger.config.json /
COPY  /schema.graphql /
RUN chmod +x /graphql-gateway
ENTRYPOINT ["/graphql-gateway"]