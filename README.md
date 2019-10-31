# GraphQL Gateway

GraphQL Gateway turns a annotated schema into an executable GraphQL Server.
Simply describe your schema and use directives to describe how the execution engine should resolve a field.
No code needs to be written, no code generation, just configuration and the execution engine does its thing.

Currently DataSources can be Static, GraphQL or JSON HTTP APIs.
You can nest GraphQL data sources into REST data sources into GraphQL data sources into...
Have a look at the schema.graphql file to fully understand the concept.

This is experimental, don't use in production! But please, try it out!

# How to?

1. Install go 1.13. (golang.org)
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

# Contribute

Please try this tool for your use case and submit issues in case it didn't work for you.
It's far from being production ready so I'm relying on your input to find missing features & bugs.

# Goals

1. Based on you input I want to evolve the engine to support as many use cases as possible.
2. When it feels "production ready" I'll focus on increasing performance because I took a few shortcuts to release this PoC faster (e.g. no parallel fetching, creation of unnecessary garbage collection)