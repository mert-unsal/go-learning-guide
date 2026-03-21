package interfaces

import (
	"fmt"
	"reflect"
	"strconv"
)
// SOLUTIONS — 06 Interfaces
func (b ExBook) StringSolution() string {
return fmt.Sprintf("%q by %s", b.Title, b.Author)
}
func (m ExMovie) StringSolution() string {
return fmt.Sprintf("%s (%d)", m.Title, m.Year)
}
func PrintAllSolution(items []ExStringer) {
for _, item := range items {
fmt.Println(item.String())
}
}
func (bw *ExBufferWriter) WriteSolution(data string) error {
bw.Buffer = append(bw.Buffer, data)
return nil
}
func WriteAllSolution(w ExWriter, items []string) error {
for _, item := range items {
if err := w.Write(item); err != nil {
return err
}
}
return nil
}
func DescribeSolution(i interface{}) string {
switch v := i.(type) {
case int:
return "int: " + strconv.Itoa(v)
case string:
return "string: " + v
case bool:
if v { return "bool: true" }
return "bool: false"
default:
return "unknown"
}
}

// ─── Exercise 4a Solution: SafeGetLabel (Guard 1) ───

func SafeGetLabelSolution(name string) Labeler {
	if name == "" {
		return nil // true nil interface: iface{nil, nil}
	}
	return &Product{Name: name}
}

// ─── Exercise 4b Solution: SafeCallLabeler (Guard 2) ───

func SafeCallLabelerSolution(l Labeler) string {
	if l == nil {
		return "no labeler" // true nil: iface{nil, nil}
	}
	p, ok := l.(*Product)
	if !ok || p == nil {
		return "no labeler" // typed nil: iface{*itab, nil}
	}
	return p.Label()
}

// ─── Exercise 4c Solution: IsTrulyNil (Guard 3) ───

func IsTrulyNilSolution(i interface{}) bool {
	if i == nil {
		return true // true nil: eface{nil, nil}
	}
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface,
		reflect.Slice, reflect.Map,
		reflect.Chan, reflect.Func:
		return v.IsNil() // nillable kinds — safe to call IsNil()
	}
	return false // non-nillable kinds (int, string, struct, etc.)
}

// ─── Exercise 5a Solution: Thermometer ───

func (t Thermometer) ReadingSolution() float64 {
	return t.Temp
}

func (t *Thermometer) CalibrateSolution(offset float64) {
	t.Temp += offset
}

// ─── Exercise 5b Solution: ReadAndCalibrate ───

func ReadAndCalibrateSolution(s Sensor, offset float64) (before, after float64) {
	before = s.Reading()
	s.Calibrate(offset)
	after = s.Reading()
	return before, after
}

// ─── Exercise 5c Solution: Celsius & Kelvin formatters ───

func (c Celsius) DisplaySolution() string {
	return fmt.Sprintf("%.1f°C", float64(c))
}

func (k *Kelvin) DisplaySolution() string {
	return fmt.Sprintf("%.1fK", float64(*k))
}

func CollectDisplayersSolution() []Displayer {
	k := Kelvin(309.8) // must store in variable to take address
	return []Displayer{
		Celsius(36.6), // value receiver → T satisfies Displayer, no & needed
		&k,            // pointer receiver → only *Kelvin satisfies Displayer
	}
}