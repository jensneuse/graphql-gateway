# nats example

## start nats server

````shell script
docker run -p -d 4222:4222 nats
````

## start producer

```shell script
go run main.go nats startProducer
```

## start consumer

Just in case you'd like to watch the raw nats messages.

```shell script
go run main.go nats startConsumer
```