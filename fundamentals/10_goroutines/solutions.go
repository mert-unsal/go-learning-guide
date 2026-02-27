package goroutines
import "sync"
// ============================================================
// SOLUTIONS — 10 Goroutines
// ============================================================
func RunConcurrentlySolution(n int, fn func(id int)) {
var wg sync.WaitGroup
for i := 0; i < n; i++ {
wg.Add(1)
i := i // capture loop variable — critical in Go < 1.22
go func() {
defer wg.Done()
fn(i)
}()
}
wg.Wait()
}
func (c *ExCounter) IncSolution() {
c.mu.Lock()
defer c.mu.Unlock()
c.value++
}
func (c *ExCounter) ValueSolution() int {
c.mu.Lock()
defer c.mu.Unlock()
return c.value
}
func SumConcurrentSolution(nums []int) int {
if len(nums) == 0 {
return 0
}
mid := len(nums) / 2
var mu sync.Mutex
total := 0
var wg sync.WaitGroup
sumHalf := func(slice []int) {
defer wg.Done()
s := 0
for _, v := range slice {
s += v
}
mu.Lock()
total += s
mu.Unlock()
}
wg.Add(2)
go sumHalf(nums[:mid])
go sumHalf(nums[mid:])
wg.Wait()
return total
}
func RunOnceSolution(setup func()) {
var once sync.Once
once.Do(setup) // runs setup
once.Do(setup) // does nothing — already ran
once.Do(setup) // does nothing — already ran
}