# 01 — The Framework

> A repeatable, structured approach to any system design interview.
> This is the most important chapter. Master this before anything else.

---

## Table of Contents

1. [Why You Need a Framework](#1-why-you-need-a-framework)
2. [The Five Phases](#2-the-five-phases)
3. [Phase 1: Clarify — Requirements & Scope](#3-phase-1-clarify--requirements--scope)
4. [Phase 2: Estimate — Back-of-Envelope Math](#4-phase-2-estimate--back-of-envelope-math)
5. [Phase 3: High-Level Design — The Big Picture](#5-phase-3-high-level-design--the-big-picture)
6. [Phase 4: Deep Dive — The Interesting Parts](#6-phase-4-deep-dive--the-interesting-parts)
7. [Phase 5: Tradeoffs & Production Readiness](#7-phase-5-tradeoffs--production-readiness)
8. [Communication Anti-Patterns](#8-communication-anti-patterns)
9. [The Evaluation Rubric — What Interviewers Score](#9-the-evaluation-rubric--what-interviewers-score)
10. [Quick Reference Card](#10-quick-reference-card)

---

## 1. Why You Need a Framework

System design interviews are 45-60 minutes of **ambiguity**. The interviewer gives
you a vague problem ("Design an order management system") and watches how you think.

Without a framework:
```
  "Design an order management system"
      ↓
  Panic → jump to database schema → forget about scale → miss edge cases
      ↓
  Interviewer thinks: "This person doesn't think structurally"
```

With a framework:
```
  "Design an order management system"
      ↓
  Clarify → Estimate → Design → Deep-dive → Tradeoffs
      ↓
  Interviewer thinks: "This person thinks like an architect"
```

The framework isn't rigid — it's a **skeleton** you adapt. The interviewer may
redirect you ("let's focus on the payment flow"), and you adjust. But you always
know where you are and what's next.

---

## 2. The Five Phases

```
┌─────────────────────────────────────────────────────────────────────┐
│                    SYSTEM DESIGN INTERVIEW                         │
│                                                                     │
│   Phase 1: CLARIFY          (~5 min)    "What exactly are we       │
│   ─────────────────                      building?"                │
│   Functional requirements                                          │
│   Non-functional requirements                                      │
│   Scope boundaries                                                  │
│                                                                     │
│   Phase 2: ESTIMATE         (~3 min)    "How big is this?"         │
│   ─────────────────                                                 │
│   Users, QPS, storage                                               │
│   Read/write ratio                                                  │
│   Data size calculations                                            │
│                                                                     │
│   Phase 3: HIGH-LEVEL       (~10 min)   "The 30,000 ft view"      │
│   ───────────────────                                               │
│   Core components                                                   │
│   Data flow                                                         │
│   API design                                                        │
│                                                                     │
│   Phase 4: DEEP DIVE        (~20 min)   "The hard parts"          │
│   ──────────────────                                                │
│   Database schema                                                   │
│   Scaling strategy                                                  │
│   Edge cases & failure modes                                        │
│                                                                     │
│   Phase 5: TRADEOFFS        (~5 min)    "What would you change?"  │
│   ──────────────────                                                │
│   What you'd do differently with more time                          │
│   Monitoring & observability                                        │
│   Cost optimization                                                 │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

The time splits are approximate. The interviewer may spend 30 minutes in Phase 4
if the deep dive is interesting. That's fine — the framework lets you know when
to naturally transition.

---

## 3. Phase 1: Clarify — Requirements & Scope

**Goal**: Turn a vague prompt into a concrete problem. Show that you don't
assume — you ask.

### Functional Requirements (What the system DOES)

Ask questions that define the **user-facing behavior**:

```
  "Design an order management system"

  Questions to ask:
  ├── Who are the users? (customers, merchants, internal ops?)
  ├── What's the core flow? (browse → cart → checkout → payment → fulfill?)
  ├── Do we need real-time order tracking?
  ├── Do we handle returns/refunds?
  ├── Multi-currency? Multi-region?
  ├── What's the order promise? (delivery ETA commitment)
  └── Do we need inventory management or is it a separate system?
```

**Pro tip**: Frame requirements as **user stories**, not technical features:
- ✅ "A customer can place an order and receive a delivery ETA"
- ❌ "We need a message queue for order events"

### Non-Functional Requirements (How the system BEHAVES)

These are what make the problem **hard**. Always ask about:

```
  ┌─────────────────────────────────────────────────────────────────┐
  │  Non-Functional Requirement     Question to Ask                 │
  ├─────────────────────────────────────────────────────────────────┤
  │  Scale                          How many orders/day? QPS?       │
  │  Availability                   What's the target? 99.9%?      │
  │  Latency                        P99 for checkout? <500ms?      │
  │  Consistency                    Can we show stale inventory?    │
  │  Durability                     Can we EVER lose an order?     │
  │  Regional                       Single region or global?        │
  │  Compliance                     PCI-DSS for payments? GDPR?    │
  └─────────────────────────────────────────────────────────────────┘
```

### Scope Boundaries (What we're NOT building)

Explicitly state what you're excluding to avoid scope creep:

```
  "For this interview, I'll focus on:
   ✅ Order placement, payment, and fulfillment tracking
   ❌ Product catalog (assume it exists)
   ❌ User authentication (assume handled by API gateway)
   ❌ Recommendation engine
   Does that scope sound right?"
```

**Why this matters**: The interviewer is testing your ability to **manage
ambiguity** — the same skill you need as a tech lead scoping a project.

---

## 4. Phase 2: Estimate — Back-of-Envelope Math

**Goal**: Establish the **scale** of the problem so your design decisions
are justified. You don't need exact numbers — order of magnitude is enough.

### The Estimation Framework

Always work through these in order:

```
  1. Users → QPS → Peak QPS
  2. Storage per record → Total storage
  3. Bandwidth → Network requirements
  4. Read/Write ratio → Caching strategy
```

### Example: Order Management System

```
  Given: E-commerce platform, mid-to-large scale

  Users:
    DAU (Daily Active Users):     1M
    Orders per user per month:    2
    Orders per day:               1M × 2 / 30 ≈ 70K orders/day

  QPS:
    Orders/sec average:           70K / 86400 ≈ 1 QPS (order creation)
    Peak (10x average):           ~10 QPS for order creation
    Order reads (tracking):       ~100 QPS (customers check status)
    Order status updates:         ~50 QPS (from fulfillment pipeline)

  Storage:
    Order record size:            ~2 KB (order + items + address + payment ref)
    Orders per year:              70K × 365 ≈ 25M orders/year
    Storage per year:             25M × 2 KB ≈ 50 GB/year
    5-year retention:             250 GB — fits in a single PostgreSQL instance

  Read/Write ratio:
    Reads (tracking, history):    ~100 QPS
    Writes (create, update):     ~60 QPS
    Ratio:                        ~2:1 (not heavily read-biased)
    → Caching helps but isn't critical at this scale
```

### Key Numbers to Memorize

```
  ┌──────────────────────────────────────────────────────┐
  │  Metric                     Rough Value              │
  ├──────────────────────────────────────────────────────┤
  │  Seconds in a day           ~86,400 ≈ 10^5          │
  │  Seconds in a year          ~31M ≈ 3 × 10^7         │
  │  1 KB                       Short email / JSON doc   │
  │  1 MB                       Small image / 1M chars   │
  │  1 GB                       1K images / 1B chars     │
  │  1 TB                       1M images               │
  │  PostgreSQL single node     ~1-5 TB comfortably      │
  │  Redis single node          ~25 GB RAM typical       │
  │  SSD random read            ~100μs                   │
  │  Network round trip (DC)    ~0.5ms                   │
  │  Network round trip (cross) ~50-150ms                │
  │  Disk sequential read       ~1 GB/s (SSD)           │
  └──────────────────────────────────────────────────────┘
```

### Why Estimates Matter

The estimates drive your **architectural decisions**:

```
  50 GB/year → single PostgreSQL instance is fine → no sharding needed
  100 QPS reads → no cache needed (PostgreSQL handles this easily)
  10 QPS writes → no write bottleneck → no write-behind queue needed

  BUT if the problem was:
  500 GB/day → need partitioning or time-series storage
  100K QPS reads → need caching layer (Redis/Memorystore)
  10K QPS writes → need write buffering, async processing
```

**Say this out loud**: "At 10 QPS writes, a single PostgreSQL instance handles
this comfortably. If we needed to scale 100x, we'd add read replicas and
partition by date." This shows you think in terms of **growth**, not just
the current number.

---

## 5. Phase 3: High-Level Design — The Big Picture

**Goal**: Draw the major components and how data flows between them.
This is the "whiteboard diagram" moment.

### Start With the User Flow

Always trace the **happy path** first:

```
  Customer                API Gateway              Order Service
     │                        │                          │
     │  POST /orders          │                          │
     │───────────────────────→│   validate + auth        │
     │                        │─────────────────────────→│
     │                        │                          │  create order
     │                        │                          │  reserve inventory
     │                        │                          │  initiate payment
     │                        │     order confirmed      │
     │   201 Created          │←─────────────────────────│
     │←───────────────────────│                          │
     │                        │                          │
```

### Then Add Components

Build the architecture incrementally. Don't draw 20 boxes at once:

```
  Step 1: Core services (what MUST exist)
  ┌────────────┐     ┌────────────┐     ┌────────────┐
  │  API GW    │────→│  Order Svc │────→│ PostgreSQL │
  └────────────┘     └────────────┘     └────────────┘

  Step 2: Add supporting services
  ┌────────────┐     ┌────────────┐     ┌────────────┐
  │  API GW    │────→│  Order Svc │────→│ PostgreSQL │
  └────────────┘     └─────┬──────┘     └────────────┘
                           │
                    ┌──────┴──────┐
                    ↓             ↓
              ┌──────────┐  ┌──────────┐
              │ Inventory │  │ Payment  │
              │ Service   │  │ Service  │
              └──────────┘  └──────────┘

  Step 3: Add async processing
  ┌────────────┐     ┌────────────┐     ┌────────────┐
  │  API GW    │────→│  Order Svc │────→│ PostgreSQL │
  └────────────┘     └─────┬──────┘     └────────────┘
                           │
                    ┌──────┴──────┐
                    ↓             ↓
              ┌──────────┐  ┌──────────┐
              │ Inventory │  │ Payment  │
              │ Service   │  │ Service  │
              └──────────┘  └──────────┘
                           │
                           ↓
                    ┌──────────────┐     ┌──────────────┐
                    │  Pub/Sub     │────→│ Fulfillment  │
                    │  (events)    │     │ Service      │
                    └──────────────┘     └──────────────┘
```

### API Design (Show It Early)

Define 3-5 core API endpoints — this shows you think about the **interface**:

```
  Core APIs:

  POST   /v1/orders                 Create order
  GET    /v1/orders/{id}            Get order details
  GET    /v1/orders/{id}/status     Get order status (lightweight)
  PATCH  /v1/orders/{id}/cancel     Cancel order
  GET    /v1/orders?customer={id}   List customer's orders (paginated)

  Internal APIs (service-to-service):
  POST   /internal/orders/{id}/fulfill    Fulfillment updates status
  POST   /internal/orders/{id}/refund     Payment triggers refund
```

### Data Model (High Level)

Show the core entities and relationships — not every column, just the shape:

```
  ┌─────────────┐     ┌─────────────┐     ┌─────────────────┐
  │   orders    │     │ order_items  │     │ order_events    │
  ├─────────────┤     ├─────────────┤     ├─────────────────┤
  │ id (PK)     │←───→│ order_id(FK)│     │ order_id (FK)   │
  │ customer_id │     │ product_id  │     │ event_type      │
  │ status      │     │ quantity    │     │ payload (JSONB) │
  │ total       │     │ unit_price  │     │ created_at      │
  │ created_at  │     │ subtotal    │     └─────────────────┘
  │ updated_at  │     └─────────────┘
  │ promise_eta │
  └─────────────┘
```

---

## 6. Phase 4: Deep Dive — The Interesting Parts

**Goal**: This is where you demonstrate **depth**. The interviewer will pick
1-2 areas and go deep. Prepare to go deep in any of these:

### Deep Dive Menu (Be Ready For Any)

```
  1. Database: schema design, indexing, query patterns, consistency
  2. Concurrency: race conditions in inventory, double-ordering
  3. Distributed transactions: saga pattern, compensation
  4. Caching: what to cache, invalidation, consistency
  5. Scaling: horizontal scaling, partitioning, read replicas
  6. Failure modes: what happens when Payment service is down?
  7. Event ordering: how to guarantee order of status updates
  8. Idempotency: retry-safe order creation
```

### Example Deep Dive: Idempotent Order Creation

**The problem**: Customer clicks "Place Order" twice (network glitch,
impatient user). Without idempotency, two orders are created.

```
  Without idempotency:
    Click 1 → POST /orders → Order #1001 created ✅
    Click 2 → POST /orders → Order #1002 created ❌ (duplicate!)

  With idempotency key:
    Click 1 → POST /orders (Idempotency-Key: abc-123) → Order #1001 ✅
    Click 2 → POST /orders (Idempotency-Key: abc-123) → Returns #1001 ✅
```

Implementation:

```sql
  -- Idempotency table
  CREATE TABLE idempotency_keys (
      key         TEXT PRIMARY KEY,
      order_id    UUID NOT NULL,
      created_at  TIMESTAMPTZ DEFAULT NOW(),
      expires_at  TIMESTAMPTZ DEFAULT NOW() + INTERVAL '24 hours'
  );

  -- In the order creation flow:
  BEGIN;
    -- Try to insert idempotency key (fails if duplicate)
    INSERT INTO idempotency_keys (key, order_id)
    VALUES ('abc-123', gen_random_uuid())
    ON CONFLICT (key) DO NOTHING
    RETURNING order_id;

    -- If we got an order_id, it's a new request → create order
    -- If no rows returned, it's a duplicate → look up existing order
  COMMIT;
```

**Say this**: "I'd use an idempotency key in the request header. The client
generates a UUID, we store it with the order. On retry, we detect the duplicate
at the database level using a unique constraint. This is a standard pattern —
Stripe uses the same approach for payment idempotency."

### Example Deep Dive: Order State Machine

Orders have a lifecycle — this is a **state machine**:

```
  ┌──────────┐   payment    ┌──────────┐  picked    ┌──────────┐
  │ CREATED  │─────ok──────→│CONFIRMED │──────────→│ PICKED   │
  └──────────┘              └──────────┘            └──────────┘
       │                         │                       │
       │ payment                 │ cancel                │ ship
       │ failed                  │ request               │
       ↓                         ↓                       ↓
  ┌──────────┐              ┌──────────┐            ┌──────────┐
  │ FAILED   │              │CANCELLED │            │ SHIPPED  │
  └──────────┘              └──────────┘            └──────────┘
                                                         │
                                                         │ deliver
                                                         ↓
                                                    ┌──────────┐
                                                    │DELIVERED │
                                                    └──────────┘
                                                         │
                                                         │ return
                                                         ↓
                                                    ┌──────────┐
                                                    │ RETURNED │
                                                    └──────────┘
```

**Implementation pattern — enforce valid transitions**:

```go
  var validTransitions = map[OrderStatus][]OrderStatus{
      Created:   {Confirmed, Failed},
      Confirmed: {Picked, Cancelled},
      Picked:    {Shipped, Cancelled},
      Shipped:   {Delivered},
      Delivered: {Returned},
      // Failed, Cancelled, Returned are terminal — no outgoing transitions
  }

  func (o *Order) TransitionTo(next OrderStatus) error {
      allowed := validTransitions[o.Status]
      for _, s := range allowed {
          if s == next {
              o.Status = next
              return nil
          }
      }
      return fmt.Errorf("invalid transition: %s → %s", o.Status, next)
  }
```

### Example Deep Dive: Saga Pattern for Order Fulfillment

An order involves multiple services (inventory, payment, fulfillment). You
can't use a distributed transaction (2PC is slow and fragile). Use a **saga**:

```
  Choreography-based saga (event-driven):

  Order Service          Inventory Service       Payment Service
       │                       │                       │
       │  OrderCreated event   │                       │
       │──────────────────────→│                       │
       │                       │  reserve stock        │
       │                       │  InventoryReserved    │
       │                       │──────────────────────→│
       │                       │                       │  charge card
       │                       │                       │  PaymentCharged
       │                       │                       │──────→
       │←──────────────────────────────────────────────│
       │  Update: CONFIRMED    │                       │

  If payment fails → compensation:
       │                       │                       │
       │                       │                       │  PaymentFailed
       │                       │←──────────────────────│
       │                       │  release stock        │
       │←──────────────────────│                       │
       │  Update: FAILED       │                       │
```

**Key principle**: Every saga step has a **compensating action**. If step 3
fails, you undo step 2, then undo step 1. This is eventual consistency in
action.

```
  Step                    Compensating Action
  ─────────────────       ─────────────────────
  Reserve inventory   →   Release inventory
  Charge payment      →   Refund payment
  Create shipment     →   Cancel shipment
```

---

## 7. Phase 5: Tradeoffs & Production Readiness

**Goal**: Show that you think beyond the "happy path." This is what
separates senior from staff-level candidates.

### Always Discuss These

```
  1. What are the single points of failure?
     "If the Order database goes down, everything stops. I'd add a
      read replica for queries and consider multi-zone Cloud SQL."

  2. What would you monitor?
     "Order creation latency (P50, P99), payment failure rate,
      inventory reservation timeout rate, saga completion rate."

  3. What would you do differently at 100x scale?
     "At 7M orders/day, I'd partition the orders table by date,
      add a Redis cache for hot order lookups, and consider CQRS
      for the read path."

  4. What are the security concerns?
     "PCI-DSS compliance for payment data — never store card numbers.
      Use a payment processor (Stripe/Adyen) and store only tokens."

  5. How do you handle deployment?
     "Terraform for infrastructure, GitOps for service deployment.
      Blue-green or canary deployments to avoid downtime."
```

### The Tradeoffs Matrix

For any decision, be ready to articulate **what you chose and what you gave up**:

```
  Decision              Chose              Gave Up
  ─────────────────     ─────────────      ──────────────
  Database              PostgreSQL         Horizontal write scaling
  Consistency           Strong (orders)    Higher latency
  Architecture          Microservices      Operational complexity
  Communication         Async (Pub/Sub)    Immediate consistency
  Saga pattern          Choreography       Central visibility
```

**The magic phrase**: "I chose X because at our scale, Y is more important
than Z. If the requirements changed to [specific scenario], I'd reconsider
and switch to [alternative]."

---

## 8. Communication Anti-Patterns

### ❌ Things That Lose Points

```
  ❌ Jumping to the solution without asking questions
     → Shows you assume, don't clarify

  ❌ Drawing 15 boxes at once without explaining data flow
     → Shows memorization, not understanding

  ❌ "We'll use Kafka" without explaining WHY
     → Name-dropping technologies without justification

  ❌ Ignoring the interviewer's hints
     → They're steering you toward interesting problems

  ❌ Saying "I don't know" and stopping
     → Better: "I haven't worked with X, but I'd approach it by..."

  ❌ Over-engineering for the problem's scale
     → "We need Kubernetes, Kafka, Redis, and Cassandra" for 10 QPS
     → Shows you can't match solutions to problems

  ❌ Never mentioning failure modes
     → Production systems fail. Discussing failures shows maturity
```

### ✅ Things That Win Points

```
  ✅ "Before I start designing, can I ask a few clarifying questions?"
     → Shows structured thinking

  ✅ "At this scale (~10 QPS), a single PostgreSQL instance is sufficient.
      If we needed to scale to 10K QPS, here's what I'd change..."
     → Shows you match solutions to problems AND think about growth

  ✅ "There's a tradeoff here between consistency and availability.
      For orders, I'd choose consistency because..."
     → Shows you understand distributed systems fundamentals

  ✅ "If the Payment service is down, the saga would pause.
      We'd need a dead letter queue and an alert..."
     → Shows you think about failure modes

  ✅ Drawing incrementally and narrating as you go
     → Shows your thought process, not just the answer

  ✅ "This is similar to how Stripe handles idempotent payments..."
     → Shows real-world knowledge (but don't name-drop without substance)
```

---

## 9. The Evaluation Rubric — What Interviewers Score

Most companies score system design on 4-5 axes:

```
  ┌──────────────────────────────────────────────────────────────────┐
  │  Axis                    What They're Looking For                │
  ├──────────────────────────────────────────────────────────────────┤
  │  Problem Exploration     Did you clarify requirements?           │
  │                          Did you identify the hard parts?        │
  │                          Did you scope appropriately?            │
  ├──────────────────────────────────────────────────────────────────┤
  │  Technical Depth         Can you go deep in at least 2 areas?   │
  │                          Do you understand the internals?        │
  │                          Can you discuss specific algorithms?    │
  ├──────────────────────────────────────────────────────────────────┤
  │  Architecture            Is the design sound?                    │
  │                          Are components well-decomposed?         │
  │                          Does data flow make sense?              │
  ├──────────────────────────────────────────────────────────────────┤
  │  Tradeoffs               Can you articulate what you gained      │
  │                          and what you gave up?                   │
  │                          Do you know alternatives?               │
  ├──────────────────────────────────────────────────────────────────┤
  │  Communication           Is your thinking structured?            │
  │                          Can you explain clearly?                │
  │                          Do you collaborate with interviewer?    │
  └──────────────────────────────────────────────────────────────────┘
```

**Senior level (L5)**: Strong in Architecture + Technical Depth.
**Staff level (L6)**: Strong in ALL axes, especially Tradeoffs.
**Principal (L7+)**: All of the above + organizational/business context.

---

## 10. Quick Reference Card

```
  ╔═══════════════════════════════════════════════════════════════════╗
  ║  THE SYSTEM DESIGN FRAMEWORK — POCKET REFERENCE                  ║
  ╠═══════════════════════════════════════════════════════════════════╣
  ║                                                                   ║
  ║  Phase 1: CLARIFY (5 min)                                        ║
  ║  □ Functional requirements (user stories, not tech)              ║
  ║  □ Non-functional (scale, latency, consistency, availability)    ║
  ║  □ Scope boundaries (what's in/out)                              ║
  ║                                                                   ║
  ║  Phase 2: ESTIMATE (3 min)                                       ║
  ║  □ Users → QPS → Peak QPS                                       ║
  ║  □ Storage per record → Total storage → Growth                   ║
  ║  □ Read/Write ratio → Caching strategy                           ║
  ║                                                                   ║
  ║  Phase 3: HIGH-LEVEL DESIGN (10 min)                             ║
  ║  □ Core components (start simple, add incrementally)             ║
  ║  □ Data flow (trace the happy path)                              ║
  ║  □ API design (3-5 core endpoints)                               ║
  ║  □ Data model (core entities and relationships)                  ║
  ║                                                                   ║
  ║  Phase 4: DEEP DIVE (20 min)                                     ║
  ║  □ Database design (schema, indexes, queries)                    ║
  ║  □ Scaling strategy (what breaks first?)                         ║
  ║  □ Failure modes (what if X goes down?)                          ║
  ║  □ Edge cases (idempotency, race conditions, retries)            ║
  ║                                                                   ║
  ║  Phase 5: TRADEOFFS (5 min)                                      ║
  ║  □ What I chose and what I gave up                               ║
  ║  □ What I'd change at 100x scale                                 ║
  ║  □ Monitoring & observability                                    ║
  ║  □ Security & compliance                                         ║
  ║                                                                   ║
  ║  GOLDEN RULES:                                                    ║
  ║  • There is no right answer — only well-reasoned tradeoffs       ║
  ║  • Match the solution to the scale (don't over-engineer)         ║
  ║  • Narrate your thinking out loud                                ║
  ║  • Listen to interviewer hints — they're guiding you             ║
  ║  • "I chose X because at our scale, Y matters more than Z"      ║
  ║                                                                   ║
  ╚═══════════════════════════════════════════════════════════════════╝
```
