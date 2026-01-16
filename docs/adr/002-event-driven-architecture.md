# ADR 002: Event-Driven Architecture

## Status
Accepted

## Context
Synchronous communication between services leads to tight coupling
and cascading failures.

## Decision
Adopt an event-driven architecture using NATS JetStream for
inter-service communication.

## Rationale
- Services can evolve independently
- Improved resilience
- Better scalability under high load

## Consequences
### Positive
- Loose coupling
- Asynchronous processing
- Better fault tolerance

### Negative
- Eventual consistency
- More complex debugging
