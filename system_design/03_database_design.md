# 03 — Database Design

> PostgreSQL deep dive for system design interviews.
> Schema design, indexing strategies, partitioning, and the patterns
> that matter at scale. Your interview stack includes PostgreSQL —
> know it well enough to justify every decision.

---

## Table of Contents

1. [PostgreSQL Architecture — How It Works](#1-postgresql-architecture--how-it-works)
2. [Schema Design Principles](#2-schema-design-principles)
3. [Indexing — The Performance Lever](#3-indexing--the-performance-lever)
4. [Query Patterns and EXPLAIN](#4-query-patterns-and-explain)
5. [Partitioning — When One Table Isn't Enough](#5-partitioning--when-one-table-isnt-enough)
6. [Connection Pooling](#6-connection-pooling)
7. [Read Replicas and Scaling](#7-read-replicas-and-scaling)
8. [PostgreSQL vs Other Databases — When to Choose What](#8-postgresql-vs-other-databases--when-to-choose-what)
9. [Order Management Schema — Complete Design](#9-order-management-schema--complete-design)
10. [Quick Reference Card](#10-quick-reference-card)

---

## 1. PostgreSQL Architecture — How It Works

### Process Model

```
  Client                PostgreSQL
    │                      │
    │── connect ──→   postmaster (main process)
    │                      │
    │                      ├── forks backend process (1 per connection)
    │←── session ──→  backend process
    │                      │
    │                      │   Shared Memory:
    │                      │   ┌────────────────────────┐
    │                      │   │ Shared Buffers (cache) │
    │                      │   │ WAL Buffers            │
    │                      │   │ Lock Tables            │
    │                      │   └────────────────────────┘
    │                      │
    │                      │   Disk:
    │                      │   ┌────────────────────────┐
    │                      │   │ Data Files (heap)      │
    │                      │   │ WAL Files (journal)    │
    │                      │   │ Index Files            │
    │                      │   └────────────────────────┘
```

### MVCC (Multi-Version Concurrency Control)

PostgreSQL doesn't use read locks. Instead, it keeps **multiple versions**
of each row:

```
  Transaction T1: UPDATE orders SET status = 'shipped' WHERE id = 1;

  Before:
  ┌────────────────────────────────────────┐
  │ Row: id=1, status='confirmed'          │
  │ xmin=100 (created by txn 100)         │
  │ xmax=∞   (not deleted yet)            │
  └────────────────────────────────────────┘

  After:
  ┌────────────────────────────────────────┐
  │ Old version: id=1, status='confirmed'  │  ← still visible to
  │ xmin=100, xmax=200                    │     concurrent reads
  └────────────────────────────────────────┘
  ┌────────────────────────────────────────┐
  │ New version: id=1, status='shipped'    │  ← visible only after
  │ xmin=200, xmax=∞                      │     T1 commits
  └────────────────────────────────────────┘
```

**Why it matters for interviews**: MVCC means readers don't block writers
and writers don't block readers. But old row versions accumulate →
**VACUUM** is needed to clean them up. Dead tuple bloat is a common
PostgreSQL production issue.

### WAL (Write-Ahead Log)

Every change goes to WAL **before** data files. This guarantees durability:

```
  Write path:
  1. Write change to WAL buffer (memory)
  2. Flush WAL buffer to WAL file on disk (fsync)
  3. Return "committed" to client
  4. Eventually, write dirty pages from shared buffers to data files (checkpoint)

  If crash after step 2: replay WAL to recover → no data loss
  If crash after step 3: same — WAL has the record
```

**WAL is also used for replication**: stream WAL records to read replicas.

---

## 2. Schema Design Principles

### Normalize First, Denormalize When Needed

```
  Normalized (3NF):                    Denormalized:
  ┌─────────┐   ┌──────────────┐      ┌──────────────────────────┐
  │ orders  │   │ order_items  │      │ orders_denormalized      │
  │─────────│   │──────────────│      │──────────────────────────│
  │ id      │←─→│ order_id     │      │ id                       │
  │ cust_id │   │ product_id   │      │ cust_id                  │
  │ status  │   │ quantity     │      │ status                   │
  │ total   │   │ unit_price   │      │ items_json (JSONB)       │
  └─────────┘   └──────────────┘      │ total                    │
                                       └──────────────────────────┘

  Normalized:
  ✅ No data duplication
  ✅ Easy to update (change item price in one place)
  ❌ JOINs for every order query

  Denormalized:
  ✅ Single query, no JOINs
  ✅ Faster reads
  ❌ Update anomalies (must update items everywhere)
  ❌ More storage
```

**For order management**: Start normalized. Denormalize the **read path**
if performance requires it (e.g., store a JSONB snapshot of order items
in the orders table for the "order details" API, but keep normalized
order_items for inventory and reporting queries).

### Use the Right Data Types

```
  Column Purpose          Type                   Why
  ──────────────          ────                   ───
  Primary key             UUID or BIGSERIAL      UUID: distributed-safe, no sequence bottleneck
                                                  BIGSERIAL: smaller, faster, sortable by time
  Money                   NUMERIC(12,2)          NEVER use FLOAT (rounding errors)
  Timestamps              TIMESTAMPTZ            Always store with timezone
  Status/enum             TEXT + CHECK           Flexible, easy to migrate
  Flexible attributes     JSONB                  Indexed, queryable, schema-flexible
  IP addresses            INET                   Built-in validation and operators
```

### JSONB — When and When Not

```
  ✅ Use JSONB for:
  • Order metadata (gift message, special instructions)
  • Event payloads (flexible schema per event type)
  • API response caching
  • Feature flags per entity

  ❌ Don't use JSONB for:
  • Core business fields you query often (use columns)
  • Relationships (use foreign keys)
  • Fields that need strong type validation
  • Fields that need to be aggregated (SUM, AVG)
```

---

## 3. Indexing — The Performance Lever

### Index Types

```
  Type          When to Use                    How It Works
  ──────        ─────────────                  ─────────────
  B-Tree        Default. Equality and range    Balanced tree, O(log n) lookup
                WHERE status = 'shipped'
                WHERE created_at > '2024-01-01'

  Hash          Equality only (rare)           Hash table, O(1) lookup
                WHERE id = 'abc-123'           B-Tree is usually fast enough

  GIN           JSONB fields, full-text,       Inverted index
                array contains                 Maps values → rows
                WHERE tags @> '{"vip"}'

  GiST          Geospatial, range types        Generalized search tree
                WHERE location <@ box          Used with PostGIS

  BRIN          Very large tables with         Block Range Index
                naturally ordered data          Tiny index, good for time-series
                WHERE created_at > '2024-01'   Orders table sorted by time
```

### Indexing Strategies for Order Management

```sql
  -- Primary lookups
  CREATE INDEX idx_orders_customer ON orders (customer_id);
  CREATE INDEX idx_orders_status ON orders (status);

  -- Compound index for common query pattern
  CREATE INDEX idx_orders_customer_status ON orders (customer_id, status);
  -- Serves: WHERE customer_id = ? AND status = ?
  -- Also serves: WHERE customer_id = ? (leftmost prefix)
  -- Does NOT serve: WHERE status = ? (need separate index)

  -- Partial index — only index what you query
  CREATE INDEX idx_orders_active ON orders (customer_id)
  WHERE status NOT IN ('delivered', 'cancelled', 'returned');
  -- Smaller index, faster scans — most orders are completed

  -- Index for time-range queries
  CREATE INDEX idx_orders_created ON orders (created_at DESC);

  -- JSONB index for order metadata
  CREATE INDEX idx_orders_metadata ON orders USING GIN (metadata);
```

### The Index Selection Thought Process

```
  1. What queries will run most often?
     → Index those columns

  2. What's the cardinality (number of distinct values)?
     → High cardinality (customer_id): good for B-Tree
     → Low cardinality (status): consider partial index

  3. What columns appear together in WHERE clauses?
     → Compound index (leftmost prefix rule applies)

  4. Are there range queries?
     → Put range column LAST in compound index
     → CREATE INDEX idx ON orders (status, created_at)
        serves: WHERE status = 'shipped' AND created_at > '2024-01-01'

  5. Does the query need to sort?
     → Index can serve ORDER BY if columns match
     → Eliminates sort step entirely

  6. Is the table very large with time-ordered data?
     → Consider BRIN index (tiny, great for time-series)
```

### Index Anti-Patterns

```
  ❌ Indexing every column
     → Each index costs write performance (must update on INSERT/UPDATE)
     → Each index consumes storage and memory

  ❌ Forgetting compound index order
     → INDEX (a, b) serves WHERE a = ? AND b = ?
     → INDEX (a, b) serves WHERE a = ?
     → INDEX (a, b) does NOT serve WHERE b = ?

  ❌ Not using partial indexes
     → If 90% of orders are 'delivered', indexing all statuses wastes space
     → Partial index on active orders only: smaller, faster

  ❌ Missing covering indexes
     → If query only needs id and status, include both in index
     → CREATE INDEX idx ON orders (customer_id) INCLUDE (status)
     → Index-only scan: no need to read the table at all
```

---

## 4. Query Patterns and EXPLAIN

### Reading EXPLAIN Output

```sql
  EXPLAIN ANALYZE
  SELECT * FROM orders WHERE customer_id = 'cust-123' AND status = 'shipped';

  -- Output:
  Index Scan using idx_orders_customer_status on orders
    Index Cond: (customer_id = 'cust-123' AND status = 'shipped')
    Rows: 5 (estimated), 5 (actual)
    Time: 0.1ms
```

### Query Patterns to Know

```
  Pattern              SQL                          Index Strategy
  ───────              ───                          ──────────────
  Point lookup         WHERE id = ?                 Primary key (automatic)
  Customer's orders    WHERE cust_id = ?            B-Tree on cust_id
  Status filter        WHERE status = ?             Partial index
  Time range           WHERE created > ?            B-Tree or BRIN
  Pagination           ORDER BY id LIMIT 20         Cursor-based pagination
  Aggregation          COUNT(*) WHERE status = ?    Partial index
```

### Pagination — Cursor vs Offset

```sql
  -- ❌ Offset pagination (slow at high offsets)
  SELECT * FROM orders ORDER BY id LIMIT 20 OFFSET 10000;
  -- Must scan and discard 10,000 rows!

  -- ✅ Cursor pagination (constant time)
  SELECT * FROM orders
  WHERE id > 'last-seen-id'
  ORDER BY id
  LIMIT 20;
  -- Seeks directly to the position using the index
```

**Always use cursor-based pagination in system design interviews.**

---

## 5. Partitioning — When One Table Isn't Enough

### PostgreSQL Native Partitioning

```sql
  -- Range partitioning by date (most common for order tables)
  CREATE TABLE orders (
      id          UUID PRIMARY KEY,
      customer_id UUID NOT NULL,
      status      TEXT NOT NULL,
      total       NUMERIC(12,2),
      created_at  TIMESTAMPTZ NOT NULL
  ) PARTITION BY RANGE (created_at);

  -- Create monthly partitions
  CREATE TABLE orders_2024_01 PARTITION OF orders
      FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
  CREATE TABLE orders_2024_02 PARTITION OF orders
      FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');
  -- ... create ahead of time (automate with pg_partman)

  -- Queries with created_at filter only scan relevant partitions
  SELECT * FROM orders WHERE created_at >= '2024-06-01';
  -- Only scans orders_2024_06 and later — partition pruning
```

### When to Partition

```
  Partition when:
  ✅ Table > 100GB and growing
  ✅ Queries almost always filter by the partition key (date)
  ✅ You need to archive/drop old data quickly (DROP TABLE is instant)
  ✅ Maintenance (VACUUM, REINDEX) becomes slow on large table

  Don't partition when:
  ❌ Table < 10GB (single table is fine)
  ❌ Queries frequently cross partition boundaries
  ❌ You'd need more than ~1000 partitions (management overhead)
```

---

## 6. Connection Pooling

### The Problem

PostgreSQL forks a **process per connection** (~10MB each). 100 connections = 1GB
just for connection processes. Cloud SQL limits connections (e.g., 500 for a
small instance).

```
  Without pooling:
  100 application pods × 10 connections each = 1000 connections
  → Exceeds Cloud SQL limit → connection refused errors

  With pooling (PgBouncer / Cloud SQL Auth Proxy):
  100 pods → PgBouncer → 50 actual PostgreSQL connections
  → Multiplexes application connections over fewer DB connections
```

### Pooling Modes

```
  Mode              How It Works                  When to Use
  ──────            ────────────                  ──────────
  Session           Connection bound to session   Temp tables, LISTEN/NOTIFY
  Transaction       Connection returned after txn Default — best for most cases
  Statement         Connection returned after stmt Only simple queries, no txns
```

### For GCP

```
  Application → Cloud SQL Auth Proxy → Cloud SQL (PostgreSQL)
                  (handles auth +       (managed PostgreSQL)
                   connection pooling)

  Or use AlloyDB for PostgreSQL (Google's PostgreSQL-compatible DB)
  with built-in connection pooling.
```

---

## 7. Read Replicas and Scaling

### Read Scaling Pattern

```
  Write Path:                      Read Path:
  API → Primary DB                 API → Read Replica
  (strong consistency)             (eventual consistency, ~ms lag)

  ┌───────────────┐     streaming     ┌──────────────────┐
  │   Primary     │────replication───→│  Read Replica    │
  │ (writes+reads)│                   │  (reads only)    │
  └───────────────┘                   └──────────────────┘

  Split traffic by query type:
  • Order creation   → Primary (must be consistent)
  • Order status     → Replica (stale by ~100ms is fine)
  • Order history    → Replica (historical data doesn't change)
  • Analytics        → Replica (never hit primary with reports)
```

### Replication Lag

```
  Primary: INSERT order #1001 at 12:00:00.000
  Replica: sees order #1001 at 12:00:00.050 (50ms lag)

  Problem: Customer places order → immediate redirect to "My Orders"
  → reads from replica → order not there yet!

  Solutions:
  1. Read-your-writes: route to primary for N seconds after write
  2. Session stickiness: same user always hits same node
  3. Version token: client sends version, replica checks if caught up
```

---

## 8. PostgreSQL vs Other Databases — When to Choose What

```
  Need                          Best Choice          Why Not PostgreSQL
  ──────                        ──────────           ──────────────────
  Relational + ACID             PostgreSQL ✅         (it IS the best choice)
  Key-value cache               Redis/Memorystore    PG too slow for sub-ms
  Document store                MongoDB or PG JSONB  PG JSONB is often enough
  Time-series (metrics)         InfluxDB/TimescaleDB PG works but not optimized
  Full-text search              Elasticsearch        PG tsvector works for simple
  Graph queries                 Neo4j                PG recursive CTEs for simple
  Global strong consistency     Google Spanner       PG is single-region
  Wide column / huge scale      Bigtable             PG doesn't scale that wide
  Analytics / OLAP              BigQuery             PG row-store too slow
```

**For your interview**: PostgreSQL + Redis covers 90% of use cases.
Add BigQuery for analytics. Add Pub/Sub for events. That's usually enough.

---

## 9. Order Management Schema — Complete Design

```sql
  -- Core order table
  CREATE TABLE orders (
      id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      customer_id     UUID NOT NULL,
      status          TEXT NOT NULL DEFAULT 'created',
      total_amount    NUMERIC(12, 2) NOT NULL,
      currency        TEXT NOT NULL DEFAULT 'USD',
      shipping_addr   JSONB NOT NULL,
      promise_eta     TIMESTAMPTZ,
      idempotency_key TEXT UNIQUE,
      created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

      CONSTRAINT valid_status CHECK (
          status IN ('created', 'confirmed', 'picked', 'shipped',
                     'delivered', 'cancelled', 'failed', 'returned')
      )
  ) PARTITION BY RANGE (created_at);

  -- Order line items
  CREATE TABLE order_items (
      id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
      order_id    UUID NOT NULL REFERENCES orders(id),
      product_id  UUID NOT NULL,
      quantity    INT NOT NULL CHECK (quantity > 0),
      unit_price  NUMERIC(12, 2) NOT NULL,
      subtotal    NUMERIC(12, 2) GENERATED ALWAYS AS (quantity * unit_price) STORED
  );

  -- Event sourcing table — every state change is recorded
  CREATE TABLE order_events (
      id          BIGSERIAL PRIMARY KEY,
      order_id    UUID NOT NULL,
      event_type  TEXT NOT NULL,
      payload     JSONB NOT NULL DEFAULT '{}',
      actor       TEXT NOT NULL,            -- 'system', 'customer', 'ops'
      created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
  ) PARTITION BY RANGE (created_at);

  -- Outbox for reliable event publishing
  CREATE TABLE outbox (
      id          BIGSERIAL PRIMARY KEY,
      event_type  TEXT NOT NULL,
      payload     JSONB NOT NULL,
      published   BOOLEAN NOT NULL DEFAULT FALSE,
      created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
  );

  -- Indexes
  CREATE INDEX idx_orders_customer ON orders (customer_id);
  CREATE INDEX idx_orders_status_active ON orders (status)
      WHERE status NOT IN ('delivered', 'cancelled', 'returned');
  CREATE INDEX idx_order_items_order ON order_items (order_id);
  CREATE INDEX idx_order_events_order ON order_events (order_id, created_at);
  CREATE INDEX idx_outbox_unpublished ON outbox (id)
      WHERE published = FALSE;
```

---

## 10. Quick Reference Card

```
  ╔═══════════════════════════════════════════════════════════════════╗
  ║  DATABASE DESIGN — POCKET REFERENCE                              ║
  ╠═══════════════════════════════════════════════════════════════════╣
  ║                                                                   ║
  ║  PostgreSQL strengths: ACID, JSONB, partitioning, extensions     ║
  ║  PostgreSQL limits: single-region, ~5TB comfortable, ~10K QPS    ║
  ║                                                                   ║
  ║  Schema: normalize first, denormalize read path if needed        ║
  ║  Money: NUMERIC(12,2) — NEVER float                              ║
  ║  Timestamps: always TIMESTAMPTZ                                  ║
  ║  PKs: UUID (distributed) or BIGSERIAL (simple)                   ║
  ║                                                                   ║
  ║  Indexes: B-Tree (default), GIN (JSONB), BRIN (time-series)     ║
  ║  Compound index: leftmost prefix rule                            ║
  ║  Partial index: only index what you query                        ║
  ║  Covering index: INCLUDE for index-only scans                    ║
  ║                                                                   ║
  ║  Partitioning: range by date for orders                          ║
  ║  Pagination: cursor-based (not OFFSET)                           ║
  ║  Pooling: PgBouncer or Cloud SQL Proxy, transaction mode         ║
  ║  Replication: read replicas for read scaling + DR                ║
  ║                                                                   ║
  ║  MVCC: readers don't block writers (but watch for dead tuples)   ║
  ║  WAL: durability guarantee + replication mechanism               ║
  ║                                                                   ║
  ╚═══════════════════════════════════════════════════════════════════╝
```
