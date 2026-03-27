# 02 — Distributed Systems Fundamentals

> The mental models behind every system design decision.
> Not academic theory — practical intuition for building systems that work.

---

## Table of Contents

1. [CAP Theorem — What It Actually Means](#1-cap-theorem--what-it-actually-means)
2. [Consistency Models — The Spectrum](#2-consistency-models--the-spectrum)
3. [Consensus — How Nodes Agree](#3-consensus--how-nodes-agree)
4. [Partitioning — Splitting Data](#4-partitioning--splitting-data)
5. [Replication — Copying Data](#5-replication--copying-data)
6. [Distributed Transactions](#6-distributed-transactions)
7. [Time and Ordering in Distributed Systems](#7-time-and-ordering-in-distributed-systems)
8. [Failure Modes — Everything Fails](#8-failure-modes--everything-fails)
9. [Quick Reference Card](#9-quick-reference-card)

---

## 1. CAP Theorem — What It Actually Means

CAP states that during a **network partition**, you must choose between
**Consistency** and **Availability**. You can't have both.

```
  ┌─────────────────────────────────────────────────────────────────┐
  │                     CAP Theorem                                 │
  │                                                                 │
  │              Consistency (C)                                    │
  │              /           \                                      │
  │             /             \                                     │
  │            /               \                                    │
  │    Availability (A) ─── Partition Tolerance (P)                │
  │                                                                 │
  │  During a partition, choose ONE:                                │
  │    CP: Consistent but some requests fail (refuse stale data)   │
  │    AP: Available but may return stale data                      │
  │                                                                 │
  │  CA doesn't exist in distributed systems                       │
  │  (partitions WILL happen — network is unreliable)              │
  └─────────────────────────────────────────────────────────────────┘
```

### The Misconception

Most people think CAP means "pick 2 of 3." That's wrong. In a distributed
system, partitions are **inevitable**. The real question is:

**When a partition happens, do you return an error (CP) or stale data (AP)?**

### Real-World CAP Decisions

```
  System          Choice    Why
  ──────────      ──────    ───────────────────────────────────
  Bank transfer   CP        Wrong balance is unacceptable
  Order creation  CP        Double-charging is unacceptable
  Product catalog AP        Showing stale price for 5 seconds is fine
  Shopping cart   AP        Cart should always be available
  Inventory count AP*       Show "approximately 5 left" is fine
                            (*but reserve with strong consistency)
```

### For Order Management: Per-Operation CAP

You don't choose CP or AP for the entire system — you choose **per operation**:

```
  Operation                 Choice    Reason
  ──────────────            ──────    ──────────────────────
  Place order               CP        Cannot lose orders
  Check order status        AP        Stale by 1s is acceptable
  Update inventory count    CP        Cannot oversell
  Browse product catalog    AP        Stale by 1 min is fine
  Calculate delivery ETA    AP        Approximate is expected
```

---

## 2. Consistency Models — The Spectrum

Consistency isn't binary (consistent vs inconsistent). It's a **spectrum**:

```
  Strongest ←───────────────────────────────────────→ Weakest

  Linearizable   Sequential    Causal     Eventual
  ────────────   ──────────    ──────     ────────
  "Real-time     "Ordered      "Related   "Eventually
   ordering"      like a        events     all nodes
                   serial        in         agree"
                   program"      order"

  Latency:       lowest ←────────────────→ highest
  Availability:  lowest ←────────────────→ highest
  Complexity:    highest ←───────────────→ lowest
```

### What Each Model Means (Practically)

**Linearizable (Strongest)**:
Every read sees the most recent write, globally. As if there's one copy.
- Use for: payments, inventory reservation, leader election
- Cost: highest latency (must coordinate across nodes)
- Example: Google Spanner (uses TrueTime for global ordering)

**Sequential Consistency**:
Operations appear in some total order, consistent with each client's order.
- Use for: order status updates (within one order, events are ordered)
- Cost: high, but less than linearizable

**Causal Consistency**:
If event A caused event B, everyone sees A before B. Unrelated events may
appear in any order.
- Use for: social media feeds, comment threads
- Cost: moderate

**Eventual Consistency**:
All replicas will eventually have the same data, but reads may be stale.
- Use for: analytics, caches, read replicas, DNS
- Cost: lowest latency, highest availability

### For Your Interview

When the interviewer asks about consistency, say:

> "For the order creation path, I'd use strong consistency —
> we can't risk double-charging. For the order tracking read path,
> I'd accept eventual consistency with a read replica that's at most
> a few seconds behind. The delivery ETA is inherently approximate,
> so eventual consistency is fine there."

---

## 3. Consensus — How Nodes Agree

When multiple nodes need to agree on a value (leader election, configuration),
they use a **consensus protocol**.

### Raft (Most Practical to Understand)

```
  ┌─────────┐     ┌─────────┐     ┌─────────┐
  │ Leader  │────→│Follower │     │Follower │
  │ (Node A)│     │(Node B) │     │(Node C) │
  └─────────┘     └─────────┘     └─────────┘
       │               ↑               ↑
       └───────────────┴───────────────┘
            Replicates log entries

  Write: Client → Leader → replicate to majority → commit → respond
  Read:  Client → Leader → respond (guaranteed latest)

  If Leader dies:
    Followers detect heartbeat timeout → election → new Leader
    Requires majority (2 of 3, 3 of 5) → why odd numbers
```

### You Won't Implement Consensus, But You'll Use It

```
  Service                     Uses Consensus For
  ───────────────             ──────────────────────
  PostgreSQL (streaming rep)  Leader election, WAL replication
  etcd / Consul               Distributed config, service discovery
  Kubernetes                   Control plane consistency (etcd)
  Google Spanner               Global strong consistency (Paxos)
  Apache Kafka                 Partition leader election (KRaft)
```

**For interviews**: You don't need to explain Raft internals. Know that consensus
requires a **majority quorum** (why 3 or 5 nodes), and that it adds latency
(must wait for majority acknowledgment).

---

## 4. Partitioning — Splitting Data

When one database can't hold all the data or handle all the traffic, you **partition**
(also called "sharding") the data across multiple nodes.

### Partitioning Strategies

```
  Strategy            How                  Good For               Bad For
  ──────────          ───                  ────────               ───────
  Range               order_id 1-1M,       Range queries          Hot spots
                      1M-2M, etc.          (orders by date)       (recent orders)

  Hash                hash(order_id)       Even distribution      Range queries
                      mod N                                       need scatter-gather

  Geographic          region = "US",       Data locality          Cross-region
                      region = "EU"        (GDPR compliance)      queries

  Time-based          orders_2024_01,      Archival, TTL          Queries across
                      orders_2024_02       Cold data pruning      time ranges
```

### For Order Management

```
  Table            Partitioning Strategy    Why
  ─────            ─────────────────────    ───
  orders           Time-based (monthly)     Old orders rarely accessed, easy archival
  order_items      Co-located with orders   Always queried together
  order_events     Time-based (daily)       Event log grows fast, old events archived
  customers        Hash (customer_id)       Even distribution of queries
  inventory        None (small dataset)     Fits in one node comfortably
```

**Key insight**: Partition by how you **query**, not by how you **insert**.
If you always query orders by customer, partition by customer_id.
If you always query orders by date, partition by date.

---

## 5. Replication — Copying Data

Replication puts copies of data on multiple nodes for **availability** and
**read scaling**.

### Replication Topologies

```
  Single Leader (most common):
  ┌─────────┐     ┌──────────┐
  │ Leader  │────→│ Replica  │   Writes → Leader only
  │ (R/W)   │     │ (R only) │   Reads  → Leader or Replica
  └─────────┘     └──────────┘
       │          ┌──────────┐
       └─────────→│ Replica  │   Pros: simple, consistent writes
                  │ (R only) │   Cons: leader is bottleneck
                  └──────────┘

  Multi-Leader:
  ┌─────────┐ ←──→ ┌─────────┐
  │Leader A │      │Leader B │   Both accept writes
  │(Region 1)│      │(Region 2)│   Pros: low-latency multi-region
  └─────────┘      └─────────┘   Cons: conflict resolution needed

  Leaderless:
  ┌────────┐  ┌────────┐  ┌────────┐
  │ Node A │  │ Node B │  │ Node C │   Any node accepts writes
  └────────┘  └────────┘  └────────┘   Read quorum: R + W > N
                                        Pros: highly available
                                        Cons: complex, eventual consistency
```

### For Order Management (GCP)

```
  Cloud SQL (PostgreSQL) with read replicas:

  ┌───────────────┐     ┌──────────────────┐
  │ Primary       │────→│ Read Replica     │
  │ (us-central1) │     │ (us-central1)    │
  │ Writes + Reads│     │ Reads only       │
  └───────────────┘     └──────────────────┘
         │
         │ (cross-region for DR)
         ↓
  ┌──────────────────┐
  │ Read Replica     │
  │ (us-east1)       │
  │ Disaster Recovery│
  └──────────────────┘

  Write path: API → Primary (strong consistency)
  Read path:  API → Read Replica (eventual consistency, ~ms lag)
  DR:         Primary fails → promote cross-region replica (~minutes)
```

---

## 6. Distributed Transactions

### The Problem

Order creation spans multiple services:
1. Create order record (Order Service)
2. Reserve inventory (Inventory Service)
3. Charge payment (Payment Service)

Each service has its own database. You can't use a single SQL transaction.

### Option 1: Two-Phase Commit (2PC) — Usually Avoid

```
  Coordinator       Order DB       Inventory DB     Payment DB
       │                │               │               │
       │── PREPARE ────→│               │               │
       │── PREPARE ─────────────────→│               │
       │── PREPARE ──────────────────────────────→│
       │                │               │               │
       │←── READY ──────│               │               │
       │←── READY ──────────────────│               │
       │←── READY ─────────────────────────────│
       │                │               │               │
       │── COMMIT ─────→│               │               │
       │── COMMIT ──────────────────→│               │
       │── COMMIT ───────────────────────────────→│

  Problem: If coordinator crashes between PREPARE and COMMIT,
  all participants are BLOCKED holding locks. Terrible for availability.
```

**2PC is almost never the right answer** in a system design interview
(unless they specifically ask about it).

### Option 2: Saga Pattern — The Standard Answer

See Chapter 01 Phase 4 for the saga diagram. There are two flavors:

```
  Choreography (event-driven):
    Each service listens for events and acts independently
    ✅ Loosely coupled, no central coordinator
    ❌ Hard to track overall progress, debugging is difficult

  Orchestration (central coordinator):
    A saga orchestrator tells each service what to do
    ✅ Easy to track, clear flow, easier debugging
    ❌ Orchestrator is a single point of failure/coupling

  For order management, orchestration is usually better because:
  • Order flow has a clear sequential structure
  • Business needs visibility into "where is this order in the pipeline?"
  • Compensating actions need to be reliable and ordered
```

### Option 3: Outbox Pattern — Reliable Event Publishing

The **outbox pattern** ensures that database writes and event publishes are
atomic (without distributed transactions):

```
  Instead of:
    1. INSERT order into orders table     ← succeeds
    2. PUBLISH OrderCreated to Pub/Sub    ← might fail!
    Result: order exists but event never published

  Outbox pattern:
    1. BEGIN transaction
    2. INSERT order into orders table
    3. INSERT event into outbox table     ← same transaction!
    4. COMMIT

    Separate process (CDC or poller):
    5. Read outbox table
    6. Publish to Pub/Sub
    7. Mark outbox entry as published

  Result: if the transaction commits, the event is guaranteed
  to eventually be published (even if step 6 fails, it retries)
```

```sql
  CREATE TABLE outbox (
      id          BIGSERIAL PRIMARY KEY,
      event_type  TEXT NOT NULL,           -- 'OrderCreated'
      payload     JSONB NOT NULL,          -- event data
      published   BOOLEAN DEFAULT FALSE,
      created_at  TIMESTAMPTZ DEFAULT NOW()
  );

  -- Poller query (run every 100ms):
  SELECT * FROM outbox
  WHERE published = FALSE
  ORDER BY id
  LIMIT 100
  FOR UPDATE SKIP LOCKED;    -- concurrent pollers won't conflict
```

---

## 7. Time and Ordering in Distributed Systems

### The Fundamental Problem

In a distributed system, there is **no global clock**. Two events on different
machines can't reliably be ordered by wall-clock time.

```
  Node A clock: 12:00:00.001   ← "I processed order first"
  Node B clock: 12:00:00.000   ← "No, I processed order first"

  Who's right? Neither — clocks drift. NTP syncs to ~1-10ms accuracy.
  At Uber's QPS, 10ms is thousands of events.
```

### Solutions

**Logical clocks (Lamport timestamps)**: Increment a counter on every event.
If event A caused event B, A's counter < B's counter. Simple but only gives
partial ordering.

**Vector clocks**: Each node maintains a vector of counters for all nodes.
Gives causal ordering but vectors grow with number of nodes.

**Hybrid logical clocks (HLC)**: Combines wall-clock time with logical
counter. Used by CockroachDB and Google Spanner.

### For Interviews

You rarely need to explain clock algorithms. But know this:

> "We can't rely on wall-clock time for ordering events across services.
> Instead, I'd use an event sequence number per order — each event
> increments the sequence. This gives us total ordering within a single
> order, which is all we need for an order management system."

---

## 8. Failure Modes — Everything Fails

### Types of Failures

```
  Failure Type          Example                     Handling
  ────────────          ───────                     ────────
  Crash failure         Process dies                Restart, failover
  Omission failure      Message dropped             Retry with timeout
  Timing failure        Response too slow            Timeout + circuit breaker
  Byzantine failure     Node returns wrong data     Checksums, consensus
                        (rare in trusted networks)
```

### Network Failures

```
  1. Message loss:     Packet dropped → retry
  2. Message delay:    Packet arrives late → timeout
  3. Message duplicate: Retry causes duplicate → idempotency
  4. Partition:         Network split → CAP decision
  5. DNS failure:       Can't resolve hostname → cache DNS, fallback
```

### The "At Least Once" vs "At Most Once" vs "Exactly Once" Problem

```
  At most once:    Send message, don't retry if no response
                   → May lose messages
                   → Use for: logging, metrics (losing one is fine)

  At least once:   Retry until acknowledged
                   → May duplicate messages
                   → Use for: most operations (with idempotency)

  Exactly once:    Extremely hard in distributed systems
                   → Usually achieved via: at-least-once + idempotent receiver
                   → Use for: payments, order creation
```

**The practical answer**: Design for **at-least-once delivery** and make your
receivers **idempotent**. This effectively gives you exactly-once semantics
without the complexity of true exactly-once protocols.

---

## 9. Quick Reference Card

```
  ╔═══════════════════════════════════════════════════════════════════╗
  ║  DISTRIBUTED SYSTEMS — POCKET REFERENCE                          ║
  ╠═══════════════════════════════════════════════════════════════════╣
  ║                                                                   ║
  ║  CAP: During partition, choose Consistency OR Availability       ║
  ║  → Orders: CP (can't lose them)                                  ║
  ║  → Catalog: AP (stale is acceptable)                             ║
  ║                                                                   ║
  ║  Consistency Spectrum:                                            ║
  ║  Linearizable > Sequential > Causal > Eventual                   ║
  ║  (stronger = slower + less available)                             ║
  ║                                                                   ║
  ║  Partitioning: Split by how you QUERY, not how you INSERT        ║
  ║  Replication: Single-leader for most cases, read replicas        ║
  ║  Transactions: Saga pattern (not 2PC) for cross-service          ║
  ║  Events: Outbox pattern for reliable publishing                  ║
  ║  Delivery: At-least-once + idempotent receivers                  ║
  ║  Time: Don't trust wall clocks across nodes                      ║
  ║                                                                   ║
  ║  ALWAYS ASK: "What happens when this node/service/network fails?"║
  ║                                                                   ║
  ╚═══════════════════════════════════════════════════════════════════╝
```
