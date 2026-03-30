package interfaces

import (
	"bytes"
	"fmt"
	"io"
)

func main() {
	var w io.Writer
	fmt.Println("A:", w == nil)

	var buf *bytes.Buffer
	w = buf
	fmt.Println("B:", w == nil)
	fmt.Println("C:", buf == nil)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("D: PANIC ->", r)
		}
	}()
	w.Write([]byte("boom"))
}
