# Deep Dive: Go Map Internals — hmap, bmap, Growth & the Concurrent Fatal

> Everything the runtime does when you read a key, insert a pair, grow
> the bucket array, or iterate — and why concurrent access kills your process.

---

## Table of Contents

1. [The Runtime Struct: `runtime.hmap`](#1-the-runtime-struct-runtimehmap)
2. [Bucket Structure: `runtime.bmap`](#2-bucket-structure-runtimebmap)
3. [Hash Function](#3-hash-function)
4. [Step-by-Step: Map Lookup (`m[key]`)](#4-step-by-step-map-lookup-mkey)
5. [Step-by-Step: Map Assignment (`m[key] = value`)](#5-step-by-step-map-assignment-mkey--value)
6. [Map Growth — Incremental Evacuation](#6-map-growth--incremental-evacuation)
7. [Iteration Randomization](#7-iteration-randomization)
8. [Delete Operation](#8-delete-operation)
9. [Nil Map vs Empty Map](#9-nil-map-vs-empty-map)
10. [Concurrent Access — The Fatal Race](#10-concurrent-access--the-fatal-race)
11. [Key Type Requirements](#11-key-type-requirements)
12. [Performance Characteristics](#12-performance-characteristics)
13. [Quick Reference Card](#13-quick-reference-card)

---

## 1. The Runtime Struct: `runtime.hmap`

A map variable in Go is a **pointer to an `hmap` struct** — exactly like channels.
When you write `m := make(map[string]int)`, `m` is `*hmap`, not `hmap`.
This is why passing a map to a function lets the callee modify it — both hold the same pointer.

**Source:** `runtime/map.go`
```go
type hmap struct {
    count     int            // number of live key-value pairs (len(m) returns this)
    flags     uint8          // iterator/write/grow state flags
    B         uint8          // log₂ of bucket count (buckets = 2^B)
    noverflow uint16         // approximate count of overflow buckets
    hash0     uint32         // hash seed — randomized at map creation
    buckets   unsafe.Pointer // pointer to array of 2^B bmap structs
    oldbuckets unsafe.Pointer // non-nil during growth — points to old bucket array
    nevacuate uintptr        // evacuation progress counter (bucket index)
    extra     *mapextra      // overflow bucket pre-allocation & bookkeeping
}
```

### Memory Layout

```
  m := make(map[string]int)

  Stack variable 'm':
  ┌──────────────────┐
  │ *hmap = 0xc00008 │──────────────────────────────────────────────┐
  └──────────────────┘                                              │
                                                                    ▼
  Heap: hmap struct                                                 
  ┌──────────────────────────────────────────────────────────────────┐
  │ count      = 0          // no entries yet                       │
  │ flags      = 0x00       // no active iterators, no writes       │
  │ B          = 0          // 2^0 = 1 bucket (minimum)             │
  │ noverflow  = 0          // no overflow buckets                  │
  │ hash0      = 0xa3f7c291 // random seed (unique per map)         │
  │ buckets   ─┼──► [ bmap_0 ]   (array of 2^B = 1 bucket)         │
  │ oldbuckets = nil        // not growing                          │
  │ nevacuate  = 0          // nothing to evacuate                  │
  │ extra      = nil        // no overflow bookkeeping yet          │
  └──────────────────────────────────────────────────────────────────┘
```

### Key Fields Explained

```
┌──────────────┬──────────────────────────────────────────────────────────┐
│ Field        │ Purpose                                                  │
├──────────────┼──────────────────────────────────────────────────────────┤
│ count        │ Live entries. len(m) reads this directly — O(1)          │
│ flags        │ Bit flags: hashWriting (0x04), iterator (0x01),          │
│              │ oldIterator (0x02), sameSizeGrow (0x08)                  │
│ B            │ log₂(bucket count). B=3 means 8 buckets. B=0 = 1 bucket │
│ noverflow    │ Approximate overflow bucket count — triggers same-size   │
│              │ grow when too many overflows accumulate                  │
│ hash0        │ Random seed generated at make(). Prevents hash flooding  │
│              │ attacks. Also the reason iteration order is randomized   │
│ buckets      │ Current bucket array — 2^B contiguous bmap structs       │
│ oldbuckets   │ Previous bucket array during incremental growth          │
│ nevacuate    │ Index of next old bucket to evacuate — monotonically     │
│              │ increasing until all old buckets are moved               │
│ extra        │ mapextra: preallocated overflow buckets, pointers for GC │
└──────────────┴──────────────────────────────────────────────────────────┘
```

**Why a pointer?** If `m` were a value type (like a slice header), passing it to a function
would copy the header. Mutations inside the function (inserting keys) would allocate new
buckets invisible to the caller. The Go team made maps pointer types so mutation is always
shared. This also means the zero value is `nil`, not an empty map — Section 9 covers why
that matters.

---

## 2. Bucket Structure: `runtime.bmap`

Each bucket holds exactly **8 key-value pairs**. This is the constant `bucketCnt = 8`
in `runtime/map.go`. The bucket layout is carefully designed for performance.

**Source:** `runtime/map.go`
```go
// Compile-time struct (only the tophash is declared in source):
type bmap struct {
    tophash [bucketCnt]uint8  // top 8 bits of hash for each slot
    // Followed by keys, values, and overflow pointer — laid out by compiler
}
```

The compiler generates the actual layout. Keys and values are stored in **separate arrays**,
not interleaved as `key0, val0, key1, val1, ...`. This avoids padding waste.

### Detailed Memory Layout

```
  Single bucket: bmap for map[string]int (key=string=16B, value=int=8B)

  Offset  Contents
  ──────  ────────────────────────────────────────────────────────────
  0x00    ┌─────────────────────────────────────────────────────────┐
          │ tophash[0] │ tophash[1] │ ... │ tophash[7]             │
          │  (1 byte)  │  (1 byte)  │     │  (1 byte)  = 8 bytes  │
  0x08    ├─────────────────────────────────────────────────────────┤
          │ key[0]  (string = 16 bytes: ptr + len)                 │
  0x18    │ key[1]  (16 bytes)                                     │
  0x28    │ key[2]  (16 bytes)                                     │
  ...     │ ...                                                    │
  0x78    │ key[7]  (16 bytes)                                     │
  0x88    ├─────────────────────────────────────────────────────────┤
          │ val[0]  (int = 8 bytes)                                │
  0x90    │ val[1]  (8 bytes)                                      │
  ...     │ ...                                                    │
  0xC0    │ val[7]  (8 bytes)                                      │
  0xC8    ├─────────────────────────────────────────────────────────┤
          │ overflow *bmap  (8 bytes — pointer to next bucket)     │
  0xD0    └─────────────────────────────────────────────────────────┘

  Total bucket size = 8 + (8×16) + (8×8) + 8 = 208 bytes
```

### Why Keys and Values Are Separated

```
  INTERLEAVED (how you might expect — Go does NOT do this):
  ┌──────────────────┬────────────┬─────────┬──────────────────┬────────────┐
  │ key[0] (string)  │ pad (8B)   │val[0]   │ key[1] (string)  │ pad (8B)  │...
  │ 16 bytes         │ WASTED     │ 8 bytes │ 16 bytes         │ WASTED    │
  └──────────────────┴────────────┴─────────┴──────────────────┴───────────┘
  For map[int8]int64: each pair would need 7 bytes of padding between key and value.

  SEPARATED (what Go actually does):
  ┌──────────────────────────────────┬──────────────────────────┐
  │ key[0]│key[1]│...│key[7]         │ val[0]│val[1]│...│val[7] │
  │ all 8 keys packed together       │ all 8 values packed      │
  └──────────────────────────────────┴──────────────────────────┘
  No padding waste. map[int8]int64 saves 7 × 8 = 56 bytes per bucket.
```

### Tophash — The Fast-Reject Filter

Each `tophash[i]` stores the **top 8 bits** of the hash for the key in slot `i`.
During lookup, the runtime compares `tophash` bytes before doing a full key comparison.
Since tophash fits in a single cache line (8 bytes), this rejects 255/256 (~99.6%) of
non-matching slots without ever touching the actual keys.

```
  Special tophash values (runtime/map.go):
  ┌───────────────┬────────┬──────────────────────────────────────────────┐
  │ Name          │ Value  │ Meaning                                      │
  ├───────────────┼────────┼──────────────────────────────────────────────┤
  │ emptyRest     │ 0      │ This cell AND all higher indexes are empty   │
  │ emptyOne      │ 1      │ This cell is empty (but higher ones may not) │
  │ evacuatedX    │ 2      │ Key was evacuated to first half of new array │
  │ evacuatedY    │ 3      │ Key was evacuated to second half             │
  │ evacuatedEmpty│ 4      │ Cell is empty AND bucket has been evacuated  │
  │ minTopHash    │ 5      │ Minimum valid tophash value for a real key   │
  └───────────────┴────────┴──────────────────────────────────────────────┘

  If hash top bits compute to 0-4, the runtime adds minTopHash (5) to
  disambiguate from the special sentinel values.
```

---

## 3. Hash Function

Go uses architecture-specific hash functions for performance. On `amd64` with AES-NI,
it uses **AES-based hashing** (`runtime.aeshash`). On other architectures, it falls back
to `runtime.memhash` (a variant of wyhash/xxhash-style algorithms).

**Source:** `runtime/hash64.go`, `runtime/asm_amd64.s`

### Hash Seed: `hash0`

Every map instance gets a **unique random seed** (`hash0`) assigned at creation time.
This has two consequences:

1. **Security:** prevents hash-flooding DoS attacks. An attacker cannot predict which
   bucket a key will land in, because the seed is different for every map.
2. **Randomized iteration:** since bucket assignments change with the seed, the physical
   ordering of keys varies per map instance.

```go
// runtime/map.go — simplified from makemap()
func makemap(t *maptype, hint int, h *hmap) *hmap {
    h.hash0 = uint32(rand())  // random seed per map
    // ...
}
```

### The Full Lookup Path: From Key to Slot

```
  Input: key = "user:42", hash0 = 0xa3f7c291, B = 3 (8 buckets)

  Step 1: HASH
  ┌──────────────────────────────────────────────┐
  │ hash = aeshash("user:42", hash0)             │
  │      = 0xb7e1_5a3f_928d_4c01   (64-bit hash) │
  └──────────────────────────────────────────────┘

  Step 2: BUCKET INDEX (low B bits)
  ┌──────────────────────────────────────────────┐
  │ bucket = hash & (2^B - 1)                    │
  │        = 0x...4c01 & 0x07                    │
  │        = 1  (bucket index)                   │
  └──────────────────────────────────────────────┘
           Low bits select the bucket ───────┘

  Step 3: TOPHASH (top 8 bits of hash)
  ┌──────────────────────────────────────────────┐
  │ top = uint8(hash >> (64 - 8))                │
  │     = uint8(0xb7)                            │
  │     = 183                                    │
  └──────────────────────────────────────────────┘
           Top bits form the fast-compare tag ─┘

  Step 4: SCAN bucket 1's tophash array
  ┌─────────────────────────────────────────────────────────┐
  │ bucket[1].tophash = [0x3a, 0xb7, 0x00, 0x00, ...]      │
  │                           ^^^^                          │
  │                     match at slot 1! → compare full key │
  └─────────────────────────────────────────────────────────┘

  Step 5: FULL KEY COMPARE
  ┌──────────────────────────────────────────────┐
  │ bucket[1].keys[1] == "user:42" ?             │
  │ → YES → return bucket[1].values[1]           │
  └──────────────────────────────────────────────┘
```

**Why low bits for bucket, high bits for tophash?** During growth, the bucket count
doubles (B increases by 1). The new bit that becomes significant is the next-lowest bit.
A key either stays in bucket `i` or moves to bucket `i + 2^(B-1)`. The tophash (high bits)
remains valid — no need to recompute it. This is an elegant design enabling incremental
evacuation (Section 6).

---

## 4. Step-by-Step: Map Lookup (`m[key]`)

**Source:** `runtime/map.go` → `mapaccess1()` (single return), `mapaccess2()` (comma-ok)

### Trace: `v := m["user:42"]`

```
  mapaccess1(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer

  ┌─────────────────────────────────────────────────────────────────┐
  │ 1. CHECK: is h == nil or h.count == 0?                         │
  │    → YES: return pointer to zero value (safe, not a panic)     │
  │    → NO: continue                                              │
  ├─────────────────────────────────────────────────────────────────┤
  │ 2. CHECK: flags & hashWriting != 0?                            │
  │    → YES: fatal("concurrent map read and map write") — CRASH   │
  │    → NO: continue                                              │
  ├─────────────────────────────────────────────────────────────────┤
  │ 3. HASH: hash = t.hasher(key, uintptr(h.hash0))               │
  │    Compute 64-bit hash of "user:42" with this map's seed       │
  ├─────────────────────────────────────────────────────────────────┤
  │ 4. BUCKET: m = 1<<h.B - 1  (bucket mask)                      │
  │    bucket = hash & m  → index into h.buckets array             │
  ├─────────────────────────────────────────────────────────────────┤
  │ 5. GROWING? if h.oldbuckets != nil                             │
  │    → old bucket may not be evacuated yet                       │
  │    → check if oldbuckets[hash & old_mask] is evacuated         │
  │    → if NOT evacuated, search there instead                    │
  ├─────────────────────────────────────────────────────────────────┤
  │ 6. TOPHASH: top = tophash(hash)                                │
  │    Extract top 8 bits, adjust if < minTopHash (5)              │
  ├─────────────────────────────────────────────────────────────────┤
  │ 7. SCAN: for each bucket in overflow chain:                    │
  │      for i := 0; i < 8; i++ {                                  │
  │        if b.tophash[i] != top {                                │
  │          if b.tophash[i] == emptyRest { break }  // done early │
  │          continue                                              │
  │        }                                                       │
  │        k := bucket.key(i)  // pointer arithmetic to key slot   │
  │        if t.key.equal(key, k) {                                │
  │          return bucket.value(i)  // FOUND — return value ptr   │
  │        }                                                       │
  │      }                                                         │
  │      b = b.overflow  // follow overflow chain                  │
  ├─────────────────────────────────────────────────────────────────┤
  │ 8. MISS: return pointer to zero value                          │
  │    (a static zero-sized allocation, safe to read)              │
  └─────────────────────────────────────────────────────────────────┘
```

### The Comma-Ok Idiom

```go
v, ok := m["user:42"]   // calls mapaccess2() instead of mapaccess1()
```

`mapaccess2()` is identical to `mapaccess1()` but returns a second `bool`:

```
  mapaccess2() returns:
  ┌───────┬────────────────────────────────────┐
  │ Found │ return (pointer to value, true)    │
  │ Miss  │ return (pointer to zero, false)    │
  └───────┴────────────────────────────────────┘
```

This is how you distinguish "key exists with zero value" from "key does not exist":

```go
count, ok := m["page_views"]
// ok == true,  count == 0  → key exists, value is genuinely 0
// ok == false, count == 0  → key does not exist, count is zero value
```

### Fast-Path Variants

The runtime also has type-specialized lookup functions to avoid the overhead of
generic `mapaccess1()`:

```
  runtime/map_fast32.go  → mapaccess1_fast32()   // 4-byte keys (int32, uint32)
  runtime/map_fast64.go  → mapaccess1_fast64()   // 8-byte keys (int64, uint64, pointers)
  runtime/map_faststr.go → mapaccess1_faststr()  // string keys
```

These avoid indirect calls through `t.hasher` and `t.key.equal`, using inline hash
computation and comparison. The compiler selects the appropriate variant at build time.

---

## 5. Step-by-Step: Map Assignment (`m[key] = value`)

**Source:** `runtime/map.go` → `mapassign()`

### Trace: `m["user:42"] = 100`

```
  mapassign(t *maptype, h *hmap, key unsafe.Pointer) unsafe.Pointer

  ┌─────────────────────────────────────────────────────────────────┐
  │ 1. CHECK: is h == nil?                                         │
  │    → YES: panic("assignment to entry in nil map")              │
  │    → NO: continue                                              │
  ├─────────────────────────────────────────────────────────────────┤
  │ 2. CHECK: flags & hashWriting != 0?                            │
  │    → YES: fatal("concurrent map writes") — HARD CRASH          │
  │    → NO: set flags |= hashWriting (claim write ownership)     │
  ├─────────────────────────────────────────────────────────────────┤
  │ 3. HASH: hash = t.hasher(key, uintptr(h.hash0))               │
  ├─────────────────────────────────────────────────────────────────┤
  │ 4. BUCKET: bucket = hash & (1<<h.B - 1)                       │
  │    Get pointer to target bucket in h.buckets                   │
  ├─────────────────────────────────────────────────────────────────┤
  │ 5. GROWING? if h.oldbuckets != nil                             │
  │    → evacuate one old bucket (incremental migration)           │
  │    → this is how growth is amortized over inserts/deletes      │
  ├─────────────────────────────────────────────────────────────────┤
  │ 6. SCAN for existing key OR first empty slot:                  │
  │                                                                │
  │    var inserti *uint8       // tophash slot for insertion       │
  │    var insertk, insertv    // key/value slot pointers           │
  │                                                                │
  │    for each bucket in overflow chain:                           │
  │      for i := 0; i < 8; i++ {                                  │
  │        if b.tophash[i] != top {                                │
  │          if isEmpty(b.tophash[i]) && inserti == nil {          │
  │            inserti = &b.tophash[i]  // remember first empty    │
  │            insertk = b.key(i)                                  │
  │            insertv = b.value(i)                                │
  │          }                                                     │
  │          if b.tophash[i] == emptyRest { break }                │
  │          continue                                              │
  │        }                                                       │
  │        k := b.key(i)                                           │
  │        if !t.key.equal(key, k) { continue }                   │
  │        // KEY EXISTS → update value in place                   │
  │        b.value(i) = value                                      │
  │        goto done                                               │
  │      }                                                         │
  │      b = b.overflow                                            │
  ├─────────────────────────────────────────────────────────────────┤
  │ 7. GROW CHECK: before inserting a new key, check triggers:     │
  │    a) load factor > 6.5  (count / 2^B > 6.5)                  │
  │       → double growth: allocate 2^(B+1) buckets               │
  │    b) too many overflow buckets (noverflow >= 2^min(B, 15))    │
  │       → same-size growth: repack to eliminate overflow chains  │
  │    If triggered: call hashGrow(), restart scan in new buckets  │
  ├─────────────────────────────────────────────────────────────────┤
  │ 8. INSERT: if no empty slot found in existing chain            │
  │    → allocate overflow bucket, link to chain                   │
  │    → use first slot of new overflow bucket                     │
  │                                                                │
  │    Write tophash[inserti] = top                                │
  │    Copy key into insertk slot                                  │
  │    Return insertv pointer (caller writes value through it)     │
  │    h.count++                                                   │
  ├─────────────────────────────────────────────────────────────────┤
  │ 9. DONE: clear flags &^= hashWriting                          │
  └─────────────────────────────────────────────────────────────────┘
```

### Load Factor: 6.5

The magic number 6.5 means the map targets an average of **6.5 entries per bucket**.
Since each bucket holds 8 slots, this means ~81% occupancy before triggering growth.
This was chosen empirically by the Go team as the sweet spot between:

```
  Too low (e.g., 4.0):  wastes memory — many empty slots
  Too high (e.g., 7.5): too many overflow chains — lookup becomes O(n)
  6.5: ~81% full — one overflow chain per ~3 buckets on average
```

### Overflow Bucket Chain

```
  When all 8 slots in a bucket are full and a new key hashes to the same bucket:

  buckets[3]:
  ┌─────────────────────┐     ┌─────────────────────┐
  │ tophash[0..7]       │     │ tophash[0..7]       │
  │ keys[0..7]          │     │ keys[0..7]          │  overflow bucket
  │ values[0..7]        │     │ values[0..7]        │
  │ overflow ───────────┼────►│ overflow ───────────┼──► nil
  └─────────────────────┘     └─────────────────────┘

  Each overflow bucket is another full bmap with 8 more slots.
  Long chains degrade performance — this triggers same-size growth.
```

---

## 6. Map Growth — Incremental Evacuation

**Source:** `runtime/map.go` → `hashGrow()`, `growWork()`, `evacuate()`

Go maps **never rehash everything at once**. A map with 1 million entries would cause a
multi-millisecond pause if rehashing were synchronous. Instead, Go amortizes the cost:
each subsequent `mapassign()` or `mapdelete()` evacuates 1-2 old buckets.

### Two Types of Growth

```
┌──────────────────────┬──────────────────────────────────────────────────┐
│ Trigger              │ Action                                          │
├──────────────────────┼──────────────────────────────────────────────────┤
│ count/2^B > 6.5      │ DOUBLE GROWTH: allocate 2^(B+1) new buckets.   │
│ (too full)           │ B increments by 1. Each old bucket splits into  │
│                      │ two new buckets (X = same index, Y = index +   │
│                      │ 2^oldB). Reduces load factor by half.           │
├──────────────────────┼──────────────────────────────────────────────────┤
│ noverflow too high   │ SAME-SIZE GROWTH: allocate 2^B new buckets     │
│ (too many overflow   │ (same count). B stays the same. Repacks entries │
│ chains, sparse data) │ to eliminate overflow chains caused by deletes. │
│                      │ Compacts the data without increasing capacity.  │
└──────────────────────┴──────────────────────────────────────────────────┘
```

### The Growth Sequence

```
  hashGrow() — called from mapassign() when trigger conditions are met

  ┌─────────────────────────────────────────────────────────────────┐
  │ 1. Allocate new bucket array: makeBucketArray(B+1 or B)        │
  │ 2. h.oldbuckets = h.buckets   // save old array               │
  │ 3. h.buckets = newBuckets     // install new array             │
  │ 4. h.nevacuate = 0            // start evacuating from index 0 │
  │ 5. h.noverflow = 0            // reset overflow counter        │
  │                                                                │
  │ NOTE: no entries are moved yet! The old data is still in       │
  │ oldbuckets. Evacuation happens incrementally.                  │
  └─────────────────────────────────────────────────────────────────┘
```

### Incremental Evacuation: `evacuate()`

On every `mapassign()` or `mapdelete()` call during growth, the runtime calls
`growWork()` which evacuates the bucket at index `nevacuate` (and sometimes one more):

```
  evacuate(t *maptype, h *hmap, oldbucket uintptr)

  For DOUBLE growth (B went from 3 to 4):

  Old bucket array (2^3 = 8 buckets)     New bucket array (2^4 = 16 buckets)
  ┌───────┐                               ┌───────┐
  │ old[0]│─ ─ ─ ─evacuate─ ─ ─ ─ ─ ─ ─►│ new[0]│  X destination
  │       │─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─ ─►│ new[8]│  Y destination
  ├───────┤                               ├───────┤
  │ old[1]│                               │ new[1]│
  │       │                               │ new[9]│
  ├───────┤                               ├───────┤
  │ old[2]│  ◄── nevacuate = 2            │ new[2]│
  │       │      (next to evacuate)       │new[10]│
  ├───────┤                               ├───────┤
  │  ...  │                               │  ...  │
  └───────┘                               └───────┘

  For each key in old[i]:
    new_hash_bit = hash & (1 << oldB)     // the NEW significant bit
    if new_hash_bit == 0:
      → goes to new[i]           (X: same index)
    else:
      → goes to new[i + 2^oldB]  (Y: index + old bucket count)
```

### Dual Lookup During Evacuation

While `oldbuckets != nil`, every lookup must check whether the target bucket has been
evacuated yet:

```
  mapaccess during growth:

  ┌─────────────────────────────────────────────────────────────┐
  │ bucket_idx = hash & new_mask                                │
  │ old_bucket_idx = hash & old_mask                            │
  │                                                             │
  │ if oldbuckets != nil:                                       │
  │   old_b = oldbuckets[old_bucket_idx]                        │
  │   if NOT evacuated(old_b):   // check tophash for markers  │
  │     search in old_b instead  // data hasn't moved yet       │
  │   else:                                                     │
  │     search in buckets[bucket_idx]  // already evacuated     │
  └─────────────────────────────────────────────────────────────┘
```

### Evacuation Progress

```
  Timeline of a growth event (B: 2→3, 4 old buckets → 8 new):

  Insert #1: evacuate old[0] → new[0] + new[4]    nevacuate = 1
  Insert #2: evacuate old[1] → new[1] + new[5]    nevacuate = 2
  Insert #3: evacuate old[2] → new[2] + new[6]    nevacuate = 3
  Insert #4: evacuate old[3] → new[3] + new[7]    nevacuate = 4
             ▲
             └─ nevacuate == 2^oldB → growth complete!
                h.oldbuckets = nil  (old array eligible for GC)
```

**Production implication:** during growth, there is a ~2x memory spike (both old and new
arrays coexist). For maps with millions of entries, this can cause OOM in memory-constrained
containers. Pre-allocate with `make(map[K]V, hint)` to avoid unexpected growth.

---

## 7. Iteration Randomization

**Source:** `runtime/map.go` → `mapiterinit()`, `mapiternext()`

Go **deliberately** randomizes map iteration order. This is not an accident or
implementation detail — it is a design decision enforced since Go 1.0 to prevent
code from depending on insertion order.

### How It Works

```go
// runtime/map.go — simplified
func mapiterinit(t *maptype, h *hmap, it *hiter) {
    // Pick a random starting bucket
    r := uintptr(rand())
    it.startBucket = r & bucketMask(h.B)
    // Pick a random starting offset WITHIN that bucket
    it.offset = uint8(r >> h.B & (bucketCnt - 1))
    it.bucket = it.startBucket
    // ...
}
```

### The Iterator Struct: `hiter`

```
  runtime.hiter
  ┌──────────────────────────────────────────────────────────────┐
  │ key         unsafe.Pointer  // current key (for range loop)  │
  │ elem        unsafe.Pointer  // current value                 │
  │ t           *maptype        // map type descriptor            │
  │ h           *hmap           // the map being iterated         │
  │ buckets     unsafe.Pointer  // bucket array at iter start     │
  │ bptr        *bmap           // current bucket pointer         │
  │ startBucket uintptr         // random starting bucket index   │
  │ offset      uint8           // random starting offset (0-7)   │
  │ wrapped     bool            // has the iterator wrapped around│
  │ B           uint8           // B at iteration start           │
  │ i           uint8           // current index within bucket    │
  │ bucket      uintptr         // current bucket index           │
  │ checkBucket uintptr         // for growth-during-iteration    │
  └──────────────────────────────────────────────────────────────┘
```

### Iteration Walk Order

```
  B = 2 (4 buckets), random startBucket = 2, random offset = 5

  Walk order (bucket:slot):
  2:5 → 2:6 → 2:7 → 2:0 → 2:1 → 2:2 → 2:3 → 2:4  (bucket 2, then overflow)
  3:5 → 3:6 → 3:7 → 3:0 → 3:1 → 3:2 → 3:3 → 3:4  (bucket 3)
  0:5 → 0:6 → 0:7 → 0:0 → 0:1 → 0:2 → 0:3 → 0:4  (bucket 0 — wrapped)
  1:5 → 1:6 → 1:7 → 1:0 → 1:1 → 1:2 → 1:3 → 1:4  (bucket 1)
  └─ starts at offset 5 within each bucket, wraps around to 0-4

  Every slot across all buckets (including overflow chains) is visited exactly once.
```

### Iteration During Growth

If the map is growing while being iterated, the iterator must handle the case where
some old buckets are evacuated and some are not:

```
  Iterator snapshot: B_at_start, buckets_at_start
  For each bucket to visit:
    if oldbuckets != nil AND old bucket not evacuated:
      → iterate from old bucket, but SKIP keys that would
        belong to the other half after evacuation
    else:
      → iterate from new bucket normally
```

This guarantees that each key is returned **exactly once**, even during growth.

---

## 8. Delete Operation

**Source:** `runtime/map.go` → `mapdelete()`

### Trace: `delete(m, "user:42")`

```
  mapdelete(t *maptype, h *hmap, key unsafe.Pointer)

  ┌─────────────────────────────────────────────────────────────────┐
  │ 1. CHECK: h == nil or h.count == 0 → return (no-op)           │
  │ 2. CHECK: flags & hashWriting → fatal("concurrent map writes")│
  │ 3. Set flags |= hashWriting                                   │
  │ 4. HASH: hash = t.hasher(key, uintptr(h.hash0))               │
  │ 5. BUCKET: bucket = hash & (1<<h.B - 1)                       │
  │ 6. GROWING? → evacuate one old bucket (incremental)            │
  │ 7. SCAN: find key in bucket chain (same as lookup)             │
  │    Found?                                                      │
  │    ├─ Clear the key memory (memclr for GC — release references)│
  │    ├─ Clear the value memory (memclr for GC)                   │
  │    ├─ Set tophash[i] = emptyOne                                │
  │    ├─ Compact: if tophash[i+1] == emptyRest, upgrade emptyOne  │
  │    │           entries backward to emptyRest for faster scans  │
  │    ├─ h.count--                                                │
  │    └─ Break                                                    │
  │    Not found? → no-op (deleting non-existent key is safe)      │
  │ 8. Clear flags &^= hashWriting                                 │
  └─────────────────────────────────────────────────────────────────┘
```

### Maps Never Shrink

```
  ┌────────────────────────────────────────────────────────────────────┐
  │ CRITICAL: delete() NEVER deallocates buckets.                     │
  │                                                                    │
  │ If you insert 1 million keys and then delete 999,999 of them,     │
  │ the map still holds memory for all the buckets that were allocated │
  │ during growth. The bucket array is NOT resized downward.           │
  │                                                                    │
  │ The buckets are only freed when the entire map is garbage          │
  │ collected (no more references to it).                              │
  │                                                                    │
  │ Workaround: create a new map and copy surviving keys:             │
  │   newMap := make(map[K]V, len(oldMap))                            │
  │   for k, v := range oldMap { newMap[k] = v }                     │
  │   oldMap = newMap  // old map becomes eligible for GC             │
  └────────────────────────────────────────────────────────────────────┘
```

**Why?** Shrinking would require evacuation in reverse — moving entries from many buckets
into fewer buckets. The Go team decided the complexity and latency risk was not worth it.
Most maps in production either grow monotonically or get replaced entirely.

---

## 9. Nil Map vs Empty Map

```go
var m1 map[string]int        // nil map — no hmap allocated
m2 := map[string]int{}       // empty map — hmap exists, count=0
m3 := make(map[string]int)   // empty map — same as m2
```

### Behavior Comparison

```
  ┌───────────────────────┬──────────────────────┬──────────────────────┐
  │ Operation             │ nil map              │ empty map            │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ v := m["key"]         │ returns zero value   │ returns zero value   │
  │                       │ (safe, no panic)     │ (safe, no panic)     │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ v, ok := m["key"]     │ ok = false (safe)    │ ok = false (safe)    │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ len(m)                │ 0 (safe)             │ 0 (safe)             │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ for k, v := range m   │ 0 iterations (safe)  │ 0 iterations (safe)  │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ delete(m, "key")      │ no-op (safe)         │ no-op (safe)         │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ m["key"] = 1          │ PANIC ❌             │ works ✅             │
  │                       │ "assignment to entry │                      │
  │                       │  in nil map"         │                      │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ json.Marshal(m)       │ "null"               │ "{}"                 │
  ├───────────────────────┼──────────────────────┼──────────────────────┤
  │ m == nil              │ true                 │ false                │
  └───────────────────────┴──────────────────────┴──────────────────────┘
```

### Why Write Panics on Nil Map

```
  var m map[string]int   // m is a nil pointer — *hmap == nil

  Stack:
  ┌──────────────────┐
  │ m = nil          │   No hmap struct exists on the heap.
  └──────────────────┘   There is no bucket array to insert into.

  m["key"] = 1
  → runtime calls mapassign(t, h, key)
  → h == nil
  → panic("assignment to entry in nil map")

  Reads are safe because mapaccess1() checks h == nil first and returns
  a pointer to a static zero value without touching any fields.
```

### The JSON Gotcha

```go
type Response struct {
    Items map[string]int `json:"items"`
}

r1 := Response{}                          // Items is nil
r2 := Response{Items: map[string]int{}}   // Items is empty

json.Marshal(r1)  // {"items":null}    ← APIs often reject null
json.Marshal(r2)  // {"items":{}}      ← clean empty object
```

**Production rule:** initialize maps in constructors to avoid nil/null serialization issues.

---

## 10. Concurrent Access — The Fatal Race

Maps are **NOT safe for concurrent use**. The runtime does not use any locking internally.
Instead, it detects concurrent access and **kills the process** — not with a panic (which
can be recovered), but with `fatal()` which is unrecoverable.

**Source:** `runtime/map.go` — the `hashWriting` flag

### How Detection Works

```
  mapassign():                          mapaccess1() (concurrent):
  ┌─────────────────────────┐           ┌─────────────────────────┐
  │ flags |= hashWriting    │           │ if flags & hashWriting  │
  │ ... do the insert ...   │           │   → fatal() ← CRASH    │
  │ flags &^= hashWriting   │           │                         │
  └─────────────────────────┘           └─────────────────────────┘

  The runtime sets a flag before writing and checks it during reads.
  This is NOT a mutex — it's a best-effort detection mechanism.
  It catches most concurrent access but is not guaranteed to catch all races.
```

### The Fatal Message

```
fatal error: concurrent map read and map write
fatal error: concurrent map writes

goroutine 42 [running]:
runtime.throw({0x4a7b20, 0x26})
    runtime/panic.go:1077 +0x48
runtime.mapaccess1_faststr(0xc000080060, {0x4a3e10, 0x5})
    runtime/map_faststr.go:12 +0x1a0

This is NOT a panic — recover() cannot catch it.
The process exits immediately. No cleanup, no deferred functions.
```

### Solutions

```
  ┌──────────────────────┬───────────────────────────────────────────────┐
  │ Pattern              │ When to use                                   │
  ├──────────────────────┼───────────────────────────────────────────────┤
  │ sync.Mutex           │ Simple cases. Lock for both read and write.   │
  │                      │ Best when writes are frequent.                │
  │                      │                                               │
  │ var mu sync.Mutex    │   mu.Lock()                                   │
  │ m := make(map[K]V)  │   m[key] = val                                │
  │                      │   mu.Unlock()                                 │
  ├──────────────────────┼───────────────────────────────────────────────┤
  │ sync.RWMutex         │ Read-heavy workloads. Multiple readers can    │
  │                      │ proceed concurrently, writers get exclusive   │
  │                      │ access. Best when reads >> writes.            │
  │                      │                                               │
  │ var mu sync.RWMutex  │   mu.RLock()    // read path                  │
  │ m := make(map[K]V)  │   v := m[key]                                 │
  │                      │   mu.RUnlock()                                │
  │                      │                                               │
  │                      │   mu.Lock()     // write path                  │
  │                      │   m[key] = val                                │
  │                      │   mu.Unlock()                                 │
  ├──────────────────────┼───────────────────────────────────────────────┤
  │ sync.Map             │ Two specific patterns:                        │
  │                      │ (1) write-once, read-many (config caches)     │
  │                      │ (2) disjoint key sets per goroutine           │
  │                      │ For everything else, map + RWMutex is faster. │
  └──────────────────────┴───────────────────────────────────────────────┘
```

### sync.Map Internals

`sync.Map` uses a **two-map architecture** to minimize lock contention for reads:

```
  sync.Map
  ┌───────────────────────────────────────────────────────────────────┐
  │ read  atomic.Pointer[readOnly]  // lock-free read path           │
  │   └─ readOnly {                                                  │
  │        m     map[any]*entry     // immutable snapshot            │
  │        amended bool             // true if dirty has new keys    │
  │      }                                                           │
  │                                                                   │
  │ mu    sync.Mutex                // protects dirty                 │
  │ dirty map[any]*entry            // mutable map, nil when clean   │
  │ misses int                      // read misses since last promote│
  └───────────────────────────────────────────────────────────────────┘

  Load("key"):
    1. Try read.m["key"] → atomic load, no lock → FAST PATH
    2. Miss + amended? → mu.Lock(), try dirty["key"], misses++
    3. misses >= len(dirty)? → PROMOTE: read.m = dirty, dirty = nil
       (amortized O(1) — promotion copies the map pointer, not entries)

  Store("key", val):
    1. If key in read.m → atomic CAS the *entry pointer → no lock!
    2. Else → mu.Lock(), add to dirty (and read if needed)
```

**Key insight:** `sync.Map` shines when most operations are reads to existing keys.
The read path is entirely lock-free (atomic pointer load). But if your workload involves
frequent writes to new keys, every Store takes the mutex, and `sync.Map` is slower than
a simple `map + sync.RWMutex`.

---

## 11. Key Type Requirements

Map keys must satisfy the `comparable` constraint — they must support `==` and `!=`.

### What Can Be a Key

```
  ┌─────────────────────┬────────┬─────────────────────────────────────┐
  │ Type                │ Key?   │ Notes                               │
  ├─────────────────────┼────────┼─────────────────────────────────────┤
  │ bool                │ ✅     │                                     │
  │ int, uint, float*   │ ✅     │ NaN != NaN — entries with NaN keys  │
  │                     │        │ are unreachable after insert!        │
  │ complex64/128       │ ✅     │ Same NaN caveat                     │
  │ string              │ ✅     │ Most common key type                │
  │ pointer (*T)        │ ✅     │ Compares ADDRESS, not pointed value │
  │ channel             │ ✅     │ Compares channel identity            │
  │ interface           │ ✅     │ Compares dynamic type + value       │
  │                     │        │ Panics at runtime if underlying     │
  │                     │        │ type is not comparable               │
  │ array [N]T          │ ✅     │ Element-wise comparison. Fixed size │
  │ struct              │ ✅*    │ Only if ALL fields are comparable   │
  ├─────────────────────┼────────┼─────────────────────────────────────┤
  │ slice []T           │ ❌     │ Not comparable (reference type)     │
  │ map[K]V             │ ❌     │ Not comparable                      │
  │ function            │ ❌     │ Not comparable                      │
  │ struct with slice/  │ ❌     │ Struct is non-comparable if any     │
  │ map/func field      │        │ field is non-comparable              │
  └─────────────────────┴────────┴─────────────────────────────────────┘
```

### The NaN Trap

```go
m := map[float64]string{}
m[math.NaN()] = "hello"
m[math.NaN()] = "world"

len(m)              // 2 — two separate entries!
v := m[math.NaN()]  // "" — cannot retrieve either one!

// NaN != NaN by IEEE 754. The key can never match on lookup.
// These entries are "zombie" — unreachable but counted.
// They'll only be freed when the map is GC'd.
```

### Pointer Keys vs Value Keys

```go
type Point struct{ X, Y int }

// POINTER keys: compare memory addresses
m1 := map[*Point]string{}
p := &Point{1, 2}
m1[p] = "A"
m1[&Point{1, 2}] = "B"   // different pointer! Two entries.
len(m1) // 2

// STRUCT keys: compare field values
m2 := map[Point]string{}
m2[Point{1, 2}] = "A"
m2[Point{1, 2}] = "B"    // same value! Overwrites.
len(m2) // 1
```

### Array Keys (underused pattern)

```go
// Arrays are comparable — great for composite keys
type CacheKey = [3]string  // [method, path, version]

cache := map[CacheKey]Response{}
cache[CacheKey{"GET", "/api/users", "v2"}] = resp

// Avoids string concatenation for composite keys.
// Fixed size is checked at compile time.
```

---

## 12. Performance Characteristics

### Cost Table

```
┌──────────────────────┬───────────────┬──────────────────────────────────────┐
│ Operation            │ Time          │ Details                              │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ Lookup (hit)         │ O(1) avg      │ Hash + tophash scan (8 compares) +   │
│ m[key]               │ ~50-150ns     │ 1 full key compare. Overflow chains  │
│                      │               │ add ~50ns per extra bucket scanned.  │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ Lookup (miss)        │ O(1) avg      │ Same path, but must scan all 8      │
│                      │ ~50-120ns     │ tophash entries + overflow chain     │
│                      │               │ before confirming miss.              │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ Insert (no growth)   │ O(1) avg      │ Hash + scan + write key/value.       │
│ m[key] = val         │ ~100-250ns    │ Includes tophash update.             │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ Insert (triggers     │ O(1) amort.   │ hashGrow() allocates new array +     │
│ growth)              │ ~500ns-5μs    │ evacuates 1-2 old buckets per call.  │
│                      │               │ One-time allocation spike.           │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ Delete               │ O(1) avg      │ Same as lookup + memclr key/value.   │
│ delete(m, key)       │ ~100-200ns    │ No memory is freed.                  │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ Iteration            │ O(n)          │ Visits every bucket + overflow.      │
│ for k, v := range m  │               │ Cache-unfriendly for sparse maps.    │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ len(m)               │ O(1)          │ Reads h.count directly.              │
├──────────────────────┼───────────────┼──────────────────────────────────────┤
│ make(map, hint)      │ O(hint)       │ Allocates ⌈hint/6.5⌉ buckets        │
│                      │               │ rounded up to next power of 2.       │
└──────────────────────┴───────────────┴──────────────────────────────────────┘
```

### Memory Overhead Per Entry

```
  Bucket overhead (per 8 entries):
    8 bytes tophash + 8 byte overflow pointer = 16 bytes overhead
    Per entry: 16/8 = 2 bytes of overhead + key size + value size

  For map[string]int (key=16B, value=8B):
    Per entry: 16 + 8 + 2 = ~26 bytes (vs 24 bytes for the raw data)
    Overhead: ~8%

  For map[int64]bool (key=8B, value=1B):
    Per entry: 8 + 1 + 2 = ~11 bytes (vs 9 bytes raw)
    But bucket alignment pads values: 8 × 1 byte values → 8 bytes
    Actual per entry: 8 + 1 + 2 = ~11 bytes
    Overhead: ~22%

  Plus: hmap struct (48 bytes) + bucket array pointer alignment
  Plus: ~19% empty slots (load factor 6.5/8 = 81% full)
```

### Map vs Slice for Small Collections

```
  ┌─────────────────────┬─────────────────────┬──────────────────────┐
  │ Collection Size     │ map[K]V             │ []struct{K,V}        │
  ├─────────────────────┼─────────────────────┼──────────────────────┤
  │ < 8 entries         │ Overkill. hmap +    │ Linear scan is fast. │
  │                     │ bucket overhead.    │ Cache-friendly. Less │
  │                     │ Hash computation    │ memory. No hashing.  │
  │                     │ cost dominates.     │ PREFER THIS.         │
  ├─────────────────────┼─────────────────────┼──────────────────────┤
  │ 8-50 entries        │ About equal.        │ Still competitive    │
  │                     │ Hash is amortized.  │ with sorted + binary │
  │                     │                     │ search.              │
  ├─────────────────────┼─────────────────────┼──────────────────────┤
  │ > 50 entries        │ Clear winner.       │ Linear scan is O(n). │
  │                     │ O(1) lookup.        │ Too slow.            │
  └─────────────────────┴─────────────────────┴──────────────────────┘
```

### Pre-Allocation: `make(map[K]V, hint)`

```go
// WITHOUT hint — 7+ growth events for 10,000 entries:
m := map[string]int{}
for i := 0; i < 10000; i++ {
    m[fmt.Sprint(i)] = i  // triggers growth at ~8, ~52, ~416, ~3328, ...
}

// WITH hint — zero growth events:
m := make(map[string]int, 10000)
for i := 0; i < 10000; i++ {
    m[fmt.Sprint(i)] = i  // all buckets pre-allocated
}
```

The hint is rounded up: `make(map[K]V, 10000)` allocates `⌈10000/6.5⌉ = 1539` buckets,
rounded to next power of 2 = 2048 buckets (B=11). This is a **hint**, not a hard cap —
the map will still grow if you exceed it. But avoiding growth eliminates the 2x memory
spike and reduces GC pressure from abandoned old bucket arrays.

**Verify with benchmarks:**
```bash
go test -bench=BenchmarkMapInsert -benchmem ./...
# Without hint: ~15 allocs/op (growth events)
# With hint:    ~1 allocs/op  (initial allocation only)
```

---

## 13. Quick Reference Card

```
MAP VARIABLE
────────────
m is *hmap (pointer to runtime.hmap struct on the heap)
Like channels: passed by pointer, zero value is nil.

BUCKET STRUCTURE
────────────────
bmap = [8]tophash + [8]keys + [8]values + *overflow
Keys and values stored separately to minimize padding.
Tophash: top 8 bits of hash → 99.6% fast rejection.

HASH & LOOKUP
─────────────
hash(key, hash0) → low B bits = bucket index, top 8 bits = tophash
Scan 8 tophash bytes → match → full key compare → return value
Type-specialized: map_fast32, map_fast64, map_faststr

GROWTH
──────
Trigger: load factor > 6.5 OR too many overflow buckets.
Double growth: 2^B → 2^(B+1), keys split by new hash bit.
Same-size growth: repack to eliminate overflow chains.
Incremental: 1-2 old buckets evacuated per insert/delete.
During growth: dual lookup (check old if not yet evacuated).

NEVER SHRINKS
─────────────
delete() clears tophash + key + value. Buckets stay allocated.
Workaround: create new map, copy survivors, drop old map.

NIL vs EMPTY
────────────
nil map:   read safe, write panics, json → "null"
empty map: read safe, write safe,  json → "{}"
Always initialize maps before writing.

CONCURRENCY
───────────
No built-in safety. Concurrent read+write → fatal (unrecoverable).
sync.Mutex: simple, works for all patterns.
sync.RWMutex: read-heavy workloads, multiple concurrent readers.
sync.Map: write-once-read-many or disjoint goroutine key sets.

KEY TYPES
─────────
Must be comparable (==). No slices, maps, or functions.
Pointers: compare address. Structs: compare all fields.
NaN float keys: unreachable after insert (NaN != NaN).

PERFORMANCE
───────────
Lookup/Delete: O(1) avg, ~50-200ns.
Insert: O(1) amortized, ~100-250ns (no growth).
Pre-allocate: make(map[K]V, hint) avoids growth events.
For < 8 entries: prefer sorted slice + linear scan.

TOOLS
─────
go test -race ./...            # detect concurrent map access
go test -bench=. -benchmem     # measure map allocation cost
go build -gcflags='-m'         # escape analysis for map values
```

---

## One-Line Summary

> A map is a `*hmap` pointer to a hash table of 8-slot buckets, using tophash bytes for
> fast rejection and incremental evacuation for growth — it never shrinks, randomizes
> iteration by design, and kills your process on concurrent read/write.
