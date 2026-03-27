# System Design — Interview Prep & Architectural Thinking

> A structured guide for senior/staff-level system design interviews.
> Focused on **Order Management domain**, **GCP stack**, **Terraform IaC**,
> and **DevOps-first mindset**. No boilerplate — only mental models that
> change how you think about distributed systems.

---

## Philosophy

System design interviews don't test whether you've memorized architectures.
They test **how you think** under ambiguity. This guide teaches:

1. **A repeatable framework** — the same structured approach for any problem
2. **Distributed systems intuition** — not formulas, but mental models
3. **Domain expertise** — order management patterns used at scale
4. **Cloud-native thinking** — GCP services as building blocks, not magic
5. **Infrastructure as Code** — Terraform as a first-class engineering practice

---

## Learning Path

| # | Chapter | What You'll Learn |
|---|---------|-------------------|
| 01 | [The Framework](./01_the_framework.md) | Structured approach to any system design interview: clarify → estimate → design → deep-dive → tradeoffs |
| 02 | [Distributed Systems Fundamentals](./02_distributed_systems.md) | CAP theorem (what it really means), consistency models, consensus, partitioning, distributed transactions |
| 03 | [Database Design](./03_database_design.md) | PostgreSQL internals, schema design, indexing strategies, partitioning, read replicas, connection pooling |
| 04 | [Event-Driven Architecture](./04_event_driven_architecture.md) | Pub/Sub patterns, event sourcing, CQRS, exactly-once semantics, idempotency, dead letter queues |
| 05 | [Microservices Patterns](./05_microservices_patterns.md) | Service decomposition, API design (REST/gRPC), resilience (circuit breaker, retry, bulkhead), saga pattern |
| 06 | [GCP Architecture](./06_gcp_architecture.md) | Cloud Run, GKE, Cloud SQL, Pub/Sub, Spanner, BigQuery, Memorystore — when to use what and why |
| 07 | [Terraform & IaC](./07_terraform_iac.md) | Terraform patterns, modules, state management, CI/CD pipelines for infrastructure, GCP provider |
| 08 | [Observability & Reliability](./08_observability_reliability.md) | SLOs/SLIs/SLAs, monitoring strategies, alerting, incident response, chaos engineering, postmortems |
| 09 | [Order Management Domain](./09_order_management_domain.md) | Order lifecycle, inventory management, fulfillment, payment flows, order promise, saga orchestration |
| 10 | [Practice Problems](./10_practice_problems.md) | End-to-end system design walkthroughs using the framework from Chapter 01 |

---

## Reading Order

```
Start here (MANDATORY):
  01 The Framework         ← learn this first, use it for everything

Foundations (any order):
  02 Distributed Systems   ← the "why" behind every design decision
  03 Database Design       ← you'll need this for every problem
  04 Event-Driven          ← modern systems are event-driven

Architecture (any order):
  05 Microservices         ← service decomposition and resilience
  06 GCP Architecture      ← your tech stack
  07 Terraform & IaC       ← DevOps-first mindset

Production (any order):
  08 Observability         ← how you prove your system works

Domain + Practice (do last):
  09 Order Management      ← ties everything together for your domain
  10 Practice Problems     ← rehearse with the framework
```

---

## How to Use This Guide

**For interview prep**: Start with Chapter 01 (The Framework). Practice it out
loud — system design interviews are verbal. Then read the domain chapter (09)
and do the practice problems (10).

**For deep understanding**: Read chapters 02-08 in any order. Each chapter is
self-contained but cross-references others. The depth here goes beyond
interview prep — it's how production systems actually work.

**The golden rule**: In a system design interview, there is no "right answer."
There are only **well-reasoned tradeoffs**. Every chapter teaches you to
articulate tradeoffs, not memorize solutions.
