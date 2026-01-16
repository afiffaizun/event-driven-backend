# System Architecture

## Overview
This platform uses event-driven microservices to ensure scalability and loose coupling.

## Services
- Auth Service
- Order Service
- Notification Service

## Communication
- REST for synchronous requests
- NATS JetStream for asynchronous events

## Data Management
Each service owns its database to avoid tight coupling.

## Reliability
- At-least-once delivery
- Retry and DLQ
- Idempotent consumers

## Observability
- Metrics via Prometheus
- Dashboards in Grafana
