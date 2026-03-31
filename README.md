# Custom Rust Ingress Controller

## Project Overview
### Motivation
I came up with this project to build a high-throughput hand built project in Rust because I love the language and want to have an exploratory development project to get extra familiar with K8s, Go, and build something cool in Rust

**The Goal:** See how hard I can push a custom Rust ingress controller before it breaks<br>
**The Rule:** No Pre-Optimization! Use metrics to determine bottlenecks and iterate

### Tech Stack
#### Infrastructure and Orchestration
- Kubernetes through KinD (Kubernetes in Docker) with simple horizontal scaling
- Containerized with Docker for multi-stage builds, dev and load test
- Skaffold for auto deployment and live reloading 

#### Application Stack
- Rust for custom ingress controller
- Go backend being system under test
- Postgres for persistent database storage

#### Testing
- K6 for load testing
- Request headers parsed by K6 into metrics

### Why Rust/Go
- **Rust**: My all-time favorite language because of its unique ownership and borrowing concepts and its ability to optimize for low resource systems and high throughput systems alike
- **Go**: Concurrency and shared resources are ideal for a simple, performant backend service that can operate at exceptionally high RPS

As the intended deployment for this is local, it adds certain constraints to how hard I can load test this system. Since Rust and Go compile to small binaries, it allows for more testing at a higher RPS. This enables me to see the auto-scaling working as intended, testing both the Ingress Controller and the Go service with realistic resource constraints on a single machine.
