# 📦 Module 08 — Context: Cancellation, Timeouts & Values

> **Topics covered:** context.Context interface · cancellation propagation · timeouts · WithValue chain · typed keys · select+Done pattern · custom Context implementation
>
> **Deep dive:** [Chapter 19 — Context: The Interface Design Masterclass](../../../learnings/19_context_interface_masterclass.md)

---

## 🗺️ Learning Path

```
1. Read learnings/19_context_interface_masterclass.md   ← How context works under the hood
2. Open exercises.go                                    ← Implement the 12 exercises
3. Run go test -race -v ./...                           ← Make them all pass
```

---

## 📚 What You Will Learn

| Concept | Exercise | Under the Hood |
|---------|----------|---------------|
| Typed context keys | Ex 1 | Private key types prevent cross-package collisions |
| Context value chain | Ex 2 | Each `WithValue` creates a linked-list node — O(n) lookup |
| Cancel and check error | Ex 3 | `cancel()` closes Done channel, sets Err to `Canceled` |
| Parent→child propagation | Ex 4 | Cancel parent → all children cancelled |
| Child→parent isolation | Ex 5 | Cancel child → parent unaffected |
| Timeout = DeadlineExceeded | Ex 6 | `timerCtx` fires after duration |
| Deadline inspection | Ex 7 | `WithTimeout` is shorthand for `WithDeadline` |
| Early bailout pattern | Ex 8 | Check `ctx.Err()` before each expensive operation |
| select + ctx.Done() | Ex 9 | **THE** fundamental context pattern in production Go |
| Fan-out with cancel | Ex 10 | First result wins, cancel stops the rest |
| Nested timeouts | Ex 11 | Effective deadline = earlier of parent and child |
| Custom Context impl | Ex 12 | Implement the 4-method interface yourself |

---

## ✏️ Exercises

| # | Function | What to implement |
|---|----------|------------------|
| 1 | `WithRequestID / GetRequestID` | Store and retrieve with typed key |
| 2 | `ChainValues / LookupAll` | Build value chain, retrieve all |
| 3 | `CancelAndCheck()` | Cancel context, return error |
| 4 | `ParentCancelsChild()` | **Prove downward propagation** |
| 5 | `ChildDoesNotCancelParent()` | **Prove upward isolation** |
| 6 | `TimeoutExpired(d)` | Wait for timeout, return DeadlineExceeded |
| 7 | `CheckDeadline(d)` | Inspect deadline on a timed context |
| 8 | `ProcessItems(ctx, items)` | **Early bailout on cancellation** |
| 9 | `SelectWithContext(ctx, ch)` | **select + ctx.Done() pattern** |
| 10 | `FirstResult(tasks)` | **Fan-out: first wins, cancel rest** |
| 11 | `NestedTimeout(t1, t2)` | Effective deadline with nested timeouts |
| 12 | `AlwaysCancelled()` | **Implement context.Context interface** |

---

## 🧪 Run Tests

```bash
go test -race -v ./exercises/stdlib/08_context/
```

---

## ✅ Done? Next Step

```bash
go test -race -v ./exercises/stdlib/09_sync/
```
