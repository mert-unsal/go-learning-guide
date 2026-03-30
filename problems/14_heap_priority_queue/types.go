package heap_priority_queue

type IntMaxHeap []int

func (h IntMaxHeap) Len() int {
	return 0
}
func (h IntMaxHeap) Less(i, j int) bool {
	return false
}
func (h IntMaxHeap) Swap(i, j int) {
}
func (h *IntMaxHeap) Push(x interface{}) {
}
func (h *IntMaxHeap) Pop() interface{} {
	return nil
}

type IntMinHeap []int

func (h IntMinHeap) Len() int {
	return 0
}
func (h IntMinHeap) Less(i, j int) bool {
	return false
}
func (h IntMinHeap) Swap(i, j int) {
}
func (h *IntMinHeap) Push(x interface{}) {
}
func (h *IntMinHeap) Pop() interface{} {
	return nil
}
