FROM golang AS builder
WORKDIR /app
ADD go.mod /app
ADD go.sum /app
RUN go mod download
ADD . /app/
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o /graphql-gateway main.go
RUN chmod +x /graphql-gateway
COPY /logger.config.json /
COPY /schema.graphql /
COPY /person.wasm /
EXPOSE 9111
WORKDIR /
ENTRYPOINT ["/graphql-gateway"]