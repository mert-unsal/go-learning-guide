package hackerrank
import (
"reflect"
"testing"
)
func TestMiniMaxSum(t *testing.T) {
tests := []struct {
arr     []int
wantMin int
wantMax int
}{
{[]int{1, 2, 3, 4, 5}, 10, 14},
{[]int{7, 69, 2, 221, 8974}, 299, 9271},
{[]int{1, 1, 1, 1, 1}, 4, 4},
}
for _, tt := range tests {
gotMin, gotMax := MiniMaxSum(tt.arr)
if gotMin != tt.wantMin || gotMax != tt.wantMax {
t.Errorf("MiniMaxSum(%v) = (%d,%d), want (%d,%d)", tt.arr, gotMin, gotMax, tt.wantMin, tt.wantMax)
}
}
}
func TestFizzBuzz(t *testing.T) {
got := FizzBuzz(15)
want := []string{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"}
if !reflect.DeepEqual(got, want) {
t.Errorf("FizzBuzz(15) = %v, want %v", got, want)
}
}
func TestDiagonalDifference(t *testing.T) {
if got := DiagonalDifference([][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}); got != 0 {
t.Errorf("got %d want 0", got)
}
if got := DiagonalDifference([][]int{{11, 2, 4}, {4, 5, 6}, {10, 8, -12}}); got != 15 {
t.Errorf("got %d want 15", got)
}
}
func TestCountingValleys(t *testing.T) {
if got := CountingValleys("UDDDUDUU"); got != 1 {
t.Errorf("got %d want 1", got)
}
if got := CountingValleys("DDUUUUDD"); got != 1 {
t.Errorf("got %d want 1", got)
}
}
func TestSockMerchant(t *testing.T) {
if got := SockMerchant([]int{10, 20, 20, 10, 10, 30, 50, 10, 20}); got != 3 {
t.Errorf("got %d want 3", got)
}
}
func TestJumpingOnClouds(t *testing.T) {
if got := JumpingOnClouds([]int{0, 0, 1, 0, 0, 1, 0}); got != 4 {
t.Errorf("got %d want 4", got)
}
}
func TestRepeatedString(t *testing.T) {
if got := RepeatedString("aba", 10); got != 7 {
t.Errorf("got %d want 7", got)
}
if got := RepeatedString("a", 1000000000000); got != 1000000000000 {
t.Errorf("got %d want 1000000000000", got)
}
}
func TestCaesarCipher(t *testing.T) {
if got := CaesarCipher("middle-Outz", 2); got != "okffng-Qwvb" {
t.Errorf("got %q want okffng-Qwvb", got)
}
if got := CaesarCipher("abc", 3); got != "def" {
t.Errorf("got %q want def", got)
}
if got := CaesarCipher("xyz", 3); got != "abc" {
t.Errorf("got %q want abc", got)
}
}
func TestPangram(t *testing.T) {
if got := Pangram("We promptly judged antique ivory buckles for the next prize"); got != "pangram" {
t.Errorf("got %q want pangram", got)
}
if got := Pangram("The quick brown fox jumps over the lazy dog"); got != "pangram" {
t.Errorf("got %q want pangram", got)
}
if got := Pangram("hello world"); got != "not pangram" {
t.Errorf("got %q want not pangram", got)
}
}
func TestArrayManipulation(t *testing.T) {
if got := ArrayManipulation(5, [][]int{{1, 2, 100}, {2, 5, 100}, {3, 4, 100}}); got != 200 {
t.Errorf("got %d want 200", got)
}
if got := ArrayManipulation(10, [][]int{{1, 5, 3}, {4, 8, 7}, {6, 9, 1}}); got != 10 {
t.Errorf("got %d want 10", got)
}
}
func TestMaximumToys(t *testing.T) {
if got := MaximumToys([]int{1, 12, 5, 111, 200, 1000, 10}, 50); got != 4 {
t.Errorf("got %d want 4", got)
}
if got := MaximumToys([]int{100, 200}, 50); got != 0 {
t.Errorf("got %d want 0", got)
}
}
func TestClimbingLeaderboard(t *testing.T) {
got1 := ClimbingLeaderboard([]int{100, 100, 50, 40, 40, 20, 10}, []int{5, 25, 50, 120})
want1 := []int{6, 4, 2, 1}
if !reflect.DeepEqual(got1, want1) {
t.Errorf("test1: got %v want %v", got1, want1)
}
got2 := ClimbingLeaderboard([]int{100, 90, 90, 80}, []int{70, 90, 95})
want2 := []int{4, 2, 2}
if !reflect.DeepEqual(got2, want2) {
t.Errorf("test2: got %v want %v", got2, want2)
}
}
func TestAlmostSorted(t *testing.T) {
if got := AlmostSorted([]int{2, 1}); got != "swap 1 2" {
t.Errorf("[2,1] got %q want swap 1 2", got)
}
if got := AlmostSorted([]int{1, 5, 4, 3, 2, 6}); got != "reverse 2 5" {
t.Errorf("[1,5,4,3,2,6] got %q want reverse 2 5", got)
}
if got := AlmostSorted([]int{1, 2, 3}); got != "yes" {
t.Errorf("[1,2,3] got %q want yes", got)
}
if got := AlmostSorted([]int{3, 4, 1, 2}); got != "no" {
t.Errorf("[3,4,1,2] got %q want no", got)
}
if got := AlmostSorted([]int{3, 1, 2}); got != "no" {
t.Errorf("[3,1,2] got %q want no", got)
}
}