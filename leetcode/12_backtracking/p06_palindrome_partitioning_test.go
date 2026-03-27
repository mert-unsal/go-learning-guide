package backtracking

import "testing"

func TestPartition(t *testing.T) {
	got := Partition("aab")
	if len(got) != 2 {
		t.Errorf("Partition(aab) returned %d partitions, want 2", len(got))
	}
}
