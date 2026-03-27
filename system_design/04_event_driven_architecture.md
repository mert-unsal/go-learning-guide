# 04 — Event-Driven Architecture

> How to decouple services, guarantee delivery, and build systems that
> react to change. Events are the backbone of modern order management.

---

## Table of Contents

1. [Why Events — The Coupling Problem](#1-why-events--the-coupling-problem)
2. [Event Types — Commands vs Events vs Queries](#2-event-types--commands-vs-events-vs-queries)
3. [Messaging Patterns](#3-messaging-patterns)
4. [GCP Pub/Sub — Your Event Backbone](#4-gcp-pubsub--your-event-backbone)
5. [Delivery Guarantees and Idempotency](#5-delivery-guarantees-and-idempotency)
6. [Event Sourcing — Events as the Source of Truth](#6-event-sourcing--events-as-the-source-of-truth)
7. [CQRS — Separate Read and Write Models](#7-cqrs--separate-read-and-write-models)
8. [Dead Letter Queues and Error Handling](#8-dead-letter-queues-and-error-handling)
9. [Event Schema Evolution](#9-event-schema-evolution)
10. [Quick Reference Card](#10-quick-reference-card)

---

## 1. Why Events — The Coupling Problem

### Synchronous (Direct) Communication

```
  Order Service ──HTTP──→ Inventory Service
                 ──HTTP──→ Payment Service
                 ──HTTP──→ Notification Service
                 ──HTTP──→ Analytics Service

  Problems:
  ├── Order Service knows about ALL downstream services
  ├── Adding a new consumer requires changing Order Service
  ├── If Notification is slow, Order creation is slow
  ├── If Analytics is down, Order creation fails (or needs error handling)
  └── Temporal coupling: all services must be up simultaneously
```

### Asynchronous (Event-Driven) Communication

```
  Order Service ──event──→ Pub/Sub ──→ Inventory Service
                                   ──→ Payment Service
                                   ──→ Notification Service
                                   ──→ Analytics Service

  Benefits:
  ├── Order Service publishes ONE event, doesn't know who listens
  ├── Adding a new consumer: just subscribe (no Order Service change)
  ├── If Notification is slow, Orders aren't affected
  ├── If Analytics is down, events queue up and process when it recovers
  └── Temporal decoupling: services don't need to be up simultaneously
```

### The Key Insight

Events invert the dependency direction:

```
  Synchronous: Order Service depends on Inventory, Payment, Notification
  Event-driven: Inventory, Payment, Notification depend on Order events

  This is the same principle as Go's interface design:
  "Accept interfaces, return structs" → consumers define what they need
  Events: producer defines what happened, consumers decide how to react
```

---

## 2. Event Types — Commands vs Events vs Queries

```
  ┌────────────────────────────────────────────────────────────────┐
  │  Type        Intent              Example                      │
  ├────────────────────────────────────────────────────────────────┤
  │  Command     "Do this"           CreateOrder, CancelOrder     │
  │              Directed to one      ReserveInventory             │
  │              service              ChargePayment                │
  │              Imperative verb                                   │
  ├────────────────────────────────────────────────────────────────┤
  │  Event       "This happened"     OrderCreated, OrderShipped   │
  │              Broadcast to many    PaymentCharged               │
  │              Past tense           InventoryReserved            │
  │              Immutable fact                                    │
  ├────────────────────────────────────────────────────────────────┤
  │  Query       "Tell me"           GetOrderStatus               │
  │              Request/response     ListCustomerOrders           │
  │              Synchronous                                       │
  └────────────────────────────────────────────────────────────────┘
```

**For system design interviews**: Use **events** (past tense) for decoupled
communication. Use **commands** when you need a specific service to act.
Use **queries** for reads (synchronous HTTP/gRPC).

### Event Structure

```json
  {
    "event_id": "evt-abc-123",           // unique, for deduplication
    "event_type": "OrderCreated",
    "aggregate_id": "order-456",         // the entity this event belongs to
    "aggregate_type": "Order",
    "version": 1,                        // for ordering within aggregate
    "timestamp": "2024-01-15T10:30:00Z",
    "data": {
      "customer_id": "cust-789",
      "items": [...],
      "total": 99.99
    },
    "metadata": {
      "correlation_id": "req-xyz",       // trace across services
      "caused_by": "evt-previous-123"    // causality chain
    }
  }
```

---

## 3. Messaging Patterns

### Point-to-Point (Queue)

```
  Producer → [Queue] → Consumer

  One message, one consumer. Used for task distribution.
  Example: Order → [Fulfillment Queue] → Fulfillment Worker

  If 5 workers: each message goes to exactly ONE worker (load balancing)
```

### Publish-Subscribe (Topic)

```
  Producer → [Topic] → Subscription A → Consumer A
                     → Subscription B → Consumer B
                     → Subscription C → Consumer C

  One message, many consumers (each subscription gets a copy).
  Example: OrderCreated → Inventory, Payment, Notification, Analytics

  If 5 workers per subscription: messages load-balanced within subscription
```

### Fan-Out / Fan-In

```
  Fan-out: One event triggers multiple parallel processes
  ┌──────────┐     ┌────────────┐
  │ OrderPaid │────→│ Ship Item A│
  │          │────→│ Ship Item B│   Items ship from different warehouses
  │          │────→│ Ship Item C│
  └──────────┘     └────────────┘

  Fan-in: Multiple events converge to one process
  ┌────────────┐
  │ Item A Done│────→┌────────────────┐
  │ Item B Done│────→│ Order Complete │   All items delivered → mark complete
  │ Item C Done│────→│ (aggregator)   │
  └────────────┘     └────────────────┘
```

---

## 4. GCP Pub/Sub — Your Event Backbone

### Architecture

```
  Publisher → Topic → Subscription → Subscriber

  ┌────────────┐     ┌─────────────────────────────┐
  │ Order Svc  │     │         Topic:               │
  │            │────→│    order-events              │
  └────────────┘     │                               │
                     │  ┌─ Sub: inventory-sub ────→ Inventory Svc
                     │  ├─ Sub: payment-sub ──────→ Payment Svc
                     │  ├─ Sub: notify-sub ───────→ Notification Svc
                     │  └─ Sub: analytics-sub ────→ BigQuery (via Dataflow)
                     └─────────────────────────────┘
```

### Key Features

```
  Feature                   Detail
  ──────────                ──────
  Delivery                  At-least-once (default)
  Ordering                  Per-key ordering (optional)
  Retention                 7 days (configurable up to 31)
  Dead letter               Built-in DLQ support
  Filtering                 Attribute-based message filtering
  Push vs Pull              Pull (subscriber controls pace) or Push (HTTP endpoint)
  Throughput                ~10M messages/sec per topic
  Latency                   ~100ms typical
  Exactly-once              Supported within Cloud Dataflow
```

### Message Ordering

By default, Pub/Sub does **not guarantee order**. For order management,
you need events for the same order to arrive in sequence:

```
  Without ordering:
    OrderCreated → PaymentCharged → OrderShipped
    Might arrive as: PaymentCharged → OrderShipped → OrderCreated
    → Subscriber sees "shipped" before "created" → broken state machine

  With ordering key:
    Publish with ordering_key = order_id
    → All events for order-456 arrive in publish order
    → Events for different orders may interleave (that's fine)
```

```go
  // Publishing with ordering key
  result := topic.Publish(ctx, &pubsub.Message{
      Data:        eventJSON,
      OrderingKey: orderID,    // ensures order within this key
      Attributes: map[string]string{
          "event_type": "OrderCreated",
      },
  })
```

---

## 5. Delivery Guarantees and Idempotency

### The At-Least-Once Reality

Pub/Sub guarantees at-least-once delivery. This means your subscriber
**will receive duplicates**:

```
  Scenario:
  1. Pub/Sub delivers message to subscriber
  2. Subscriber processes message
  3. Subscriber sends ACK
  4. Network hiccup — ACK doesn't reach Pub/Sub
  5. Pub/Sub thinks message wasn't processed
  6. Pub/Sub re-delivers message
  7. Subscriber processes AGAIN → duplicate!
```

### Making Subscribers Idempotent

```
  Strategy 1: Deduplication table
  ┌────────────────────────────────────┐
  │  processed_events                  │
  │  ──────────────────                │
  │  event_id TEXT PRIMARY KEY         │
  │  processed_at TIMESTAMPTZ         │
  └────────────────────────────────────┘

  On receive:
    INSERT INTO processed_events (event_id) VALUES ('evt-123')
    ON CONFLICT DO NOTHING
    RETURNING event_id;

    If rows returned → new event → process it
    If no rows → duplicate → skip it

  Strategy 2: Idempotent operations
  Instead of:  UPDATE inventory SET quantity = quantity - 1
  Use:         UPDATE inventory SET quantity = 5 WHERE order_id = 'xyz'
               (setting absolute value, not relative — replay-safe)

  Strategy 3: Version checks
  UPDATE orders SET status = 'shipped', version = 4
  WHERE id = 'order-456' AND version = 3;
  -- If version doesn't match, it's a duplicate or out-of-order
```

---

## 6. Event Sourcing — Events as the Source of Truth

Instead of storing **current state** (traditional CRUD), store **every event
that ever happened**. The current state is derived by replaying events.

```
  Traditional (state-based):
    orders table: {id: 456, status: 'shipped', total: 99.99}
    → You know the current status but NOT the history

  Event sourced:
    order_events table:
    ┌──────┬──────────────────┬────────────────────────────┐
    │ seq  │ event_type       │ data                       │
    ├──────┼──────────────────┼────────────────────────────┤
    │ 1    │ OrderCreated     │ {items: [...], total: 99}  │
    │ 2    │ PaymentCharged   │ {payment_id: "pay-123"}    │
    │ 3    │ OrderConfirmed   │ {}                         │
    │ 4    │ ItemPicked       │ {item_id: "item-1"}        │
    │ 5    │ OrderShipped     │ {tracking: "TR-456"}       │
    └──────┴──────────────────┴────────────────────────────┘

    Current state: replay events 1-5 → status = "shipped"
    → You have complete audit trail + can rebuild state at any point
```

### When to Use Event Sourcing

```
  ✅ Use when:
  • Audit trail is mandatory (financial, compliance)
  • You need to reconstruct past states ("what was the order at 3 PM?")
  • Multiple services need to react to state changes
  • Domain is naturally event-oriented (order lifecycle)

  ❌ Don't use when:
  • Simple CRUD with no history requirements
  • Read-heavy with simple queries (rebuilding state is expensive)
  • Team is unfamiliar (significant learning curve)
```

### Hybrid Approach (Most Practical)

```
  Best of both worlds:
  1. Store current state in orders table (for fast reads)
  2. Store events in order_events table (for history/audit)
  3. Update both in the same transaction

  This is NOT pure event sourcing but gives you:
  ✅ Fast reads (query the orders table directly)
  ✅ Complete audit trail (query order_events)
  ✅ Event publishing (read from order_events or outbox)
  ✅ Simpler than pure event sourcing
```

---

## 7. CQRS — Separate Read and Write Models

**Command Query Responsibility Segregation**: use different models
(even different databases) for writes and reads.

```
  Without CQRS:
  ┌──────────────┐     ┌──────────────┐
  │  API Server  │────→│  PostgreSQL  │  Same schema for reads and writes
  │ (read+write) │←────│  (one DB)    │
  └──────────────┘     └──────────────┘

  With CQRS:
  ┌──────────────┐     ┌──────────────┐
  │ Write API    │────→│  PostgreSQL  │  Normalized, optimized for writes
  └──────────────┘     └──────────────┘
                              │ events
                              ↓
                       ┌──────────────┐
                       │   Pub/Sub    │
                       └──────┬───────┘
                              ↓
  ┌──────────────┐     ┌──────────────┐
  │ Read API     │←────│  Redis or    │  Denormalized, optimized for reads
  └──────────────┘     │  Read DB     │
                       └──────────────┘
```

### When CQRS Makes Sense

```
  ✅ Read and write patterns are very different
     (writes: complex validation, business rules
      reads: simple lookups, different query shapes)

  ✅ Read QPS >> Write QPS (need different scaling strategy)

  ✅ Read model needs different shape
     (write: normalized tables
      read: pre-joined, denormalized documents)

  ❌ Simple CRUD applications
  ❌ Team unfamiliar with eventual consistency
  ❌ Read and write patterns are similar
```

### For Order Management (Practical CQRS)

```
  Write side (PostgreSQL):
    Normalized schema, enforces business rules, ACID transactions
    Used by: order creation, status updates, cancellation

  Read side (Redis or read replica):
    Denormalized order view, pre-computed aggregations
    Used by: order tracking page, customer order history, dashboard

  Sync mechanism:
    Order events (via Pub/Sub) → read model updater → Redis/read DB
    Eventual consistency: read model lags by ~100ms
```

---

## 8. Dead Letter Queues and Error Handling

### The Problem

What happens when a subscriber can't process a message?

```
  Pub/Sub → Subscriber: "OrderCreated for order-456"
  Subscriber: tries to process → ERROR (maybe DB is down)
  Subscriber: NACKs the message
  Pub/Sub: re-delivers
  Subscriber: ERROR again
  Pub/Sub: re-delivers
  ... infinite loop of failures ("poison pill" message)
```

### Dead Letter Queue (DLQ)

```
  Normal flow:
  Topic → Subscription → Subscriber
                              ↓ (success)
                           ACK → done

  After N failures:
  Topic → Subscription → Subscriber
                              ↓ (fail × 5)
                       Dead Letter Topic → DLQ Subscription
                                                ↓
                                          Alert + manual review
                                          or automated retry later
```

### GCP Pub/Sub DLQ Configuration

```
  Subscription config:
    max_delivery_attempts: 5
    dead_letter_topic: "projects/my-proj/topics/order-events-dlq"

  After 5 failed deliveries:
  → Message moves to DLQ topic
  → Alert fires (Cloud Monitoring)
  → On-call engineer investigates
  → Fix the issue, republish from DLQ
```

### Error Classification

```
  Error Type          Action              Example
  ──────────          ──────              ───────
  Transient           Retry (with backoff) DB connection timeout
  Permanent           Send to DLQ          Invalid message format
  Poison pill         Send to DLQ          Business logic bug
  Rate limited        Backoff + retry      Third-party API limit
  Dependency down     Retry later          Payment service offline
```

---

## 9. Event Schema Evolution

### The Problem

Your event schema will change over time. Subscribers must handle old AND new formats.

```
  V1: OrderCreated { customer_id, total }
  V2: OrderCreated { customer_id, total, currency }     ← new field
  V3: OrderCreated { customer_id, total_amount, currency } ← renamed field
```

### Compatibility Rules

```
  Backward compatible (SAFE):
  ├── Add optional field with default     ✅
  ├── Add new event type                  ✅
  └── Deprecate field (keep it, add new)  ✅

  Breaking change (DANGEROUS):
  ├── Remove field                         ❌
  ├── Rename field                         ❌
  ├── Change field type                    ❌
  └── Change field semantics               ❌
```

### Versioning Strategies

```
  Strategy 1: Schema version in event
  {
    "schema_version": 2,
    "event_type": "OrderCreated",
    "data": { ... }
  }
  → Subscriber checks version, handles each differently

  Strategy 2: Event type versioning
  "OrderCreated.v1", "OrderCreated.v2"
  → Different subscriptions or routing per version

  Strategy 3: Always additive (recommended)
  Never remove or rename fields. Only add new ones.
  Old subscribers ignore new fields (forward compatible).
  New subscribers handle missing fields with defaults.
```

---

## 10. Quick Reference Card

```
  ╔═══════════════════════════════════════════════════════════════════╗
  ║  EVENT-DRIVEN ARCHITECTURE — POCKET REFERENCE                    ║
  ╠═══════════════════════════════════════════════════════════════════╣
  ║                                                                   ║
  ║  Events: past tense ("OrderCreated"), immutable facts            ║
  ║  Commands: imperative ("CreateOrder"), directed to one service   ║
  ║                                                                   ║
  ║  Pub/Sub: topic + subscription model                             ║
  ║  → At-least-once delivery (make subscribers idempotent!)        ║
  ║  → Use ordering_key for same-entity event ordering              ║
  ║  → DLQ after N failures, alert + investigate                     ║
  ║                                                                   ║
  ║  Outbox pattern: DB write + outbox insert in one transaction     ║
  ║  → Poller/CDC reads outbox → publishes to Pub/Sub               ║
  ║  → Guarantees: if order exists, event will be published         ║
  ║                                                                   ║
  ║  Event sourcing: store events, derive state                      ║
  ║  → Audit trail, temporal queries, multiple projections           ║
  ║  → Or hybrid: state table + events table (practical)             ║
  ║                                                                   ║
  ║  CQRS: different models for reads and writes                     ║
  ║  → Use when read/write patterns diverge significantly            ║
  ║                                                                   ║
  ║  Schema evolution: always additive, never remove fields          ║
  ║                                                                   ║
  ╚═══════════════════════════════════════════════════════════════════╝
```
