package channels
import (
"sync"
"time"
)
// ============================================================
// SOLUTIONS — 11 Channels
// ============================================================
func SumAsyncSolution(nums []int, ch chan<- int) {
sum := 0
for _, v := range nums {
sum += v
}
ch <- sum // send result — caller receives it
}
func GenerateSolution(n int) <-chan int {
ch := make(chan int)
go func() {
defer close(ch) // always close from the sender side
for i := 1; i <= n; i++ {
ch <- i
}
}()
return ch
}
func SquareSolution(in <-chan int) <-chan int {
out := make(chan int)
go func() {
defer close(out)
for v := range in { // range exits when in is closed
out <- v * v
}
}()
return out
}
func MergeSolution(a, b <-chan int) <-chan int {
out := make(chan int)
var wg sync.WaitGroup
forward := func(ch <-chan int) {
defer wg.Done()
for v := range ch {
out <- v
}
}
wg.Add(2)
go forward(a)
go forward(b)
// Close output when both inputs are drained
go func() {
wg.Wait()
close(out)
}()
return out
}
func WithTimeoutSolution(ch <-chan int, maxWaitMs int) (int, bool) {
select {
case v := <-ch:
return v, true
case <-time.After(time.Duration(maxWaitMs) * time.Millisecond):
return 0, false
}
}