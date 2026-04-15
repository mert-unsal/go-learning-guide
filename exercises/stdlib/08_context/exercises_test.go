package context_exercises

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: WithRequestID / GetRequestID
// ────────────────────────────────────────────────────────────

func TestWithRequestID(t *testing.T) {
	tests := []struct {
		id string
	}{
		{"req-abc-123"},
		{""},
		{"trace-00000"},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			ctx := WithRequestID(context.Background(), tt.id)
			got := GetRequestID(ctx)
			if got != tt.id {
				t.Errorf("❌ GetRequestID after WithRequestID(%q) = %q, want %q\n\t\t"+
					"Hint: Use context.WithValue(ctx, RequestKey{}, id) to store, "+
					"ctx.Value(RequestKey{}).(string) to retrieve. "+
					"Private key types prevent collisions. See learnings/19 §7",
					tt.id, got, tt.id)
			} else {
				t.Logf("✅ GetRequestID = %q", got)
			}
		})
	}

	// Verify missing key returns empty string
	t.Run("missing_key", func(t *testing.T) {
		got := GetRequestID(context.Background())
		if got != "" {
			t.Errorf("❌ GetRequestID on Background() = %q, want \"\"\n\t\t"+
				"Hint: ctx.Value() returns nil for missing keys. "+
				"Type-assert with comma-ok: v, ok := ctx.Value(key).(string)",
				got)
		} else {
			t.Logf("✅ GetRequestID on Background() = \"\" (correct)")
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 2: ChainValues / LookupAll
// ────────────────────────────────────────────────────────────

func TestChainValues(t *testing.T) {
	pairs := [][2]string{
		{"userID", "42"},
		{"traceID", "abc-123"},
		{"role", "admin"},
	}
	ctx := ChainValues(pairs)

	t.Run("lookup_all", func(t *testing.T) {
		got := LookupAll(ctx, []string{"userID", "traceID", "role"})
		for _, p := range pairs {
			if got[p[0]] != p[1] {
				t.Errorf("❌ LookupAll[%q] = %q, want %q\n\t\t"+
					"Hint: Each WithValue creates a linked-list node. "+
					"Value() walks the chain — O(n). See learnings/19 §3",
					p[0], got[p[0]], p[1])
			} else {
				t.Logf("✅ LookupAll[%q] = %q", p[0], p[1])
			}
		}
	})

	t.Run("missing_key", func(t *testing.T) {
		got := LookupAll(ctx, []string{"nonexistent"})
		if _, exists := got["nonexistent"]; exists {
			t.Errorf("❌ LookupAll[\"nonexistent\"] should not exist in result\n\t\t"+
				"Hint: ctx.Value() returns nil for missing keys. Omit from map")
		} else {
			t.Logf("✅ missing key correctly omitted")
		}
	})

	t.Run("empty_chain", func(t *testing.T) {
		ctx := ChainValues(nil)
		got := LookupAll(ctx, []string{"anything"})
		if len(got) != 0 {
			t.Errorf("❌ LookupAll on empty chain should return empty map, got %v", got)
		} else {
			t.Logf("✅ empty chain returns empty map")
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 3: CancelAndCheck
// ────────────────────────────────────────────────────────────

func TestCancelAndCheck(t *testing.T) {
	err := CancelAndCheck()
	if err != context.Canceled {
		t.Errorf("❌ CancelAndCheck() = %v, want context.Canceled\n\t\t"+
			"Hint: After calling cancel(), ctx.Err() returns context.Canceled. "+
			"ctx.Done() channel is also closed. See learnings/19 §3",
			err)
	} else {
		t.Logf("✅ CancelAndCheck() = context.Canceled")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: ParentCancelsChild
// ────────────────────────────────────────────────────────────

func TestParentCancelsChild(t *testing.T) {
	parentErr, childErr := ParentCancelsChild()

	if parentErr != context.Canceled {
		t.Errorf("❌ parent.Err() = %v, want context.Canceled\n\t\t"+
			"Hint: Cancelling parent sets parent.Err() to Canceled",
			parentErr)
	} else {
		t.Logf("✅ parent.Err() = context.Canceled")
	}

	if childErr != context.Canceled {
		t.Errorf("❌ child.Err() = %v, want context.Canceled\n\t\t"+
			"Hint: Cancellation propagates DOWNWARD. Child of cancelled parent "+
			"is also cancelled. See learnings/19 §7",
			childErr)
	} else {
		t.Logf("✅ child.Err() = context.Canceled (propagated from parent)")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: ChildDoesNotCancelParent
// ────────────────────────────────────────────────────────────

func TestChildDoesNotCancelParent(t *testing.T) {
	parentErr, childErr := ChildDoesNotCancelParent()

	if parentErr != nil {
		t.Errorf("❌ parent.Err() = %v, want nil\n\t\t"+
			"Hint: Cancellation does NOT propagate UPWARD. "+
			"Cancelling a child leaves the parent alive. See learnings/19 §7",
			parentErr)
	} else {
		t.Logf("✅ parent.Err() = nil (not affected by child cancel)")
	}

	if childErr != context.Canceled {
		t.Errorf("❌ child.Err() = %v, want context.Canceled",
			childErr)
	} else {
		t.Logf("✅ child.Err() = context.Canceled")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: TimeoutExpired
// ────────────────────────────────────────────────────────────

func TestTimeoutExpired(t *testing.T) {
	start := time.Now()
	err := TimeoutExpired(50 * time.Millisecond)
	elapsed := time.Since(start)

	if err != context.DeadlineExceeded {
		t.Errorf("❌ TimeoutExpired(50ms) = %v, want context.DeadlineExceeded\n\t\t"+
			"Hint: WithTimeout fires after the duration. <-ctx.Done() blocks until then. "+
			"Err() returns DeadlineExceeded (not Canceled). "+
			"Always defer cancel() to release the timer. See learnings/19 §3",
			err)
	} else {
		t.Logf("✅ TimeoutExpired(50ms) = context.DeadlineExceeded (took %v)", elapsed)
	}

	if elapsed < 40*time.Millisecond {
		t.Errorf("❌ returned too fast (%v) — did you actually wait for the timeout?", elapsed)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: CheckDeadline
// ────────────────────────────────────────────────────────────

func TestCheckDeadline(t *testing.T) {
	before := time.Now()
	deadline, ok := CheckDeadline(5 * time.Second)

	if !ok {
		t.Errorf("❌ CheckDeadline(5s) ok = false, want true\n\t\t"+
			"Hint: WithTimeout creates a timerCtx that overrides Deadline(). "+
			"It returns (deadline, true). Background() returns (zero, false). "+
			"See learnings/19 §3")
		return
	}
	t.Logf("✅ CheckDeadline(5s) ok = true")

	// Deadline should be approximately 5s from now
	expected := before.Add(5 * time.Second)
	diff := deadline.Sub(expected).Abs()
	if diff > 100*time.Millisecond {
		t.Errorf("❌ deadline is %v away from expected — should be ~5s from now", diff)
	} else {
		t.Logf("✅ deadline is within 100ms of expected")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: ProcessItems
// ────────────────────────────────────────────────────────────

func TestProcessItems(t *testing.T) {
	t.Run("no_cancel", func(t *testing.T) {
		ctx := context.Background()
		items := []int{1, 2, 3, 4, 5}
		got := ProcessItems(ctx, items)
		want := []int{1, 4, 9, 16, 25}
		if fmt.Sprint(got) != fmt.Sprint(want) {
			t.Errorf("❌ ProcessItems(Background, [1..5]) = %v, want %v\n\t\t"+
				"Hint: Square each item. With no cancellation, all items are processed",
				got, want)
		} else {
			t.Logf("✅ ProcessItems with no cancel = %v", got)
		}
	})

	t.Run("cancel_after_3", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// Cancel after a tiny delay — ProcessItems should check ctx.Err()
		// We'll cancel immediately and expect at most 0 items
		cancel()
		got := ProcessItems(ctx, items)

		if len(got) > 0 {
			t.Errorf("❌ ProcessItems after cancel: got %d items, want 0\n\t\t"+
				"Hint: Check ctx.Err() BEFORE processing each item. "+
				"If cancelled, stop and return what you have. See learnings/19 §7",
				len(got))
		} else {
			t.Logf("✅ ProcessItems after cancel = [] (immediate bailout)")
		}
	})

	t.Run("empty_items", func(t *testing.T) {
		got := ProcessItems(context.Background(), nil)
		if len(got) != 0 {
			t.Errorf("❌ ProcessItems(nil) = %v, want empty", got)
		} else {
			t.Logf("✅ ProcessItems(nil) = []")
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 9: SelectWithContext
// ────────────────────────────────────────────────────────────

func TestSelectWithContext(t *testing.T) {
	t.Run("value_arrives", func(t *testing.T) {
		ctx := context.Background()
		ch := make(chan string, 1)
		ch <- "hello"

		got, err := SelectWithContext(ctx, ch)
		if err != nil {
			t.Errorf("❌ SelectWithContext returned error %v when value available\n\t\t"+
				"Hint: select on ctx.Done() and ch. If ch sends, return (val, nil)", err)
		} else if got != "hello" {
			t.Errorf("❌ SelectWithContext = %q, want \"hello\"", got)
		} else {
			t.Logf("✅ SelectWithContext = \"hello\"")
		}
	})

	t.Run("context_cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // cancel immediately
		ch := make(chan string) // unbuffered, nothing will send

		_, err := SelectWithContext(ctx, ch)
		if err != context.Canceled {
			t.Errorf("❌ SelectWithContext returned %v, want context.Canceled\n\t\t"+
				"Hint: select { case <-ctx.Done(): return \"\", ctx.Err() ... }\n\t\t"+
				"This is THE fundamental context pattern. See learnings/19 §7",
				err)
		} else {
			t.Logf("✅ SelectWithContext = context.Canceled")
		}
	})

	t.Run("timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		ch := make(chan string) // never sends

		_, err := SelectWithContext(ctx, ch)
		if err != context.DeadlineExceeded {
			t.Errorf("❌ SelectWithContext with timeout returned %v, want DeadlineExceeded", err)
		} else {
			t.Logf("✅ SelectWithContext with timeout = DeadlineExceeded")
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 10: FirstResult
// ────────────────────────────────────────────────────────────

func TestFirstResult(t *testing.T) {
	t.Run("fastest_wins", func(t *testing.T) {
		tasks := []func() int{
			func() int { time.Sleep(200 * time.Millisecond); return 1 },
			func() int { return 42 }, // instant — should win
			func() int { time.Sleep(200 * time.Millisecond); return 3 },
		}
		got := FirstResult(tasks)
		if got != 42 {
			t.Errorf("❌ FirstResult = %d, want 42 (the instant task)\n\t\t"+
				"Hint: Launch each task in a goroutine, send result to buffered channel, "+
				"return first received. Use context.WithCancel to signal others to stop. "+
				"Buffer the channel (cap=len(tasks)) so slower goroutines don't block forever",
				got)
		} else {
			t.Logf("✅ FirstResult = 42 (instant task won)")
		}
	})

	t.Run("single_task", func(t *testing.T) {
		tasks := []func() int{
			func() int { return 99 },
		}
		got := FirstResult(tasks)
		if got != 99 {
			t.Errorf("❌ FirstResult = %d, want 99", got)
		} else {
			t.Logf("✅ FirstResult single task = 99")
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 11: NestedTimeout
// ────────────────────────────────────────────────────────────

func TestNestedTimeout(t *testing.T) {
	t.Run("parent_shorter", func(t *testing.T) {
		// Parent: 1s, Child: 5s → effective = ~1s (parent wins)
		got := NestedTimeout(1*time.Second, 5*time.Second)
		if got > 1100*time.Millisecond || got < 800*time.Millisecond {
			t.Errorf("❌ NestedTimeout(1s, 5s) = %v, want ~1s\n\t\t"+
				"Hint: Child deadline cannot exceed parent's deadline. "+
				"WithTimeout on a parent with shorter deadline → parent wins. "+
				"See learnings/19 §3 — timerCtx",
				got)
		} else {
			t.Logf("✅ NestedTimeout(1s, 5s) = %v (parent wins)", got)
		}
	})

	t.Run("child_shorter", func(t *testing.T) {
		// Parent: 5s, Child: 1s → effective = ~1s (child wins)
		got := NestedTimeout(5*time.Second, 1*time.Second)
		if got > 1100*time.Millisecond || got < 800*time.Millisecond {
			t.Errorf("❌ NestedTimeout(5s, 1s) = %v, want ~1s\n\t\t"+
				"Hint: When child's timeout is shorter, it wins. "+
				"The effective deadline is always the earlier of parent and child",
				got)
		} else {
			t.Logf("✅ NestedTimeout(5s, 1s) = %v (child wins)", got)
		}
	})
}

// ────────────────────────────────────────────────────────────
// Exercise 12: AlwaysCancelled
// ────────────────────────────────────────────────────────────

func TestAlwaysCancelled(t *testing.T) {
	ctx := AlwaysCancelled()

	t.Run("err", func(t *testing.T) {
		if ctx.Err() != context.Canceled {
			t.Errorf("❌ AlwaysCancelled().Err() = %v, want context.Canceled\n\t\t"+
				"Hint: Implement the Context interface (4 methods). "+
				"Err() returns context.Canceled. See learnings/19 §4",
				ctx.Err())
		} else {
			t.Logf("✅ Err() = context.Canceled")
		}
	})

	t.Run("done_closed", func(t *testing.T) {
		select {
		case <-ctx.Done():
			t.Logf("✅ Done() channel is closed (select didn't block)")
		default:
			t.Errorf("❌ Done() channel is not closed — select hit default\n\t\t"+
				"Hint: Create a channel with make(chan struct{}), close it immediately. "+
				"Return it from Done(). A closed channel always succeeds in select")
		}
	})

	t.Run("deadline", func(t *testing.T) {
		_, ok := ctx.Deadline()
		if ok {
			t.Errorf("❌ Deadline() ok = true, want false")
		} else {
			t.Logf("✅ Deadline() ok = false (no deadline)")
		}
	})

	t.Run("value", func(t *testing.T) {
		if ctx.Value("anything") != nil {
			t.Errorf("❌ Value() should return nil")
		} else {
			t.Logf("✅ Value() = nil")
		}
	})
}
