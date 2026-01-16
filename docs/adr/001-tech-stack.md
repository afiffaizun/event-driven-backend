# ADR 001: Technology Stack Selection

## Status
Accepted

## Context
We need to build a scalable, maintainable backend platform that supports
event-driven communication, high performance, and cloud-native deployment.

## Decision
- Language: Go (1.25)
- API: REST
- Messaging: NATS JetStream
- Database: PostgreSQL
- Cache: Redis
- Container & Orchestration: Docker & Kubernetes
- CI/CD: GitLab CI
- GitOps: ArgoCD

## Rationale
- Go provides high performance, simplicity, and strong concurrency support.
- NATS JetStream is lightweight, easy to operate, and fits event-driven workloads.
- PostgreSQL offers reliability and strong transactional guarantees.
- Kubernetes enables scalable and resilient deployments.

## Consequences
### Positive
- High scalability and loose coupling
- Easier horizontal scaling
- Clear service ownership

### Negative
- Operational complexity increases
- Requires strong observability and monitoring
