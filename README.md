This codebase is **deprecated**!
Please use https://github.com/jensneuse/graphql-go-tools instead.

# GraphQL Gateway

GraphQL Gateway turns an annotated schema into an executable GraphQL Server.
Simply describe your schema and use directives to describe how the execution engine should resolve a field.
No code needs to be written, no code generation, just configuration and the execution engine does its thing.

Currently DataSources can be Static, GraphQL or JSON HTTP APIs.
You can nest GraphQL data sources into REST data sources into GraphQL data sources into...
Have a look at the schema.graphql file to fully understand the concept.

This is experimental, don't use in production! But please, try it out!

# Run

Run the gateway using docker.

```shell script
docker run -p 0.0.0.0:9111:9111 jensneuse/graphql-gateway serve
```

# Docs

https://jens-neuse.gitbook.io/graphql-gateway/

# Develop

1. Install go 1.15 (golang.org)
2. git clone https://github.com/jensneuse/graphql-gateway.git
3. cd graphql-gateway
4. run:

```shell script
go run main.go serve
```

Use GraphQL Playground to explore your API.

Update the schema and make use of the directives.
Take inspiration from the schema.graphql file on how to use the directives.
If something doesn't work for you please supply a schema file so I can reconstruct the problem.

# Motivation

Writing resolvers on top of existing data sources is repetitive.
The logic can be encapsulated into directives.
If a directive doesn't give enough flexibility you're free to open an issue and/or PR.

Before wrapping existing data sources think about the costs of writing and maintaining the resolver code.
Implementing, extending and maintaining a bunch of directives scales a lot better than writing individual resolvers.
There's a lot of tooling for writing efficient REST services. Also it's very easy to apply caching (e.g. using Varnish) to a REST API.

With the help of this gateway you'll get the best of both worlds, easy to reason about backend services and the benefits of GraphQL at the front end layer:
Write simple CRUD based REST APIs, apply server side caching and turn them into a GraphQL API with very little effort.

# Architecture

Query -> Lexing -> Parsing -> Normalization -> Validation -> Create Query Execution Plan -> Execute Query Plan -> Return Result

Caching of execution plans will be possible to increase performance.
This will enable skipping validation and query planning for recurring operations.
After normalization the input can be hashed and the hash will be used to fill a query plan map / retrieve a cached query plan.

# *any DataSource

*any DataSource can be supported implementing the necessary interfaces.
Currently supported:
- GraphQL (multiple GraphQL services can be combined)
- static (static embedded data)
- HTTP JSON
- HTTP JSON Streaming (uses polling to create a stream)
- MQTT
- Nats
- Webassembly (resolve a Request using WASI compliant modules)

# Contribute

Please try this tool for your use case and submit issues in case it didn't work for you.
It's far from being production ready so I'm relying on your input to find missing features & bugs.

# Goals

1. Based on you input I want to evolve the engine to support as many use cases as possible.
2. When it feels "production ready" I'll focus on increasing performance because I took a few shortcuts to release this PoC faster (e.g. no parallel fetching, creation of unnecessary garbage collection)

# Comparison

### nautilus gateway
https://github.com/nautilus/gateway

Nautilus is based around the assumption that all your services are GraphQL Services. It uses the Node interface ([Relay Global Object Identification](https://facebook.github.io/relay/graphql/objectidentification.htm)) to automatically federate between multiple services. This approach is quite similar to the approach Apollo Federation took. So all services have to comply to the Relay Global Object spec and you're ready to go. Nautilus gateway will analyze your services via introspection at startup time and generate the final gateway schema.

In contrast this library goes a complete different approach. The basic assumption is that your public gateway schema should not be an artifact. The gateway schema is the contract between gateway and clients. The gateway doesn't make any assumptions on your services other than complying to some protocol and some spec. It's a lot more manual work at the beginning but this gives a lot of advantages. Because the gateway is the single source of truth regarding the schema you cannot easily break the contract. With federation an upstream service can directly break the contract. Additionally you're not limited to GraphQL upstreams. While GraphQL upstreams are supported it's only a matter of implementing another DataSource interface to support more upstream protocols, e.g. SOAP, MQTT, KAFKA, Redis etc.. On top of that, because the gateway schema is the single source of truth you'll get two additional benefits. First, you can swap DataSources for a given schema without changing the contract. E.g. you could replace a REST User-Service with a GraphQL User-Service without changing the contract between client and gateway. This could help to easily transition from legacy to more modern architectures. Second, because the gateway owns the schema you can apply features like rate limiting, authorization etc. at the gateway level.

To conclude, if you have only GraphQL services complying to the Relay Global Object spec and there are no special requirements like rate limiting or authZ you're best off using nautilus as it's a lot easier to get started.

If you have a heterogenous mesh of services and more specific requirements that would require customizing the schema you'd want to put in the extra effort to configure it all manually in exchange for the benefits described above.

As a sidenote, nautilus relies on the amazing gqlgen library as well as gqlparser from vektah, both amazing tools which are quite mature I think. In comparison I've implemented the GraphQL spec myself ([graphql-go-tools](https://github.com/jensneuse/graphql-go-tools)) for this specific use case. I've paid extra attention to implement features like parsing/lexing/validation/ast-walking in a zero garbage collection fashion which ultimately leads to better performance and more consistent latencies than every available implementation. I won't add benchmarks here because I'm biased but you're invited to prove me wrong.
