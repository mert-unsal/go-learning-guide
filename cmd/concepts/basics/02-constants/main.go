// Constants in Go — demonstrates const, iota, and bitwise flag patterns.
//
// Constants are declared with 'const' and CANNOT be changed.
// They must be known at compile time.
package main

import "fmt"

const MaxSize = 100
const AppName = "LearnGo"
const Pi = 3.14159265358979

// Weekday --- iota: auto-incrementing constant generator ---
type Weekday int

const (
	Sunday    Weekday = iota // 0
	Monday                   // 1
	Tuesday                  // 2
	Wednesday                // 3
	Thursday                 // 4
	Friday                   // 5
	Saturday                 // 6
)

// iota with bit shifting — very common for flags/permissions
type Permission uint

const (
	Read    Permission = 1 << iota // 1  (001)
	Write                          // 2  (010)
	Execute                        // 4  (100)
)

func main() {
	fmt.Println("MaxSize:", MaxSize)
	fmt.Println("Monday:", Monday) // prints 1
	fmt.Println("Read:", Read, "Write:", Write, "Execute:", Execute)

	// Checking permissions with bitwise AND
	myPerms := Read | Write                        // 3 (011)
	fmt.Println("Can read?", myPerms&Read != 0)    // true
	fmt.Println("Can exec?", myPerms&Execute != 0) // false
}
