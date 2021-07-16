# Roadmap

Exo provides an exoskeleton that supports and protects your services.

(NOTE: THIS DOCUMENT IS VERY ASPIRATIONAL)

In dev, Exo acts as a supervisor. In prod, you can use your favorite
orchestration system. IO is instrumented with proxies, giving your service new
functionality without any code changes in your process. All this functionality
is exposed through a management console that lives as a sidecar to your
service.

## Features

- Admin Console
- Logging
- HTTP Request Logging
- Error Tracking
- Online Help
- Metrics
- Lots more...

## Getting Started

TODO: Write migration guides for...

- Procfiles
- Docker Compose
- Minikube

## Concepts

### Components.

Like resources in an Infrastructure-as-Code system.

Each component has a typed and has a spec, plus some metadata.
