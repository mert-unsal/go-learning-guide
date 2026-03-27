# 09 — Order Management Domain

> The domain-specific chapter. Everything you need to design an order
> management system, order promise service, or fulfillment pipeline.
> This ties together concepts from every other chapter.

---

## Table of Contents

1. [Order Management — The Big Picture](#1-order-management--the-big-picture)
2. [Order Lifecycle — State Machine](#2-order-lifecycle--state-machine)
3. [Order Promise — Delivery ETA](#3-order-promise--delivery-eta)
4. [Inventory Management](#4-inventory-management)
5. [Payment Integration](#5-payment-integration)
6. [Fulfillment Pipeline](#6-fulfillment-pipeline)
7. [Order Splitting and Multi-Warehouse](#7-order-splitting-and-multi-warehouse)
8. [Cancellation and Returns](#8-cancellation-and-returns)
9. [Full Architecture — Putting It All Together](#9-full-architecture--putting-it-all-together)
10. [Interview Walkthrough — Design an Order System](#10-interview-walkthrough--design-an-order-system)

---

## 1. Order Management — The Big Picture

### What Is an Order Management System (OMS)?

An OMS is the central nervous system of e-commerce. It coordinates:

```
  Customer places order
       │
       ├── Validate items are available (Inventory)
       ├── Calculate delivery promise (Order Promise)
       ├── Process payment (Payment)
       ├── Route to warehouse (Fulfillment)
       ├── Track shipment (Logistics)
       ├── Handle returns/refunds (Returns)
       └── Notify customer at each step (Notifications)

  Surrounding systems:
  ┌─────────────────────────────────────────────────────────────┐
  │                                                             │
  │   Product Catalog ──→ ┌─────────┐ ──→ Fulfillment/WMS     │
  │   Pricing Engine ───→ │  ORDER  │ ──→ Shipping/Logistics  │
  │   Customer Service ──→│ MANAGE- │ ──→ Payment Processor   │
  │   Promotions ────────→│  MENT   │ ──→ Notification Service│
  │   Fraud Detection ──→ │ SYSTEM  │ ──→ Analytics/BI        │
  │   Tax Service ──────→ └─────────┘ ──→ Customer Service    │
  │                                                             │
  └─────────────────────────────────────────────────────────────┘
```

### Key Entities

```
  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
  │    Order     │     │  Order Item  │     │   Shipment   │
  ├──────────────┤     ├──────────────┤     ├──────────────┤
  │ id           │←──→ │ order_id     │     │ order_id     │
  │ customer_id  │     │ product_id   │     │ warehouse_id │
  │ status       │     │ quantity     │     │ carrier      │
  │ total        │     │ unit_price   │     │ tracking_no  │
  │ promise_eta  │     │ status       │     │ shipped_at   │
  │ payment_ref  │     │ warehouse_id │     │ delivered_at │
  └──────────────┘     └──────────────┘     └──────────────┘

  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
  │   Payment    │     │  Inventory   │     │ Order Event  │
  ├──────────────┤     ├──────────────┤     ├──────────────┤
  │ order_id     │     │ product_id   │     │ order_id     │
  │ amount       │     │ warehouse_id │     │ event_type   │
  │ method       │     │ available    │     │ payload      │
  │ status       │     │ reserved     │     │ timestamp    │
  │ provider_ref │     │ on_hand      │     │ actor        │
  └──────────────┘     └──────────────┘     └──────────────┘
```

---

## 2. Order Lifecycle — State Machine

### Complete State Machine

```
  ┌──────────┐
  │ PENDING  │ ← Customer submitted order (items in cart → order)
  └────┬─────┘
       │ validate items + address + fraud check
       ↓
  ┌──────────┐   payment     ┌──────────┐
  │ CREATED  │────failed───→│  FAILED  │ (terminal)
  └────┬─────┘              └──────────┘
       │ payment authorized
       ↓
  ┌──────────┐   cancel     ┌───────────┐
  │CONFIRMED │───request──→│ CANCELLING│──→ CANCELLED (terminal)
  └────┬─────┘              └───────────┘
       │ inventory reserved + assigned to warehouse
       ↓
  ┌──────────┐
  │ PICKING  │ ← Items being picked from warehouse shelves
  └────┬─────┘
       │ all items picked and packed
       ↓
  ┌──────────┐
  │  PACKED  │ ← Ready for carrier pickup
  └────┬─────┘
       │ carrier picks up, tracking number assigned
       ↓
  ┌──────────┐
  │ SHIPPED  │ ← In transit with carrier
  └────┬─────┘
       │ carrier confirms delivery
       ↓
  ┌──────────┐   return    ┌──────────┐
  │DELIVERED │──request──→│ RETURNED │ (terminal)
  └──────────┘             └──────────┘
```

### State Transition Rules

```
  From          Allowed To              Trigger
  ────          ──────────              ───────
  PENDING       CREATED, FAILED         Payment authorization result
  CREATED       CONFIRMED, FAILED       Payment capture + inventory check
  CONFIRMED     PICKING, CANCELLING     Warehouse assignment
  PICKING       PACKED, CANCELLING      Pick completion
  PACKED        SHIPPED                 Carrier pickup
  SHIPPED       DELIVERED               Delivery confirmation
  DELIVERED     RETURNED                Return request (within window)

  Terminal states: FAILED, CANCELLED, DELIVERED, RETURNED
  (no outgoing transitions)
```

### Implementation: Event-Sourced State Machine

```go
  // Every state transition is an event
  type OrderEvent struct {
      OrderID   string    `json:"order_id"`
      EventType string    `json:"event_type"`
      Payload   any       `json:"payload"`
      Actor     string    `json:"actor"`
      Timestamp time.Time `json:"timestamp"`
  }

  // Event types map to transitions
  // PaymentAuthorized:   PENDING → CREATED
  // PaymentCaptured:     CREATED → CONFIRMED
  // WarehouseAssigned:   CONFIRMED → PICKING
  // AllItemsPicked:      PICKING → PACKED
  // CarrierPickedUp:     PACKED → SHIPPED
  // DeliveryConfirmed:   SHIPPED → DELIVERED
  // PaymentFailed:       PENDING/CREATED → FAILED
  // CancellationApproved: CONFIRMED/PICKING → CANCELLED
  // ReturnApproved:      DELIVERED → RETURNED
```

---

## 3. Order Promise — Delivery ETA

### What Is Order Promise?

Order Promise calculates **when the customer will receive their order** and
commits to it. This ETA is shown on the product page, in checkout, and in
order confirmation.

```
  Customer sees: "Order by 2 PM, get it by Thursday"
                  ↑                    ↑
                  cutoff time          promise ETA

  The promise must account for:
  ├── Inventory availability (is it in stock? which warehouse?)
  ├── Warehouse processing time (pick + pack: ~4-8 hours)
  ├── Carrier transit time (depends on origin → destination)
  ├── Carrier cutoff times (last pickup at 5 PM)
  ├── Business days (no delivery on weekends/holidays)
  └── Current capacity (can the warehouse handle it today?)
```

### Promise Calculation Architecture

```
  ┌─────────────┐
  │ Product Page│
  │ / Checkout  │
  └──────┬──────┘
         │ GET /promise?product=X&zip=12345
         ↓
  ┌─────────────────────────────┐
  │     Promise Service         │
  │                             │
  │  1. Check inventory         │──→ Inventory Service
  │     (which warehouse?)      │     "product X is in warehouse A, B"
  │                             │
  │  2. Calculate processing    │──→ Warehouse Capacity Service
  │     (when can they ship?)   │     "warehouse A can ship by tomorrow"
  │                             │
  │  3. Calculate transit       │──→ Carrier Service
  │     (how long to deliver?)  │     "UPS Ground: A→12345 = 3 days"
  │                             │
  │  4. Apply business rules    │
  │     (cutoffs, holidays)     │
  │                             │
  │  Result: "Delivers Thursday"│
  └─────────────────────────────┘
```

### Promise Data Model

```sql
  -- Promise rules per warehouse-carrier combination
  CREATE TABLE promise_rules (
      warehouse_id    UUID NOT NULL,
      carrier_id      UUID NOT NULL,
      destination_zip TEXT NOT NULL,
      processing_hours INT NOT NULL,     -- warehouse pick+pack time
      transit_days     INT NOT NULL,      -- carrier transit time
      cutoff_time      TIME NOT NULL,     -- order-by time for same-day processing
      PRIMARY KEY (warehouse_id, carrier_id, destination_zip)
  );

  -- Carrier transit time zones
  CREATE TABLE carrier_zones (
      carrier_id       UUID NOT NULL,
      origin_zip_range TEXT NOT NULL,     -- "90000-90999"
      dest_zip_range   TEXT NOT NULL,     -- "10000-10999"
      transit_days     INT NOT NULL,
      service_level    TEXT NOT NULL,     -- 'ground', 'express', 'overnight'
      PRIMARY KEY (carrier_id, origin_zip_range, dest_zip_range, service_level)
  );
```

### Promise Accuracy and SLOs

```
  Key metrics:
  ├── Promise accuracy: % of orders delivered by promised ETA
  │   Target: >95% on time
  │
  ├── Promise availability: can we show a promise on the product page?
  │   Target: 99.9% (fallback: "Delivers in 3-5 business days")
  │
  └── Promise latency: how fast can we calculate?
      Target: P99 < 200ms (it's on the critical checkout path)

  If promise is wrong:
  ├── Too optimistic: customer angry, trust lost, support tickets
  ├── Too pessimistic: customer buys from competitor
  └── Best: slightly conservative (under-promise, over-deliver)
```

### Caching Strategy for Promise

```
  Product page traffic >> order creation traffic
  Promise calculation involves multiple service calls (expensive)
  → Cache aggressively

  Cache key: product_id + destination_zip + service_level
  Cache TTL: 5-15 minutes (inventory can change)
  Cache location: Redis/Memorystore

  Invalidation:
  ├── Inventory change event → invalidate affected products
  ├── Carrier disruption → invalidate affected zones
  └── TTL expiry (safety net)
```

---

## 4. Inventory Management

### Inventory Types

```
  ┌─────────────────────────────────────────────────────────┐
  │                    Inventory Levels                     │
  │                                                         │
  │  On-hand:     Physical items in warehouse               │
  │  Reserved:    Allocated to confirmed orders (not shipped)│
  │  Available:   on_hand - reserved (what can be sold)     │
  │  In-transit:  Being shipped from supplier to warehouse  │
  │  Backorder:   Ordered but supplier hasn't shipped yet   │
  │                                                         │
  │  available = on_hand - reserved                         │
  │  sellable = available + in_transit (if you allow pre-orders)│
  └─────────────────────────────────────────────────────────┘
```

### The Oversell Problem

```
  Scenario: 1 item left in stock, 2 customers checkout simultaneously

  Without protection:
    Customer A: SELECT available FROM inventory WHERE product = X → 1
    Customer B: SELECT available FROM inventory WHERE product = X → 1
    Customer A: UPDATE inventory SET reserved = reserved + 1 → ok
    Customer B: UPDATE inventory SET reserved = reserved + 1 → ok!
    Result: reserved = 2, on_hand = 1 → OVERSOLD

  With database-level protection:
    UPDATE inventory
    SET reserved = reserved + 1
    WHERE product_id = 'X'
    AND available > 0;      ← atomic check-and-update

    If affected rows = 0 → out of stock (return error)
    If affected rows = 1 → successfully reserved

  Or with SELECT FOR UPDATE:
    BEGIN;
    SELECT * FROM inventory WHERE product_id = 'X' FOR UPDATE;
    -- Row is now locked — other transactions wait
    UPDATE inventory SET reserved = reserved + 1;
    COMMIT;
```

### Inventory Reservation Flow

```
  Order created
      │
      ├── Reserve inventory (atomic decrement available)
      │   └── If out of stock → FAIL order
      │
      ├── Payment charged
      │   └── If payment fails → RELEASE reservation
      │
      ├── Order shipped
      │   └── Decrement on_hand (item physically left warehouse)
      │
      └── Order cancelled
          └── RELEASE reservation (increment available back)
```

---

## 5. Payment Integration

### Payment Flow

```
  Two-phase payment (industry standard):

  Phase 1: AUTHORIZE (at checkout)
  ┌──────────┐     ┌──────────────┐     ┌──────────────┐
  │ Order Svc│────→│ Payment Svc  │────→│ Stripe/Adyen │
  │          │     │              │     │ (processor)  │
  │ "auth    │     │ "can this    │     │ "yes, funds  │
  │  $99.99" │     │  card pay?"  │     │  available"  │
  └──────────┘     └──────────────┘     └──────────────┘
  → No money moves yet! Just a hold on the customer's card.

  Phase 2: CAPTURE (at shipment)
  ┌──────────┐     ┌──────────────┐     ┌──────────────┐
  │ Order Svc│────→│ Payment Svc  │────→│ Stripe/Adyen │
  │          │     │              │     │              │
  │ "capture │     │ "charge the  │     │ "money moved │
  │  $99.99" │     │  authorized  │     │  to merchant"│
  └──────────┘     │  amount"     │     └──────────────┘
                   └──────────────┘

  Why two phases?
  ├── Don't charge until you ship (legal requirement in many regions)
  ├── Easy cancellation: just void the authorization
  ├── Partial capture: ship 2 of 3 items, capture only those
  └── Auth expires (typically 7 days) — capture before it expires
```

### Payment Idempotency

```
  Every payment API call needs an idempotency key:

  POST /v1/payments/authorize
  Idempotency-Key: order-456-auth-1
  {
    "amount": 9999,
    "currency": "USD",
    "payment_method": "pm_xxx"
  }

  If the same Idempotency-Key is sent again:
  → Stripe returns the original result (no double-charge)
```

### PCI-DSS Compliance

```
  NEVER store:
  ├── Card numbers
  ├── CVV/CVC codes
  ├── Full magnetic stripe data
  └── PIN numbers

  ALWAYS:
  ├── Use a payment processor (Stripe, Adyen, Braintree)
  ├── Store only payment tokens/references
  ├── Use processor's hosted payment form (reduces PCI scope)
  └── Tokenize card data at the edge (client-side SDK)
```

---

## 6. Fulfillment Pipeline

### From Order to Delivery

```
  Order CONFIRMED
       │
       ↓
  ┌────────────────┐
  │ Warehouse      │ Select optimal warehouse based on:
  │ Assignment     │ ├── Inventory availability
  └────┬───────────┘ ├── Distance to customer (shipping cost)
       │             └── Current capacity (load balancing)
       ↓
  ┌────────────────┐
  │ Pick List      │ Generate picking instructions for warehouse worker
  │ Generation     │ (aisle, shelf, bin location)
  └────┬───────────┘
       │
       ↓
  ┌────────────────┐
  │ Picking        │ Worker picks items from shelves
  │                │ Scan barcode to confirm correct item
  └────┬───────────┘
       │
       ↓
  ┌────────────────┐
  │ Packing        │ Pack items, generate shipping label
  │                │ Select box size, add packing materials
  └────┬───────────┘
       │
       ↓
  ┌────────────────┐
  │ Carrier        │ Carrier picks up package
  │ Handoff        │ Tracking number assigned
  └────┬───────────┘
       │
       ↓
  ┌────────────────┐
  │ In Transit     │ Carrier provides tracking updates
  │ Tracking       │ (via webhooks or polling API)
  └────┬───────────┘
       │
       ↓
  ┌────────────────┐
  │ Delivered      │ Carrier confirms delivery
  │                │ (photo, signature, GPS)
  └────────────────┘
```

---

## 7. Order Splitting and Multi-Warehouse

### When Orders Split

```
  Customer orders: Item A + Item B + Item C

  Item A: available in Warehouse WEST
  Item B: available in Warehouse EAST
  Item C: available in both WEST and EAST

  Options:
  1. Ship all from WEST (Item B ships cross-country — slow, expensive)
  2. Ship all from EAST (Item A ships cross-country — slow, expensive)
  3. Split: A+C from WEST, B from EAST (fastest, but two shipments)

  Decision factors:
  ├── Customer promise (which option meets the ETA?)
  ├── Shipping cost (two shipments > one, but faster)
  ├── Customer experience (multiple packages vs one delayed package)
  └── Business rules (company policy on splitting)
```

### Split Order Data Model

```
  Order #456
  ├── Shipment #456-A (Warehouse WEST)
  │   ├── Item A × 1
  │   └── Item C × 1
  │
  └── Shipment #456-B (Warehouse EAST)
      └── Item B × 1

  Order status: derived from shipment statuses
  ├── All shipments delivered → Order DELIVERED
  ├── Any shipment shipped → Order PARTIALLY_SHIPPED
  ├── All shipments picked → Order PACKED
  └── Mixed states → show per-shipment status to customer
```

---

## 8. Cancellation and Returns

### Cancellation Rules

```
  Can cancel?
  ├── PENDING/CREATED:     YES (void payment auth, release inventory)
  ├── CONFIRMED:           YES (void payment, release inventory)
  ├── PICKING:             MAYBE (if not picked yet, cancel; else too late)
  ├── PACKED:              NO (too late, suggest return after delivery)
  ├── SHIPPED:             NO (in transit, suggest return after delivery)
  └── DELIVERED:           NO (use return flow instead)

  Cancellation is a saga (compensation pattern):
  1. Mark order CANCELLING
  2. Void payment authorization → PaymentVoided event
  3. Release inventory reservation → InventoryReleased event
  4. Both successful → Mark order CANCELLED
  5. If either fails → retry with exponential backoff
```

### Return Flow

```
  Customer requests return (within return window, e.g., 30 days)
       │
       ├── Validate return eligibility
       │   ├── Within return window?
       │   ├── Item returnable? (not customized, not final-sale)
       │   └── Reason provided?
       │
       ├── Generate return label (or schedule pickup)
       │
       ├── Customer ships item back
       │
       ├── Warehouse receives return
       │   ├── Inspect item condition
       │   └── Restock if sellable
       │
       └── Process refund
           ├── Full refund: original payment method
           ├── Partial refund: deduct restocking fee
           └── Store credit: gift card/credit
```

---

## 9. Full Architecture — Putting It All Together

```
  ┌─────────────────────────────────────────────────────────────────────┐
  │                        GCP Infrastructure                          │
  │                                                                     │
  │  ┌──────────────┐                                                  │
  │  │ Cloud Load   │                                                  │
  │  │ Balancer     │                                                  │
  │  └──────┬───────┘                                                  │
  │         │                                                           │
  │  ┌──────┴───────┐     ┌──────────────┐     ┌──────────────┐       │
  │  │ API Gateway  │     │ Promise Svc  │     │ Inventory Svc│       │
  │  │ (Cloud Run)  │────→│ (Cloud Run)  │     │ (Cloud Run)  │       │
  │  └──────┬───────┘     └──────────────┘     └──────────────┘       │
  │         │                                                           │
  │  ┌──────┴───────┐                                                  │
  │  │ Order Svc    │────→ Cloud SQL (PostgreSQL)                      │
  │  │ (Cloud Run)  │        ├── Primary (writes)                      │
  │  └──────┬───────┘        └── Read Replica (reads)                  │
  │         │                                                           │
  │         │ events                                                    │
  │         ↓                                                           │
  │  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐       │
  │  │  Pub/Sub     │────→│ Fulfillment  │     │ Notification │       │
  │  │ (events)     │────→│ Svc          │     │ Svc          │       │
  │  │              │────→│ (Cloud Run)  │     │ (Cloud Run)  │       │
  │  └──────────────┘     └──────────────┘     └──────────────┘       │
  │                                                                     │
  │  ┌──────────────┐     ┌──────────────┐                             │
  │  │ Memorystore  │     │ BigQuery     │ ← analytics events         │
  │  │ (Redis)      │     │ (analytics)  │                             │
  │  │ (cache)      │     └──────────────┘                             │
  │  └──────────────┘                                                  │
  │                                                                     │
  │  All infrastructure managed by Terraform                           │
  └─────────────────────────────────────────────────────────────────────┘
```

### Service Responsibilities

```
  Service          Owns                    Database
  ───────          ────                    ────────
  Order Svc        Order lifecycle         Cloud SQL (PostgreSQL)
  Promise Svc      Delivery ETA calc       Cloud SQL + Redis cache
  Inventory Svc    Stock levels            Cloud SQL (PostgreSQL)
  Payment Svc      Payment processing      Cloud SQL + Stripe API
  Fulfillment Svc  Warehouse operations    Cloud SQL (PostgreSQL)
  Notification Svc Email/SMS/Push          Pub/Sub (stateless)
```

---

## 10. Interview Walkthrough — Design an Order System

### Complete 45-Minute Walkthrough

**Phase 1: Clarify (5 min)**

> "Let me make sure I understand the scope. We're designing an order
> management system for an e-commerce platform. Key questions:
>
> - Is this B2C (customer-facing) or B2B?
> - Do we handle the full lifecycle (create → deliver → return)?
> - Do we need an order promise (delivery ETA)?
> - Multi-region or single region?
> - What's the expected scale? Orders per day?
>
> I'll assume: B2C, full lifecycle, with order promise, single region
> (GCP us-central1), ~100K orders/day. I'll exclude product catalog
> and user auth — assume they exist."

**Phase 2: Estimate (3 min)**

> "Let me do quick math:
> - 100K orders/day ÷ 86400 ≈ 1.2 QPS average, ~12 QPS peak
> - Order reads (tracking): ~10x writes ≈ 120 QPS peak
> - Storage: 2KB per order × 100K/day × 365 = ~73 GB/year
> - At this scale, single PostgreSQL instance handles everything
> - Promise API (product pages): ~1000 QPS (cacheable)"

**Phase 3: High-Level Design (10 min)**

> Draw the architecture (simplified):
> - API Gateway → Order Service → PostgreSQL
> - Order Service → Pub/Sub → Fulfillment, Notification, Analytics
> - Promise Service → Inventory + Carrier APIs → Redis cache
> - Show 3-5 core API endpoints
> - Show core data model (orders, items, events tables)

**Phase 4: Deep Dive (20 min)**

> Pick 2-3 areas to go deep:
> 1. Order state machine + saga pattern for payment/inventory
> 2. Idempotent order creation (idempotency key pattern)
> 3. Order promise calculation + caching strategy
> 4. Inventory reservation (oversell prevention)

**Phase 5: Tradeoffs (5 min)**

> "Tradeoffs I made:
> - Chose eventual consistency for order tracking (read replica) but
>   strong consistency for order creation and inventory
> - Chose choreography saga — at higher complexity, I'd switch to
>   orchestration with a saga coordinator service
> - Chose Cloud SQL over Spanner — at this scale, single-region
>   PostgreSQL is simpler and cheaper. At global scale, I'd consider
>   Spanner for multi-region strong consistency
> - Monitoring: I'd track order creation latency, payment failure
>   rate, promise accuracy, and saga completion rate as key SLIs"
```
