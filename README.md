# easyhttp

## Features

- Plugin driven architecture.
- Simple, expressive, fluent API.
- Idiomatic built on top of `net/http` package.
- Context-aware hierarchical middleware layer supporting all the HTTP life cycle.
- Built-in multiplexer for easy composition capabilities.
- Easy to extend via plugins/middleware.
- Ability to easily intercept and modify HTTP traffic on-the-fly.
- Convenient helpers and abstractions over Go's HTTP primitives.
- URL template path params.
- Built-in JSON, XML and multipart bodies serialization and parsing.
- Easy to test via HTTP mocking (e.g: [gentleman-mock](https://github.com/h2non/gentleman-mock)).
- Supports data passing across plugins/middleware via its built-in context.
- Fits good while building domain-specific HTTP API clients.
- Easy to hack.
- Dependency free.
